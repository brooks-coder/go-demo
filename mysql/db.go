package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // 注册驱动，不需要内部方法
)

var db *sql.DB

func RegisterDB() error {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(10.105.11.29:3306)/shm")
	if err != nil {
		return err
	}
	return nil
}