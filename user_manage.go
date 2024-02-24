package wechat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 用户管理模块

// GetUserListResponse 获取用户列表响应
type GetUserListResponse struct {
	Total int `json:"total"` // 关注该公众账号的总用户数
	Count int `json:"count"` // 拉取的OPENID个数，最大值为10000
	Data  struct {
		OpenID []string `json:"openid"`
	} `json:"data"` // 列表数据，OPENID的列表
	NextOpenID string `json:"next_openid"` // 拉取列表的最后一个用户的OPENID
	Errcode    int    `json:"errcode"`     // 错误码
	ErrMsg     string `json:"errmsg"`      // 错误信息
}

// GetUserList 获取用户列表
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/User_Management/Getting_a_User_List.html
func GetUserList(accessToken, nextOpenID string) (*GetUserListResponse, error) {
	// 接口地址
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/get?access_token=%s&next_openid=%s",
		accessToken, nextOpenID)

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
	var responseJson GetUserListResponse
	if err = json.Unmarshal(body, &responseJson); err != nil {
		return nil, err
	}
	if responseJson.Errcode != 0 {
		return nil, ErrorHandler(ErrGetUserList, responseJson.ErrMsg, responseJson.Errcode)
	}
	return &responseJson, nil
}

// GetUserInfoResponse 获取用户基本信息响应
type GetUserInfoResponse struct {
	Subscribe      int    `json:"subscribe"`       // 用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。
	OpenID         string `json:"openid"`          // 用户的标识，对当前公众号唯一
	Language       string `json:"language"`        // 用户的语言，简体中文为zh_CN
	SubscribeTime  int    `json:"subscribe_time"`  // 用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
	UnionID        string `json:"unionid"`         // 只有在用户将公众号绑定到微信开放平台账号后，才会出现该字段。
	Remark         string `json:"remark"`          // 公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	GroupID        int    `json:"groupid"`         // 用户所在的分组ID（兼容旧的用户分组接口）
	TagIDList      []int  `json:"tagid_list"`      // 用户被打上的标签ID列表
	SubScribeScene string `json:"subscribe_scene"` // 返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENE_PROFILE_LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_WECHAT_ADVERTISEMENT 微信广告，ADD_SCENE_REPRINT 他人转载 ,ADD_SCENE_LIVESTREAM 视频号直播，ADD_SCENE_CHANNELS 视频号 , ADD_SCENE_OTHERS 其他
	QRScene        int    `json:"qr_scene"`        // 二维码扫码场景（开发者自定义）
	QRSceneStr     string `json:"qr_scene_str"`    // 二维码扫码场景描述（开发者自定义）
	Errcode        int    `json:"errcode"`         // 错误码
	ErrMsg         string `json:"errmsg"`          // 错误信息
}

// GetUserInfo 获取用户基本信息
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
func GetUserInfo(accessToken, openID string) (*GetUserInfoResponse, error) {
	// 接口地址
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN",
		accessToken, openID)

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
	var responseJson GetUserInfoResponse
	if err = json.Unmarshal(body, &responseJson); err != nil {
		return nil, err
	}
	if responseJson.Errcode != 0 {
		return nil, ErrorHandler(ErrGetUserInfo, responseJson.ErrMsg, responseJson.Errcode)
	}
	return &responseJson, nil
}
