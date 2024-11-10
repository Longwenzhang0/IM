package service

import (
	"IM/model"
	"IM/util"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// 模块化。使用结构体
type UserService struct {
}

// 注册函数
func (s *UserService) Register(
	mobile, // 手机号
	plainpwd, // 明文密码
	nickname, // 昵称
	avatar, sex string) (user model.User, err error) {

	// 检测手机号是否存在
	tmp := model.User{}
	_, err = dbEngine.Where("mobile=? ", mobile).Get(&tmp)
	if err != nil {
		fmt.Printf("error-----------> %s", err)
		return tmp, err // 错误处理
	}
	// 如果存在，则返回提示已经注册
	if tmp.Id > 0 {
		return tmp, errors.New("该手机号已经注册")
	}
	// 否则拼接插入数据
	tmp.Mobile = mobile
	// passwd =
	// md5加密，将明文密码加密之后存入数据库；
	tmp.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	tmp.Passwd = util.MakePasswd(plainpwd, tmp.Salt)
	tmp.Nickname = nickname
	tmp.Avatar = avatar
	tmp.Sex = sex
	tmp.CreateAt = time.Now()
	// 新建token
	tmp.Token = fmt.Sprintf("%08d", rand.Int31())
	// 返回新用户信息

	// 插入数据
	_, err = dbEngine.InsertOne(&tmp)
	// 前端恶意插入特殊字符
	// 数据库连接操作失败

	return tmp, err
}

// 登录函数
func (s *UserService) Login(
	mobile,
	plainpwd string) (user model.User, err error) {

	tmp := model.User{}
	_, err = dbEngine.Where("mobile=? ", mobile).Get(&tmp)
	if err != nil {
		return tmp, err // 错误处理
	}
	// 如果存在，对比密码
	if tmp.Id > 0 {
		if !util.ValidatePasswd(plainpwd, tmp.Salt, tmp.Passwd) {
			// 密码错误
			return tmp, errors.New("密码错误")
		}
		// 密码正确的情况下， 用户登录一次刷新token，安全；
		str := fmt.Sprintf("%d", time.Now().Unix())
		token := util.MD5Encode(str)
		tmp.Token = token
		// 重新写入token
		_, err := dbEngine.ID(tmp.Id).Cols("token").Update(&tmp)
		if err != nil {
			return tmp, errors.New("重写token失败")
		}
		return tmp, nil

	} else {
		return tmp, errors.New("该用户不存在")
	}

}

func (s *UserService) Find(userId int64) (user model.User) {
	// // 查找某个用户根据id查询token
	tmp := model.User{}
	_, err := dbEngine.ID(userId).Get(&tmp)
	if err != nil {
		log.Println(err.Error())
		return tmp
	}
	return tmp

}
