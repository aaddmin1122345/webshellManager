// Package database database.go
package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"webshellManager/util"
)

var db *sql.DB // 全局数据库连接

func init() {
	var err error
	db, err = connectDb()
	util.HandleError(err, "数据库连接失败!")
}

func connectDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	return db, err
}

func CreateTable() {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS info (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT NOT NULL,
        passwd TEXT NOT NULL,
        ua TEXT NOT NULL,
        other TEXT NOT NULL
    );`

	_, err := db.Exec(createTableSQL)
	util.HandleError(err, "创建数据库出错!")

	insertDataSQL := `
    INSERT INTO info (url, passwd, ua, other) VALUES (?, ?, ?, ?);`

	_, err = db.Exec(insertDataSQL, "http://test.test", "cmd", "test_ua", "备注信息")
	util.HandleError(err, "插入数据出错!")
}

func AddURL() {
	var url, passwd, ua, other string
	fmt.Printf("输入添加的url:\t")
	_, err := fmt.Scanln(&url)
	util.HandleError(err, "添加url失败!")

	fmt.Printf("输入添加的passwd:\t")
	_, err = fmt.Scanln(&passwd)
	util.HandleError(err, "添加passwd失败!")

	fmt.Printf("输入添加的ua:\t")
	_, err = fmt.Scanln(&ua)
	util.HandleError(err, "添加ua失败!")

	fmt.Printf("输入备注:\t")
	_, err = fmt.Scanln(&other)
	util.HandleError(err, "添加备注失败!")

	insertDataSQL := `
    INSERT INTO info (url, passwd, ua, other) VALUES (?, ?, ?, ?);`
	_, err = db.Exec(insertDataSQL, url, passwd, ua, other)
	util.HandleError(err, "插入数据出错!")

	fmt.Println("数据已成功添加!")
}

func DbAll() {
	rows, err := db.Query("SELECT id,url, passwd, ua FROM info; ")
	util.HandleError(err, "查询数据库错误!")
	defer func(rows *sql.Rows) {
		err := rows.Close()
		util.HandleError(err, "关闭数据库错误!")
	}(rows)

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

	rows, err := db.Query("SELECT url, passwd, ua FROM info WHERE id = ?", id)
	util.HandleError(err, "查询数据库错误!")
	defer func(rows *sql.Rows) {
		err := rows.Close()
		util.HandleError(err, "关闭数据库错误!")
	}(rows)

	for rows.Next() {
		err := rows.Scan(&url, &passwd, &ua)
		util.HandleError(err, "遍历数据库内容出错!")
		fmt.Printf("id:%d url:%s passwd:%s ua:%s\n", id, url, passwd, ua)
	}
}
