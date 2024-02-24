package wechat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetWebAuthAccessTokenResponse 通过code换取网页授权access_token的响应
type GetWebAuthAccessTokenResponse struct {
	AccessToken    string `json:"access_token"`    // 网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	Expires        string `json:"expires_in"`      // access_token接口调用凭证超时时间，单位（秒）
	RefreshToken   string `json:"refresh_token"`   // 用户刷新access_token
	OpenID         string `json:"openid"`          // 用户唯一标识，请注意，在未关注公众号时，用户访问公众号的网页，也会产生一个用户和公众号唯一的OpenID
	Scope          string `json:"scope"`           // 用户授权的作用域，使用逗号（,）分隔
	IsSnapShotUser int    `json:"is_snapshotuser"` // 是否为快照页模式虚拟账号，只有当用户是快照页模式虚拟账号时返回，值为1
	UnionID        string `json:"unionid"`         // 用户统一标识（针对一个微信开放平台账号下的应用，同一用户的 unionid 是唯一的），只有当scope为"snsapi_userinfo"时返回
	Errcode        int    `json:"errcode"`         // 错误码
	ErrMsg         string `json:"errmsg"`          // 错误信息
}

// GetWebAuthAccessToken 获取网页授权access_token
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html#1
// 说明：此功能需要的权限较高，需要在微信公众号后台配置相关信息使用，详细使用方法流程请参考官方文档
func GetWebAuthAccessToken(appID, appSecret, code string) (*GetWebAuthAccessTokenResponse, error) {
	// 接口地址
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		appID, appSecret, code)

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
	var responseJson GetWebAuthAccessTokenResponse
	if err = json.Unmarshal(body, &responseJson); err != nil {
		return nil, err
	}
	if responseJson.Errcode != 0 {
		return nil, ErrorHandler(ErrGetWebAuthAccessToken, responseJson.ErrMsg, responseJson.Errcode)
	}
	return &responseJson, nil
}
