package main

import (
	"fmt"
	"github.com/supercat0867/wechat"
)

func main() {
	sdk := wechat.NewMessageSDK()
	resp, err := sdk.GetAccessToken("", "")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.AccessToken)
}
