package wechat

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// NewMessageSDK 实例化sdk
func NewMessageSDK(appid, appsecret string) *SDK {
	sdk := &SDK{
		handlers:  make(map[MessageType]MessageHandler),
		AppID:     appid,
		AppSecret: appsecret,
	}

	// 自动维护accesstoken
	go func() {
		for {
			resp, err := sdk.GetAccessToken()
			if err != nil {
				log.Println(err)
			} else {
				sdk.AccessToken = resp.AccessToken
			}
			time.Sleep(time.Hour)
		}
	}()

	return sdk
}

// RegisterHandler 注册消息处理方法
func (s *SDK) RegisterHandler(msgType MessageType, handler MessageHandler) {
	s.handlers[msgType] = handler
}

// 解析微信xml消息到结构体
func parseWeChatMessage(data []byte) (*XMLMessage, error) {
	var msg XMLMessage
	err := xml.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// HandleWeChatMessage 处理消息
func (s *SDK) HandleWeChatMessage(data []byte, w http.ResponseWriter) {
	msg, err := parseWeChatMessage(data)
	if err != nil {
		// 处理错误
		log.Printf("xml解析失败:%v", err)
		return
	}

	genericMsg := &Message{
		ToUserName:   msg.ToUserName,
		FromUserName: msg.FromUserName,
	}

	switch msg.MsgType {
	case "text":
		genericMsg.Type = TextMessage
		genericMsg.Content = msg.Content
	case "voice":
		genericMsg.Type = VoiceMessage
		// 语音自动转文字能力被官方移除
		genericMsg.Content = msg.Recognition
		genericMsg.MediaId = msg.MediaId
	case "event":
		genericMsg.Type = EventMessage
		genericMsg.Event = msg.Event
	// 添加其他消息类型的转换
	default:
		// 处理未知消息类型
		return
	}

	// 调用对应类型的处理器
	if handler, ok := s.handlers[genericMsg.Type]; ok {
		handler(genericMsg, w)
	}
}

// BuildTextResponse 构造被动回复文本消息xml
func (s *SDK) BuildTextResponse(toUser, fromUser, content string) string {
	return fmt.Sprintf(`<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>`, toUser, fromUser, time.Now().Unix(), content)
}

// SendTextMessage 发送文本消息
// 官方文档地址：https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#%E5%AE%A2%E6%9C%8D%E6%8E%A5%E5%8F%A3-%E5%8F%91%E6%B6%88%E6%81%AF
func (s *SDK) SendTextMessage(toUser, content string) error {
	data := map[string]interface{}{
		"touser":  toUser,
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": content,
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 创建请求
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", s.AccessToken)
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

	return ErrorHandler(ErrSendMiniprogramMessage, responseJson.ErrMsg, responseJson.Errcode)
}

// SendMiniprogramMessage 发送小程序卡片
// 官方文档地址：https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#%E5%AE%A2%E6%9C%8D%E6%8E%A5%E5%8F%A3-%E5%8F%91%E6%B6%88%E6%81%AF
func (s *SDK) SendMiniprogramMessage(toUser, title, appid, pagePath, mediaId string) error {
	data := map[string]interface{}{
		"touser":  toUser,
		"msgtype": "miniprogrampage",
		"miniprogrampage": map[string]interface{}{
			"title":          title,
			"appid":          appid,
			"pagepath":       pagePath,
			"thumb_media_id": mediaId,
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 创建请求
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", s.AccessToken)
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

	return ErrorHandler(ErrSendTextMessage, responseJson.ErrMsg, responseJson.Errcode)
}

// GetAccessToken 获取access_token
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
func (s *SDK) GetAccessToken() (*AccessTokenResponse, error) {
	// 接口地址
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		s.AppID, s.AppSecret)

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

// NewTemMessage 实例化模版消息
func (s *SDK) NewTemMessage(touser, templateID, url, appID, appPagePath, clientMsgID string, msgData map[string]string) *TempMessage {
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

// SendTempMessage 发送模版消息
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
func (s *SDK) SendTempMessage(message *TempMessage) error {
	// 将消息数据序列化为JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// 创建请求
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", s.AccessToken)
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

// GetUserList 获取用户列表
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/User_Management/Getting_a_User_List.html
func (s *SDK) GetUserList(nextOpenID string) (*GetUserListResponse, error) {
	// 接口地址
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/get?access_token=%s&next_openid=%s",
		s.AccessToken, nextOpenID)

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

// GetUserInfo 获取用户基本信息
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
func (s *SDK) GetUserInfo(openID string) (*GetUserInfoResponse, error) {
	// 接口地址
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN",
		s.AccessToken, openID)

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

// GetWebAuthAccessToken 获取网页授权access_token
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html#1
// 说明：此功能需要的权限较高，需要在微信公众号后台配置相关信息使用，详细使用方法流程请参考官方文档
func (s *SDK) GetWebAuthAccessToken(code string) (*GetWebAuthAccessTokenResponse, error) {
	// 接口地址
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		s.AppID, s.AppSecret, code)

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

// DownloadAmrVoiceByMediaID 通过获取临时素材接口下载amr格式音频到指定路径
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_temporary_materials.html
func (s *SDK) DownloadAmrVoiceByMediaID(mediaID, path string) error {
	// 接口地址
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s",
		s.AccessToken, mediaID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// DownloadAmrVoiceByMediaIDAndReturnBase64 通过获取临时素材接口下载amr格式音频，并返回Base64编码字符串
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_temporary_materials.html
func (s *SDK) DownloadAmrVoiceByMediaIDAndReturnBase64(mediaID string) (string, error) {
	// 接口地址
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s",
		s.AccessToken, mediaID)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	base64Str := base64.StdEncoding.EncodeToString(data)
	return base64Str, nil
}

// AddMaterial 新增永久素材
// 官方文档地址 https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
func (s *SDK) AddMaterial(mediaType, fileUrl string) (string, error) {

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/material/add_material?access_token=%s&type=%s",
		s.AccessToken, mediaType)

	resp, err := http.Get(fileUrl)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read downloaded file: %v", err)
	}

	// Create a buffer to hold the form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("media", "test.jpg")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	// 将文件字节流写入表单
	_, err = io.Copy(part, bytes.NewReader(fileBytes))
	if err != nil {
		return "", fmt.Errorf("failed to copy file content: %v", err)
	}

	// 关闭multipart写入器，最终生成整个表单数据
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// 创建HTTP POST请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	// 设置multipart表单的Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp2, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp2.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp2.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}
	var response AddMediaResponse
	if err = json.Unmarshal(respBody, &response); err != nil {
		return "", err
	}

	if response.Errcode != 0 {
		return "", fmt.Errorf(response.ErrMsg)
	}

	return response.MediaId, nil
}
