package handlers

import (
	"demo-to-start/mysql"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type QueryUserRequest struct {
	UserID int `json:"user_id"`
}
// QueryUser 请求处理逻辑
func QueryUser(w http.ResponseWriter, r *http.Request) {
	bd, err := io.ReadAll(r.Body)
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
	user, err := mysql.GetUser(req.UserID)
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
