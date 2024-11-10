package service

import (
	"IM/model"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"log"
	"xorm.io/core"
)

var dbEngine *xorm.Engine

// 安装方式
// go get github.com/go-xorm/xorm
// go get github.com/go-sql-driver/mysql

func init() {
	// 初始化函数
	driverName := "mysql"
	// 用户名，密码，数据库地址，后面可以添加字符模式
	dataSourceName := "root:root@(127.0.0.1:3306)/chat?charset=utf8"
	err := errors.New("")
	dbEngine, err = xorm.NewEngine(driverName, dataSourceName)

	if err != nil && err.Error() != "" {
		log.Fatal(err.Error())
	}
	// 显示操作过程中的sql语句
	dbEngine.ShowSQL(true)
	// 设置数据库最大链接数
	dbEngine.SetMaxOpenConns(2)
	// 设置从字段映射到数据库列名的方式
	dbEngine.SetMapper(core.SameMapper{})

	// 自动user，同步数据库
	err = dbEngine.Sync2(new(model.User),
		new(model.Contact),
		new(model.Community))
	if err != nil {
		fmt.Println("init data base failed")
		return
	}

	fmt.Println("init data base ok")

}
