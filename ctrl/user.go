package ctrl

import (
	"IM/model"
	"IM/service"
	"IM/util"
	"fmt"
	"math/rand"
	"net/http"
)

func UserLoginHandler(writer http.ResponseWriter, request *http.Request) {
	// 可以做：
	// 数据库操作
	// 逻辑操作
	// rest aip json/xml返回

	// 如何获取参数
	request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")

	user, err := userService.Login(mobile, passwd)
	if err != nil {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, user, "")
	}

	// curl http://127.0.0.1:8080/user/login -X POST -d "mobile=1234&passwd=5678"
	// 在cmd中使用此命令验证，X POST 指定请求方法为 POST。使用 -d 指定 POST请求数据
}

var userService service.UserService

func UserRegisterHandler(writer http.ResponseWriter, request *http.Request) {
	// 测试
	// curl http://127.0.0.1:8080/user/register -d "mobile=13452131581&passwd=123456"
	request.ParseForm()

	mobile := request.PostForm.Get("mobile")
	plainpwd := request.PostForm.Get("passwd")
	nickname := fmt.Sprintf("user%06d", rand.Int31())
	avatar := ""
	sex := model.SEX_UNKONWN

	user, err := userService.Register(mobile, plainpwd, nickname, avatar, sex)

	//fmt.Printf("err===========>%s", err)

	if err != nil {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, user, "")
	}
}
