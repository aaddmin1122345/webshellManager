// Package http http.go
package http

import (
	"net/http"
	"strings"
	"webshellManager/util"
)

const (
	userAgent   = `Mozilla/5.0 ...` // 保留您的常量
	contentType = "application/x-www-form-urlencoded"
	httpURL     = "http://test.test/eval.php"
)

func MakeRequest(payload string) (*http.Response, error) {
	req, err := http.NewRequest("POST", httpURL, strings.NewReader(payload))
	util.HandleError(err, "发送payload失败!")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	return client.Do(req)
}
