package user_manage

import (
	"github.com/supercat0867/wechat"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	sdk := wechat.New("", "")
	// 获取用户列表
	userList, err := sdk.GetUserList("")
	if err != nil {
		t.Error(err)
		return
	}
	// 查询第一个用户的基础信息
	info, err := sdk.GetUserInfo(userList.Data.OpenID[0])
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info)
	return
}
