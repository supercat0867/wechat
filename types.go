package wechat

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// ErrorHandler 错误处理
func ErrorHandler(action, errmsg string, errcode int) error {
	return fmt.Errorf("%s:%s,错误码：%d", action, errmsg, errcode)
}

var (
	ErrGetAccessToken        = "access_token获取失败"
	ErrSendTempMessage       = "模版消息发送失败"
	ErrSendTextMessage       = "客服文本消息发送失败"
	ErrGetUserList           = "用户列表获取失败"
	ErrGetUserInfo           = "用户基础信息获取失败"
	ErrGetWebAuthAccessToken = "网页授权access_token获取失败"
)

type MessageType string

// 消息类型
const (
	TextMessage       MessageType = "text"       // 文本消息
	VoiceMessage      MessageType = "voice"      // 语音消息
	VideoMessage      MessageType = "video"      // 视频消息
	ShortVideoMessage MessageType = "shortvideo" // 小视频消息
	LocationMessage   MessageType = "location"   // 地理位置消息
	LinkMessage       MessageType = "link"       // 链接消息
	EventMessage      MessageType = "event"      // 事件消息
)

type MessageHandler func(msg *Message, w http.ResponseWriter)

type SDK struct {
	handlers map[MessageType]MessageHandler
}

// XMLMessage 微信xml消息格式
type XMLMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`            // 开发者微信号
	FromUserName string   `xml:"FromUserName"`          // 发送方账号（一个OpenID）
	CreateTime   int64    `xml:"CreateTime"`            // 消息创建时间 （整型）
	MsgType      string   `xml:"MsgType"`               // 消息类型，文本为text
	Content      string   `xml:"Content"`               // 文本消息内容
	MsgId        int64    `xml:"MsgId"`                 // 消息id，64位整型
	MsgDataId    string   `xml:"MsgDataId,omitempty"`   // 消息的数据ID（消息如果来自文章时才有）
	Idx          string   `xml:"Idx,omitempty"`         // 多图文时第几篇文章，从1开始（消息如果来自文章时才有）
	PicUrl       string   `xml:"PicUrl,omitempty"`      // 图片链接（由系统生成）
	MediaId      string   `xml:"MediaId,omitempty"`     // 图片消息媒体id或语音消息媒体id，可以调用获取临时素材接口拉取数据。
	Format       string   `xml:"Format,omitempty"`      // 语音格式，如amr，speex等
	Recognition  string   `xml:"Recognition,omitempty"` // 语音识别结果，UTF8编码 (已废弃)
	Event        string   `xml:"Event,omitempty"`       // 事件类型
}

type Message struct {
	Type         MessageType // 消息类型
	Content      string      // 消息内容
	ToUserName   string      // 开发者微信号
	FromUserName string      // 发送方openid
	MediaId      string      // 素材ID
	Event        string      // 事件类型
}

// AccessTokenResponse access_token响应
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"` // 获取到的凭证
	ExpiresIn   int    `json:"expires_in"`   // 凭证有效时间，单位：秒
	Errcode     int    `json:"errcode"`      // 错误码
	ErrMsg      string `json:"errmsg"`       // 错误信息
}

// SendTempMessageResponse 发送模版消息响应
type SendTempMessageResponse struct {
	Errcode int    `json:"errcode"` // 错误码
	ErrMsg  string `json:"errmsg"`  // 错误信息
	MsgID   int    `json:"msgid"`   // 消息ID
}

// TempMessage 模版消息通用格式
type TempMessage struct {
	ToUser      string                     `json:"touser"`        // 接收者openid
	TemplateID  string                     `json:"template_id"`   // 模板ID
	URL         string                     `json:"url"`           // 模板跳转链接（海外账号没有跳转能力）
	MiniProgram TempMessageMiniProgram     `json:"miniprogram"`   // 跳小程序所需数据，不需跳小程序可不用传该数据
	ClientMsgID string                     `json:"client_msg_id"` // 防重入id。对于同一个openid + client_msg_id, 只发送一条消息,10分钟有效,超过10分钟不保证效果。若无防重入需求，可不填
	Data        map[string]TempMessageData `json:"data"`          // 模板数据
}
type TempMessageData struct {
	Value string `json:"value"`
}
type TempMessageMiniProgram struct {
	AppID    string `json:"app_id"`   // 所需跳转到的小程序appid（该小程序appid必须与发模板消息的公众号是绑定关联关系，暂不支持小游戏）
	PagePath string `json:"pagePath"` // 所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar），要求该小程序已发布，暂不支持小游戏
}

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
