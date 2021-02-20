package mysql

import (
	"demo-to-start/model"
	"log"
)

// GetUser 数据库查询
func GetUser(userID int) (model.User, error) {
	var res model.User
	// 这个是单行返回的方法,Scan指定列与字段对应顺序
	err := db.QueryRow("select * from `users` where id=? ", userID).Scan(&res.ID, &res.Name, &res.Age)
	if err != nil {
		log.Println(err)
		return res, err
	}
	return res, nil
}
