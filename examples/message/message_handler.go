package main

import (
	"fmt"
	"github.com/supercat0867/wechat"
	"io/ioutil"
	"log"
	"net/http"
)

// 被动回复消息示例，实际使用中要在公众号后台配置号服务器地址等信息
func main() {
	sdk := wechat.NewMessageSDK()

	// 注册文本消息处理函数
	sdk.RegisterHandler(wechat.TextMessage, func(msg *wechat.Message, w http.ResponseWriter) {
		log.Printf("收到文本消息：%s\n,发送方openid：%s ", msg.Content, msg.FromUserName)
		// TODO 检查是否可
		responseXML := sdk.BuildTextResponse(msg.FromUserName, msg.ToUserName, "这是一条文本消息回复")
		fmt.Fprint(w, responseXML)
	})
	// 注册语音消息处理函数
	sdk.RegisterHandler(wechat.VoiceMessage, func(msg *wechat.Message, w http.ResponseWriter) {
		log.Printf("收到语音消息：%s\n,发送方openid：%s ", msg.Content, msg.FromUserName)
		fmt.Fprint(w, "success")
	})

	http.HandleFunc("/msgHandler", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				// 发生错误时的处理
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			sdk.HandleWeChatMessage(body, w)

		}
	})

	fmt.Println("启动服务器，监听/msgHandler路径")
	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
