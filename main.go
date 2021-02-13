package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // 注册驱动，不需要内部方法
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	ID   int
	Name string
	Age  int
}

var db *sql.DB

func main() {
	log.SetFlags(log.Lshortfile)
	// 数据库初始化
	err := registerDB()
	if err != nil {
		log.Println(err)
		return
	}
	// 服务初始化
	Server()
}
func registerDB() error {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(10.105.11.29:3306)/shm")
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

// GetUser 数据库查询
func GetUser(userID int) (User, error) {
	var res User
	// 这个是单行返回的方法,Scan指定列与字段对应顺序
	err := db.QueryRow("select * from `users` where id=? ", userID).Scan(&res.ID, &res.Name, &res.Age)
	if err != nil {
		log.Println(err)
		return res, err
	}
	return res, nil
}

type QueryUserRequest struct {
	UserID int `json:"user_id"`
}
type Response struct {
	Code int
	Msg  string
	Data interface{}
}

// QueryUser 请求处理逻辑
func QueryUser(w http.ResponseWriter, r *http.Request) {
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var req QueryUserRequest
	err = json.Unmarshal(bd, &req)
	if err != nil {
		log.Println(err)
		return
	}
	if req.UserID <= 0 {
		_, _ = w.Write([]byte("userID is invalid"))
		return
	}
	user, err := GetUser(req.UserID)
	if err != nil {
		log.Println(err)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	resp := Response{Code: 1, Msg: "success", Data: user}
	bd, err = json.Marshal(resp)
	if err != nil {
		log.Println(err)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	// 设置响应头
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bd)
}

// Server 网络服务
func Server() {
	// 1.注册一个处理器函数,这里没有限制Get/Post等http方法
	http.HandleFunc("/get_user", QueryUser)

	// 2.设置监听的TCP地址并启动服务
	// 参数1:TCP地址(IP+Port)
	// 参数2:handler handler参数一般会设为nil，此时会使用DefaultServeMux。
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Println("http.ListenAndServe()函数执行错误,错误为:", err.Error())
		return
	}
}
