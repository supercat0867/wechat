package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 模版消息模块

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

// NewTemMessage 实例化模版消息
func NewTemMessage(touser, templateID, url, appID, appPagePath, clientMsgID string, msgData map[string]string) *TempMessage {
	var data = make(map[string]TempMessageData)
	for key, value := range msgData {
		data[key] = TempMessageData{value}
	}
	return &TempMessage{
		ToUser:      touser,
		TemplateID:  templateID,
		URL:         url,
		MiniProgram: TempMessageMiniProgram{AppID: appID, PagePath: appPagePath},
		ClientMsgID: clientMsgID,
		Data:        data,
	}
}

// Send 发送模版消息
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
func (m *TempMessage) Send(accessToken string) error {
	// 将消息数据序列化为JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// 创建请求
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", accessToken)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/json")

	// 创建一个 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解析响应
	var responseJson SendTempMessageResponse
	if err = json.Unmarshal(body, &responseJson); err != nil {
		return err
	}

	if responseJson.Errcode == 0 {
		return nil
	}

	return ErrorHandler(ErrSendTempMessage, responseJson.ErrMsg, responseJson.Errcode)
}
