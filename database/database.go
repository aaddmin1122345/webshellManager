// Package database database.go
package database

import (
	"database/sql"
	"fmt"
	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
	"webshellManager/util"
)

// 数据库连接
var db *sql.DB

func init() {
	var err error
	db, err = connectDb()
	util.HandleError(err, "数据库连接失败!")
}

func connectDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	return db, err
}

// CreateTable 创建数据库表
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

	_, err = db.Exec(insertDataSQL, "http://test.test/eval.php", "cmd", "test_ua", "备注信息")
	util.HandleError(err, "插入数据出错!")
}

// AddURL 添加URL
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

// ALLDB 查询全部数据库
func ALLDB() {
	resultRows, err := db.Query("SELECT id,url, passwd, ua FROM info; ")
	util.HandleError(err, "查询数据库错误!")
	listDB(resultRows)
	//CloseDB()
}

// SelectDb 选择数据库记录
func SelectDb() {
	fmt.Println("以下是数据库中的全部信息!")
	ALLDB()
	var id int
	fmt.Print("输入要操作的数据库id: ")
	_, err := fmt.Scanln(&id)
	util.HandleError(err, "查询id失败!")
	resultRows, err := db.Query("SELECT id, url, passwd,ua FROM info WHERE id = ?", id)
	util.HandleError(err, "查询数据库错误!")
	listDB(resultRows)
	CloseDB()

}

// CloseDB 关闭数据库连接
func CloseDB() {
	//err := db.Close()
	defer func(db *sql.DB) {
		err := db.Close()
		util.HandleError(err, "关闭数据库连接失败!")
	}(db)

}

func listDB(resultRows *sql.Rows) {
	// 创建颜色对象，用于将ID标记为绿色
	green := color.New(color.FgGreen)
	fmt.Println("+----+--------------------------------+----------+---------+")
	fmt.Println("| ID | URL                            | Password | UA      |")
	fmt.Println("+----+--------------------------------+----------+---------+")

	for resultRows.Next() {
		var id int
		var url, passwd, ua string
		err := resultRows.Scan(&id, &url, &passwd, &ua)
		util.HandleError(err, "遍历数据库错误!!!")
		// 使用颜色对象输出绿色的ID和对齐的列
		fmt.Printf("| %-2s%-1s | %-30s | %-8s | %-7s |\n", green.Sprint(id), "", url, passwd, ua)
		fmt.Println("+----+--------------------------------+----------+---------+")

	}

}
