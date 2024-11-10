package model

import "time"

// 好友和群都存在这个表里，可以根据具体业务做拆分

const (
	CONCAT_CATE_USER     = 0x01 //用户
	CONCAT_CATE_COMUNITY = 0x02 //群组
)

type Contact struct {
	Id       int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	OwnerId  int64     `xorm:"bigint(20)" form:"ownerid" json:"ownerid"` // 本端id
	DstObj   int64     `xorm:"bigint(20)" form:"dstobj" json:"dstobj"`   // 对端id
	Cate     int       `xorm:"int(11)" form:"cate" json:"cate"`          // 区分用户互相加，还是加群
	Memo     string    `xorm:"varchar(120)" form:"memo" json:"memo"`     // 描述
	CreateAt time.Time `xorm:"datetime" form:"createat" json:"createat"` // 创建时间，用于统计使用
}
