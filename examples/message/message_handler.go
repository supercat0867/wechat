package main

import (
	"fmt"
	"github.com/supercat0867/wechat"
	"io"
	"log"
	"net/http"
)

func main() {
	sdk := wechat.New("", "")

	// 注册文本消息处理函数
	sdk.RegisterHandler(wechat.TextMessage, func(msg *wechat.Message, w http.ResponseWriter) {
		log.Printf("收到文本消息：%s\n,发送方openid：%s ", msg.Content, msg.FromUserName)
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
			body, err := io.ReadAll(r.Body)
			if err != nil {
				// 发生错误时的处理
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			sdk.HandleWeChatMessage(body, w)

		}
	})

	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
