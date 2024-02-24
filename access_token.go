package wechat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AccessTokenResponse access_token响应
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"` // 获取到的凭证
	ExpiresIn   int    `json:"expires_in"`   // 凭证有效时间，单位：秒
	Errcode     int    `json:"errcode"`      // 错误码
	ErrMsg      string `json:"errmsg"`       // 错误信息
}

// GetAccessToken 获取access_token
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
func GetAccessToken(appID, appSecret string) (*AccessTokenResponse, error) {
	// 接口地址
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		appID, appSecret)

	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应到请求体
	var responseJson AccessTokenResponse
	if err = json.Unmarshal(body, &responseJson); err != nil {
		return nil, err
	}

	if responseJson.Errcode == 0 {
		return &responseJson, nil
	}

	return nil, ErrorHandler(ErrGetAccessToken, responseJson.ErrMsg, responseJson.Errcode)
}
