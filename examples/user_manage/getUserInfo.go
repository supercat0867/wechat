package main

import (
	"fmt"
	"github.com/supercat0867/wechat"
)

// 示例 获取用户基础信息
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
	// 查询第一个用户的基础信息
	info, err := sdk.GetUserInfo(resp.AccessToken, userList.Data.OpenID[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(info)
}
