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
	// 构造消息结构
	data := map[string]string{
		"thing2":   "请假流程通知",
		"time15":   "2012-01-02",
		"phrase10": "小明",
		"thing16":  "扶老奶奶过马路",
	}
	tempMessage := wechat.NewTemMessage("obIt16lHlQiZpT5MYC_lTfFv7ZSA", "IWMM8w9XD3jqc01gXyisvG6Y6yPMfGhlGyLPWimAN2w",
		"www.baidu.com", "", "", "", data)
	// 发送模版消息
	err = tempMessage.Send(resp.AccessToken)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("模版消息发送成功！")
	}
}
