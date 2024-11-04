package examples

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

// 微信公众号设置的token
const token = "YourWeChatToken"

type WeChatMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        int64    `xml:"MsgId"`
}

// 验证签名
func validate(r *http.Request) bool {
	signature := r.URL.Query().Get("signature")
	timestamp := r.URL.Query().Get("timestamp")
	nonce := r.URL.Query().Get("nonce")

	tmpStrs := sort.StringSlice{token, timestamp, nonce}
	sort.Strings(tmpStrs)
	tmpStr := strings.Join(tmpStrs, "")

	sha1 := sha1.New()
	io.WriteString(sha1, tmpStr)
	hashcode := fmt.Sprintf("%x", sha1.Sum(nil))

	return hashcode == signature
}

// 处理微信服务器发送的GET请求（服务器验证）
func handleGet(w http.ResponseWriter, r *http.Request) {
	if validate(r) {
		echoStr := r.URL.Query().Get("echostr")
		fmt.Fprintf(w, echoStr)
	} else {
		fmt.Fprintf(w, "Failed to validate")
	}
}

// 处理微信服务器发送的POST请求（接收消息）
func handlePost(w http.ResponseWriter, r *http.Request) {
	if !validate(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't read body", http.StatusBadRequest)
		return
	}

	var msg WeChatMessage
	if err := xml.Unmarshal(body, &msg); err != nil {
		http.Error(w, "Error parsing XML", http.StatusBadRequest)
		return
	}

	// 根据消息类型处理消息，这里简单回显收到的文本消息
	if msg.MsgType == "text" {
		response := fmt.Sprintf(`<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[You said: %s]]></Content>
</xml>`, msg.FromUserName, msg.ToUserName, msg.CreateTime, msg.Content)

		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintf(w, response)
	}
}

func main() {
	http.HandleFunc("/msgHandler", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGet(w, r)
		case http.MethodPost:
			handlePost(w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Println("Server error:", err)
	}
	//if err := http.ListenAndServeTLS(":443", "", "", nil); err != nil {
	//	fmt.Println("Server error:", err)
	//}
}
