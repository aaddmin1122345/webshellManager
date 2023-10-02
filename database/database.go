// Package database database.go
package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"webshellManager/util"
)

func connectDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	util.HandleError(err, "打开数据库错误!")
	return db, err
}

func CreateTable() {
	db, err := connectDb()
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS info (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT NOT NULL,
        passwd TEXT NOT NULL,
        ua TEXT NOT NULL,
        other TEXT NOT NULL
    );`

	_, err = db.Exec(createTableSQL)
	util.HandleError(err, "创建数据库出错!")

	insertDataSQL := `
    INSERT INTO info (url, passwd, ua, other) VALUES (?, ?, ?, ?);`

	_, err = db.Exec(insertDataSQL, "http://test.test", "cmd", "test_ua", "备注信息")
	util.HandleError(err, "插入数据出错!")

	defer db.Close()
}

func AddURL() {
	var url, passwd, ua, other string
	fmt.Printf("输入添加的url:\t")
	fmt.Scanln(&url)
	fmt.Printf("输入添加的passwd:\t")
	fmt.Scanln(&passwd)
	fmt.Printf("输入添加的ua:\t")
	fmt.Scanln(&ua)
	fmt.Printf("输入备注:\t")
	fmt.Scanln(&other)
	db, err := connectDb()
	insertDataSQL := `
    INSERT INTO info (url, passwd, ua, other) VALUES (?, ?, ?, ?);`
	_, err = db.Exec(insertDataSQL, url, passwd, ua, other)
	util.HandleError(err, "插入数据出错!")
	fmt.Println("数据已成功添加!")
	defer db.Close()
}

func DbAll() {
	db, err := connectDb()
	util.HandleError(err, "连接数据库出错!")
	defer db.Close()

	rows, err := db.Query("select id,url, passwd, ua from info; ")
	util.HandleError(err, "查询数据库错误!")
	defer rows.Close()

	for rows.Next() {
		var url, passwd, ua string
		var id int
		err := rows.Scan(&id, &url, &passwd, &ua)
		util.HandleError(err, "遍历数据库内容出错!")

		fmt.Printf("id:%d url:%s passwd:%s ua:%s\n", id, url, passwd, ua)
	}
}

func SelectDb() {
	fmt.Println("以下是数据库中的全部信息!")
	DbAll()
	var id int
	var url, passwd, ua string
	fmt.Print("输入要查询的ID: ")
	_, err := fmt.Scanln(&id)
	util.HandleError(err, "查询id失败!")

	db, err := connectDb()
	util.HandleError(err, "连接数据库出错!")
	defer db.Close()

	rows, err := db.Query("select url, passwd, ua from info where id = ?", id)
	util.HandleError(err, "查询数据库错误!")
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&url, &passwd, &ua)
		util.HandleError(err, "遍历数据库内容出错!")

		fmt.Printf("id:%d url:%s passwd:%s ua:%s\n", id, url, passwd, ua)
	}
}
