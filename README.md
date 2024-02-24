# wechat sdk for golang

本项目是一个为 Go 项目提供微信公众号相关功能的 SDK，是基于微信公众号官方文档开发。

## 功能

项目目前包含以下功能：

| 模块          | 功能                 | 方法                                                                                                                          |
|-------------|--------------------|-----------------------------------------------------------------------------------------------------------------------------|
| 用户管理        | 获取用户列表             | func GetUserList(accessToken, nextOpenID string) (*GetUserListResponse, error)                                              |
|             | 获取用户基础信息           | func GetUserInfo(accessToken, openID string) (*GetUserInfoResponse, error)                                                  |
| AccessToken | 获取公众号access_token  | func GetAccessToken(appID, appSecret string) (*AccessTokenResponse, error)                                                  |
| 消息          | 实例化模版消息            | func NewTemMessage(touser, templateID, url, appID, appPagePath, clientMsgID string, msgData map[string]string) *TempMessage |
|             | 发送模版消息             | func (m *TempMessage) Send(accessToken string) error                                                                        |
| 授权          | 获取网页授权access_token | func GetWebAuthAccessToken(appID, appSecret, code string) (*GetWebAuthAccessTokenResponse, error)                           |
| 持续更新中...    |                    |                                                                                                                             |

## 快速开始

要开始使用本项目，请确保你的 Go 项目已经初始化并且可以管理依赖（使用 Go Modules）。

1. 在你的 Go 项目中，引入本项目：
   ```bash
   go get github.com/supercat0867/wechat
2. 获取access_token
    ```bash
    // 获取access_token
   resp, err := wechat.GetAccessToken("", "")
   if err != nil {
      panic(err)
   }
   fmt.Println(resp.AccessToken)
3. 获取用户列表
    ```bash
   // 获取用户列表
   userList, err := wechat.GetUserList(resp.AccessToken, "")
   if err != nil {
     panic(err)
   }
   // 查询第一个用户的基础信息
   info, err := wechat.GetUserInfo(resp.AccessToken, userList.Data.OpenID[0])
   if err != nil {
     panic(err)
   }
   fmt.Println(info)
4. 更多功能参考源码，注释已详细标注。   