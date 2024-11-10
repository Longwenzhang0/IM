package ctrl

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// 消息的类型
const (
	CMD_SINGLE_MSG = 10
	CMD_ROOM_MSG   = 11
	CMD_HEART      = 0
)

// 消息结构体
type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"`           // 消息的id
	Userid  int64  `json:"userid,omitempty" form:"userid"`   // 谁发的
	Cmd     int    `json:"cmd,omitempty" form:"cmd"`         // 群聊还是私聊
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`     // 对端用户id或者群聊id
	Media   int    `json:"media,omitempty" form:"media"`     // 消息以什么样式来展示
	Content string `json:"content,omitempty" form:"content"` // 消息的内容
	Pic     string `json:"pic,omitempty" form:"pic"`         // 预览的图片
	Url     string `json:"url,omitempty" form:"url"`         // 服务的URL
	Memo    string `json:"memo,omitempty" form:"memo"`       // 简单描述
	Amount  int    `json:"amount,omitempty" form:"amount"`   // 其他和数字相关的

}

// 核心在于形成userid到Node的映射关系
type Node struct {
	Conn *websocket.Conn
	// 并行转串行
	DataQueue chan []byte   // 用于存储数据。这个通道用于并行转串行，即缓冲发送给 WebSocket 客户端的数据。
	GroupSets set.Interface //  GroupSets：类型为 set.Interface，是一个线程安全的集合
}

// 映射关系表 userid <-----> Node
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwlocker sync.RWMutex

// 一般请求是这个格式：127.0.0.1/chat?id=1&token=xxxx
func ChatHandler(writer http.ResponseWriter, request *http.Request) {
	// 检验接入是否合法

	query := request.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	// 将string类型的id转为int64
	userId, _ := strconv.ParseInt(id, 10, 64)
	isValid := checkToken(userId, token)

	// 如果校验合法:将一个普通的HTTP请求升级为WebSocket连接，
	//以便可以在客户端和服务器之间进行双向通信。
	//这通常用于需要实时通信的Web应用，如聊天应用、游戏或实时数据传输等场景。
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValid
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 获得连接
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	// 获取用户的全部群id;用户接入时初始化groupset
	comIds := contactService.SearchCommunityIds(userId)
	// 刷新groupSet
	for _, v := range comIds {
		node.GroupSets.Add(v)
	}

	// userid 和 node绑定关系
	// map的操作频率很高，需要加读写锁；保证并发不出错
	rwlocker.Lock()
	clientMap[userId] = node
	rwlocker.Unlock()
	// 完成发送逻辑
	go sendProc(node)
	// 完成接收逻辑
	go recvProc(node)

	// 测试下
	sendMsg(userId, []byte("Hello First Message!"))

}

// 添加新的群id到用户的group set中
func AddGroupId(userId, gid int64) {
	// 取得node
	rwlocker.Lock()
	node, ok := clientMap[userId]
	if ok {
		node.GroupSets.Add(gid)
	}
	rwlocker.Unlock()
}

// 发送协程
func sendProc(node *Node) {
	for {
		select { // select 语句用于等待多个通道操作中的一个完成。等待从 node.DataQueue 通道中接收数据。没有break；会一直运行
		case data := <-node.DataQueue: // 将数据从DataQueue发送到data
			err := node.Conn.WriteMessage(websocket.TextMessage, data) // 写入数据
			if err != nil {
				// 出错打印信息并且返回
				log.Println(err.Error())
				return
			}
		}
	}
}

// 接收协程
func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}
		// 对data进一步处理
		//fmt.Printf("recv<=%s\n", data)
		dispatch(data)
		// 把消息广播到局域网
		//broadMsg(data)
		log.Printf("[ws]<==%s\n", data)
	}
}

func init() {
	go udpsendproc()
	go udprecvproc()
}

// 用于存放将要广播的数据
var udpsendchan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendchan <- data
}

func udpsendproc() {
	log.Println("start udpsendproc")
	// 使用udp协议拨号
	con, err := net.DialUDP("udp", nil,
		&net.UDPAddr{
			IP:   net.IPv4(192, 170, 0, 255),
			Port: 3333,
		})
	defer con.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 通过得到的con发送消息
	for {
		select {
		case data := <-udpsendchan:
			_, err = con.Write(data)
			log.Println("i am udp send proc--")
			if err != nil {
				log.Println(err.Error())
				return
			}

		}
	}

}

func udprecvproc() {
	log.Println("start udprecvproc")
	//todo 监听udp广播端口
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3333,
	})
	defer con.Close()
	if err != nil {
		log.Println(err.Error())
	}
	//TODO 处理端口发过来的数据
	for {
		log.Println("i am udp recv proc--")
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			log.Println("recv error is: " + err.Error())
			return
		}
		//直接数据处理
		dispatch(buf[0:n])
	}
	log.Println("stop updrecvproc")
}

// 发送消息
func sendMsg(userId int64, msg []byte) {
	// 入参是发送给谁，对端id，和消息内容
	rwlocker.RLock() // 读锁，保证map并发安全性
	node, ok := clientMap[userId]
	rwlocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}

}

func checkToken(userId int64, token string) bool {
	// 检验token是否合法，从数据库中查询并且比对
	user := userService.Find(userId)
	return user.Token == token
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	// 解析data为message
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 根据cmd对逻辑进行处理
	switch msg.Cmd {
	case CMD_SINGLE_MSG:
		sendMsg(msg.Dstid, data)
	case CMD_ROOM_MSG:
		// todo 群聊
		for _, v := range clientMap {
			if v.GroupSets.Has(msg.Dstid) {
				// 遍历所有的映射表，判断是否有群id，有的话就发过去
				v.DataQueue <- data
			}
		}
	case CMD_HEART:
		// 心跳 一般啥都不做

	}

	// 根据cmd对逻辑进行处理

}
