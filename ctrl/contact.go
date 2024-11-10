package ctrl

import (
	"IM/args"
	"IM/model"
	"IM/service"
	"IM/util"
	"net/http"
)

var contactService service.ContactService

func LoadFriendHandler(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	//如果这个用的上,那么可以直接
	util.Bind(req, &arg)

	users := contactService.SearchFriend(arg.Userid)
	util.RespOkList(w, users, len(users))
}

func LoadCommunityHandler(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	//如果这个用的上,那么可以直接
	util.Bind(req, &arg)
	comunitys := contactService.SearchCommunity(arg.Userid)
	util.RespOkList(w, comunitys, len(comunitys))
}
func JoinCommunityHandler(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg

	//如果这个用的上,那么可以直接
	util.Bind(req, &arg)
	err := contactService.JoinCommunity(arg.Userid, arg.Dstid)
	// 加群的时候刷新群组信息
	AddGroupId(arg.Userid, arg.Dstid)

	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "")
	}
}

func CreateCommunityHandler(w http.ResponseWriter, req *http.Request) {
	var arg model.Community
	util.Bind(req, &arg)
	com, err := contactService.CreatCommunity(arg)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, com, "")
	}
}

func AddFriendHandler(w http.ResponseWriter, req *http.Request) {
	//定义一个参数结构体
	/*request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")
	*/
	var arg args.ContactArg
	util.Bind(req, &arg)
	//调用service
	err := contactService.AddFriend(arg.Userid, arg.Dstid)
	//
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "好友添加成功")
	}
}
