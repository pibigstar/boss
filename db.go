package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// ====== db相关操作 ======
var db *sql.DB

//  获取mysql连接
func connectMysql() {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8mb4,utf8",
		cfg.Mysql.UserName, cfg.Mysql.Password, cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.DB)

	var err error
	db, err = sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
}

// 获取有效的用户
func listUserFromDB() ([]*User, error) {
	rows, err := db.Query("select id,username,cookie,status from user where status = 1")
	if err != nil {
		return nil, err
	}
	var users []*User
	for rows.Next() {
		var u User
		err = rows.Scan(&u.Id, &u.UserName, &u.Cookie, &u.Status)
		if err != nil {
			log.Println("list user rows scan", err.Error())
			continue
		}
		users = append(users, &u)
	}
	return users, nil
}

func listJobsFromDB() {

}
