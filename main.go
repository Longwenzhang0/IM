package main

import (
	"IM/ctrl"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

// 模板函数

func registerView() {
	// ** 表示一级目录，*表示一个文件
	tpl, err := template.ParseGlob("view/**/*")
	// 此处会返回四个值， /user/login.shtml和login.html 和 ，/user/register.shtml和register.html
	// 而login.html是无法解析的，需要添加过滤条件，将末尾是html的项目跳过
	if err != nil {
		log.Fatal(err.Error()) // 渲染出错直接打印并退出
	}

	// 遍历每一个模板
	for _, v := range tpl.Templates() {
		tplname := v.Name()
		if strings.Contains(tplname, ".html") {
			// 添加过滤条件。只渲染.shtml，而.html就跳过；
			continue
		}
		fmt.Printf("[%s tplname:%s ] \n", time.Now().Format("2006/01/02 15:04:05"), tplname)
		http.HandleFunc(tplname, func(writer http.ResponseWriter, request *http.Request) {
			tpl.ExecuteTemplate(writer, tplname, nil)
		})

	}

}

func main() {
	// 绑定请求和处理函数
	http.HandleFunc("/user/login", ctrl.UserLoginHandler)
	http.HandleFunc("/user/register", ctrl.UserRegisterHandler)
	http.HandleFunc("/contact/loadcommunity", ctrl.LoadCommunityHandler)
	http.HandleFunc("/contact/loadfriend", ctrl.LoadFriendHandler)
	http.HandleFunc("/contact/joincommunity", ctrl.JoinCommunityHandler)
	http.HandleFunc("/contact/createcommunity", ctrl.CreateCommunityHandler)
	http.HandleFunc("/contact/addfriend", ctrl.AddFriendHandler)
	http.HandleFunc("/chat", ctrl.ChatHandler)
	http.HandleFunc("/attach/upload", ctrl.UploadHandler)

	// 1. 提供静态资源目录支持
	//http.Handle("/", http.FileServer(http.Dir(".")))
	// 结果暴露了main.go http://localhost:8080/main.go，
	// 提供指定目录的静态文件支持
	http.Handle("/asset/", http.FileServer(http.Dir(".")))
	http.Handle("/mnt/", http.FileServer(http.Dir(".")))

	registerView()

	// 启动Web服务器

	http.ListenAndServe(":8080", nil)

}
