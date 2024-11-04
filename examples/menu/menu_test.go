package menu

import (
	"github.com/supercat0867/wechat"
	"testing"
)

// 创建自定义菜单
func TestCreateMenu(t *testing.T) {
	sdk := wechat.New("", "")
	menu := wechat.Menu{
		Button: []wechat.MenuButton{
			{
				Type:     "miniprogram",
				Name:     "小程序",
				Url:      "http://mp.weixin.qq.com",
				AppID:    "",
				PagePath: "pages/index/index",
			},
		},
	}
	if err := sdk.CreateMenu(menu); err != nil {
		t.Error(err)
		return
	}
	return
}
