// Package util util.go
package util

import (
	"fmt"
	"log"
	//"webshellManager/database"
	//"webshellManager/database"
)

var (
	HttpURL     string
	Passwd      string
	UserAgent   string
	ContentType = "application/x-www-form-urlencoded"
)

func Init(httpURl, passwd, userAgent string) {
	HttpURL = httpURl
	Passwd = passwd
	UserAgent = userAgent
	//fmt.Println(HttpURL, Passwd, UserAgent)
	//test := "成功!"
	//return test
	//return HttpURL, Passwd, UserAgent, err

}

func HandleError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func PrintLogo() {
	logo := `
      ____            _
     / ___| ___   ___| | __
    | |  _ / _ \ / __| |/ /
    | |_| | (_) | (__|   <
     \____|\___/ \___|_|\_\
    `
	fmt.Printf("%s\n", logo)
	fmt.Println("--help\t显示完整帮助信息\n--code\t输入要执行的 PHP 代码(省略`;`)\n--shell\t利用 system 函数执行系统命令\n--generate-shell 生成简单的web_shell\n--dbinfo 显示目前数据库信息\n--adddb\t添加数据\n--phpinfo\t查看 PHP 禁用的函数\n------------华丽的分割线-----------")
}
