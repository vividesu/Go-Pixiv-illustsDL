package utils

import (
	"net/http"
	"net/url"
)

var (
	u, err = url.Parse("http://127.0.0.1:7890") // 代理（这里用的是自己的代理
	// 使用代理
	client = http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(u),
		},
	}
)

// 设置请求头爬取
func Http_Client_SetAndDo(AimUrl string) (http_resp *http.Response) {
	req, _ := http.NewRequest("GET", AimUrl, nil)
	req.Header.Add("referer", "https://www.pixiv.net/")
	resp, _ := client.Do(req)
	return resp
}
