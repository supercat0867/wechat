package main

import (
	"fmt"
	"github.com/supercat0867/wechat"
)

// 示例 获取用户列表
func main() {
	sdk := wechat.NewMessageSDK()
	// 获取access_token
	resp, err := sdk.GetAccessToken("", "")
	if err != nil {
		panic(err)
	}
	// 获取用户列表
	userList, err := sdk.GetUserList(resp.AccessToken, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(userList.Data.OpenID)
}
