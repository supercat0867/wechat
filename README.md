# wechat-officialAccount-sdk-go

本项目是一个为Go项目提供微信公众号相关功能的非官方SDK，基于微信公众号官方文档开发。

## 功能

项目目前包含以下功能：

| 模块          | 功能                 | 方法                                                                                                                                   |
|-------------|--------------------|--------------------------------------------------------------------------------------------------------------------------------------|
| 自定义菜单       | 创建自定义菜单            | func (s *SDK) CreateMenu(menu Menu) error                                                                                            |
| 用户管理        | 获取用户列表             | func (s *SDK) GetUserList(nextOpenID string) (*GetUserListResponse, error)                                                           |
|             | 获取用户基础信息           | func (s *SDK) GetUserInfo(openID string) (*GetUserInfoResponse, error)                                                               |
| AccessToken | 获取公众号access_token  | func (s *SDK) GetAccessToken() (*AccessTokenResponse, error)                                                                         |
| 模版消息        | 实例化模版消息            | func (s *SDK) NewTemMessage(touser, templateID, url, appID, appPagePath, clientMsgID string, msgData map[string]string) *TempMessage |
|             | 发送模版消息             | func (s *SDK) SendTempMessage(message *TempMessage) error                                                                            |
| 授权          | 获取网页授权access_token | func GetWebAuthAccessToken(code string) (*GetWebAuthAccessTokenResponse, error)                                                      |
| 客服消息        | 发送文本消息             | func (s *SDK)SendTextMessage(toUser, content string) error                                                                           |
|             | 发送小程序卡片消息          | func (s *SDK) SendMiniprogramMessage(toUser, title, appid, pagePath, mediaId string) error                                           |
| 素材管理        | 下载音频文件             | func (s *SDK) DownloadVoice(mediaID, path string) error                                                                              |
|             | 新增永久素材             | func (s *SDK) AddMaterial(mediaType, fileUrl string) (string, error)                                                                 |

## 快速开始

1. 引入本项目：
   ```bash
   go get github.com/supercat0867/wechat
2. 发送模版消息
    ```bash
   func main() {
       sdk := wechat.NewMessageSDK("", "")
       // 构造消息结构
       data := map[string]string{
           "thing2":   "请假流程通知",
           "time15":   "2012-01-02",
           "phrase10": "小明",
           "thing16":  "扶老奶奶过马路",
       }
       tempMessage := sdk.NewTemMessage("obIt16lHlQiZpT5MYC_lTfFv7ZSA", "IWMM8w9XD3jqc01gXyisvG6Y6yPMfGhlGyLPWimAN2w",
   "www.baidu.com", "", "", "", data)
       // 发送模版消息
       err := sdk.SendTempMessage(tempMessage)
       if err != nil {
         panic(err)
       } else {
         fmt.Println("模版消息发送成功！")
       }
   }

3. 更多功能参考源码示例...   