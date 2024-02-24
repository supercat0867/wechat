package main

import (
	"fmt"
	"github.com/supercat0867/wechat"
)

func main() {
	// 获取access_token
	resp, err := wechat.GetAccessToken("", "")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.AccessToken)
}
