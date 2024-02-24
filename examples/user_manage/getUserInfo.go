package main

import (
	"fmt"
	"github.com/supercat0867/wechat"
)

// 示例 获取用户基础信息
func main() {
	// 获取access_token
	resp, err := wechat.GetAccessToken("", "")
	if err != nil {
		panic(err)
	}
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
}
