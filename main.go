package main

import (
	"demo-to-start/handlers"
	"demo-to-start/mysql"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Lshortfile)
	// 数据库初始化
	err := mysql.RegisterDB()
	if err != nil {
		log.Println(err)
		return
	}
	// 服务初始化
	Server()
}

// Server 网络服务
func Server() {
	// 1.注册一个处理器函数,这里没有限制Get/Post等http方法
	http.HandleFunc("/get_user", handlers.QueryUser)

	// 2.设置监听的TCP地址并启动服务
	// 参数1:TCP地址(IP+Port)
	// 参数2:handler handler参数一般会设为nil，此时会使用DefaultServeMux。
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Println("http.ListenAndServe()函数执行错误,错误为:", err.Error())
		return
	}
}
