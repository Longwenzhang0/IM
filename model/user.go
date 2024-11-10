package model

import "time"

const (
	SEX_WOMAN   = "W"
	SEX_MAN     = "M"
	SEX_UNKONWN = "U"
)

// model.SEX_WOMAN

type User struct {
	Id       int64     `xorm:"pk autoincr bigint(64)" form:"id" json:"id"`  // 用户id
	Mobile   string    `xorm:"varchar(20)" form:"mobile" json:"mobile"`     // 用户手机号
	Passwd   string    `xorm:"varchar(40)" form:"passwd" json:"-"`          // 密码=f(plainpwd + salt) salt为随机数，加密md5
	Avatar   string    `xorm:"varchar(150)" form:"avatar" json:"avatar"`    // 头像
	Sex      string    `xorm:"varchar(2)" form:"sex" json:"sex"`            // 性别
	Nickname string    `xorm:"varchar(20)" form:"nickname" json:"nickname"` // 昵称
	Salt     string    `xorm:"varchar(10)" form:"salt" json:"-"`            // 随机数
	Online   int       `xorm:"int(10)" form:"online" json:"online"`         //是否在线
	Token    string    `xorm:"varchar(40)" form:"token" json:"token"`       // /chat?id=1&token=x
	Memo     string    `xorm:"varchar(140)" form:"memo" json:"memo"`        //
	CreateAt time.Time `xorm:"datetime" form:"createat" json:"createat"`    // 统计每天用户增量；数据创建时间
}
