package wechat

import "fmt"

// ErrorHandler 错误处理
func ErrorHandler(action, errmsg string, errcode int) error {
	return fmt.Errorf("%s:%s,错误码：%d", action, errmsg, errcode)
}

var (
	ErrGetAccessToken        = "access_token获取失败"
	ErrSendTempMessage       = "模版消息发送失败"
	ErrGetUserList           = "用户列表获取失败"
	ErrGetUserInfo           = "用户基础信息获取失败"
	ErrGetWebAuthAccessToken = "网页授权access_token获取失败"
)
