// Package http http.go
package http

import (
	"net/http"
	"strings"
	"webshellManager/util"
)

func MakeRequest(payload, url, ua string) (*http.Response, error) {
	//fmt.Println("11111111111111111111:\t", url)
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	{
		util.HandleError(err, "发送payload失败!")
	}
	//fmt.Println(req)
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Content-Type", util.ContentType)

	client := &http.Client{}
	return client.Do(req)
}
