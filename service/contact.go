package service

import (
	"IM/model"
	"errors"
	"log"
	"time"
)

type ContactService struct {
}

// 自动添加好友

func (service *ContactService) AddFriend(
	userid,
	dstid int64) error {
	// 添加好友，入参为自己id和对端id，返回error类型
	// 如果自己加自己，提示失败
	if userid == dstid {
		return errors.New("不能添加自己为好友")
	}
	// 判断是否已经加了好友
	tmp := model.Contact{}
	// 查询是否已经是好友，条件的链式操作
	_, err := dbEngine.Where("OwnerId = ?", userid).
		And("DstObj = ?", dstid).
		And("Cate = ?", model.CONCAT_CATE_USER).
		Get(&tmp)
	//if err != nil {
	//	fmt.Println("==============>", err)
	//	return errors.New("查询好友关系失败")
	//}
	// 获得一条纪录
	// count() 的效率更低，使用get更快
	// 如果存在记录，则无法添加
	if tmp.Id > 0 {
		return errors.New("该用户已经被添加过了")
	}
	// 启动事务，如果要插入两条数据，必须要两条都插入成功，事务才算成功
	session := dbEngine.NewSession()
	err = session.Begin()
	if err != nil {
		return errors.New("启动事务失败")
	}
	// 插入自己的数据
	_, e2 := session.InsertOne(model.Contact{
		OwnerId:  userid,
		DstObj:   dstid,
		Cate:     model.CONCAT_CATE_USER,
		CreateAt: time.Now(),
	})
	// 插入对端的数据
	_, e3 := session.InsertOne(model.Contact{
		OwnerId:  dstid,
		DstObj:   userid,
		Cate:     model.CONCAT_CATE_USER,
		CreateAt: time.Now(),
	})
	// 判断插入数据的情况，都不返错，则添加成功，提交事务；
	if e2 == nil && e3 == nil {
		err := session.Commit()
		if err != nil {
			return errors.New("提交事务失败")
		}
		return nil
	} else {
		// 回滚
		err := session.Rollback()
		if err != nil {
			return errors.New("事务回滚失败")
		}
		// 判断返回error
		if e2 != nil {
			// 这里如果两个都插入失败了，只返回了e3
			return e2
		} else {
			return e3
		}
	}
}

// 搜索群聊
func (service *ContactService) SearchCommunity(userId int64) []model.Community {
	// 传入本端用户id，返回我加的多个群聊的切片；
	// 两步走，先查我的联系人，查出来id之后，再根据id查群的表
	contacts := make([]model.Contact, 0)
	comIds := make([]int64, 0)
	// 根据我的user_id号和联系人类型(群)，
	// 查询我自己的联系人，存到contacts里；对于contacts来说，对方的id，就是群的号码；
	err := dbEngine.Where("OwnerId = ? and Cate = ?", userId, model.CONCAT_CATE_COMUNITY).Find(&contacts)
	if err != nil {
		log.Println("query community error: " + err.Error())
	}
	for _, v := range contacts {
		// 将群号存到comID里；
		comIds = append(comIds, v.DstObj)
	}
	// 创建一个空的群的切片
	coms := make([]model.Community, 0)
	if len(comIds) == 0 {
		return coms
	}
	// 查询id列，查询comIds里存的数据，再到群聊的表里的去查所有的群；
	dbEngine.In("id", comIds).Find(&coms)
	return coms
}

func (service *ContactService) SearchCommunityIds(userId int64) (comIds []int64) {
	// 获取用户全部群id
	// 定义一个contact的切片和一个群号码的切片
	contacts := make([]model.Contact, 0)
	comIds = make([]int64, 0)
	dbEngine.Where("OwnerId = ? and Cate = ?", userId, model.CONCAT_CATE_COMUNITY).Find(&contacts)
	for _, v := range contacts {
		comIds = append(comIds, v.DstObj)
	}
	return comIds
}

// 加群
func (service *ContactService) JoinCommunity(userId, comID int64) error {
	// 入参，用户id和群id，出参error信息
	// 新增一个联系人，本端id就是自己的id，对端id是群号，类型为群聊
	cot := model.Contact{
		OwnerId: userId,
		DstObj:  comID,
		Cate:    model.CONCAT_CATE_COMUNITY,
	}
	dbEngine.Get(&cot)
	// 查询是否已经有这个群聊了；没有的话就插入一条数据
	if cot.Id == 0 {
		cot.CreateAt = time.Now()
		_, err := dbEngine.InsertOne(cot)
		return err
	} else {
		return errors.New("已经加群了")
	}
}

// 建群
func (servce *ContactService) CreatCommunity(comm model.Community) (ret model.Community, err error) {
	if len(comm.Name) == 0 {
		err = errors.New("缺少群名称")
		return ret, err
	}
	if comm.OwnerId == 0 {
		err = errors.New("请先登录")
		return ret, err
	}
	// 实例化一个群com
	com := model.Community{
		OwnerId: comm.OwnerId,
	}
	// 查询本人id创建了多少个群
	num, err := dbEngine.Count(&com)
	if num > 5 {
		err = errors.New("一个用户最多只能创建5个群")
		return com, err
	} else {
		// 先插入一个群聊项目
		comm.CreateAt = time.Now()
		session := dbEngine.NewSession()
		session.Begin()
		_, err = session.InsertOne(&comm)
		if err != nil {
			session.Rollback()
			return com, err
		}
		// 再插入一个联系人项目
		_, err = session.InsertOne(model.Contact{
			OwnerId:  comm.OwnerId,
			DstObj:   comm.Id,
			Cate:     model.CONCAT_CATE_COMUNITY,
			CreateAt: time.Now(),
		})
		if err != nil {
			session.Rollback()
		} else {
			session.Commit()
		}
		return com, err
	}
}

// 查找好友
func (service *ContactService) SearchFriend(userId int64) []model.User {
	// 创建一个空的联系人切片，id切片
	contacts := make([]model.Contact, 0)
	objIds := make([]int64, 0)
	dbEngine.Where("ownerid = ? and cate = ?", userId, model.CONCAT_CATE_USER).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, v.DstObj)
	}
	coms := make([]model.User, 0)
	if len(objIds) == 0 {
		return coms
	}
	dbEngine.In("id", objIds).Find(&coms)
	return coms
}
