package main

import (
	"container/list"
	"demo-to-start/handlers"
	"demo-to-start/mysql"
	"log"
	"net/http"
)

type LRUCache struct {
	Cap  int
	Keys map[int]*list.Element
	List *list.List
}

type pair struct{ K, V int }

func Constructor(capacity int) LRUCache {
	return LRUCache{Cap: capacity, Keys: make(map[int]*list.Element), List: list.New()}
}

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
