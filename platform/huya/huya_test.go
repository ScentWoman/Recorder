package huya

import (
	"fmt"
	"testing"
)

func TestHuya(t *testing.T) {
	testURL("https://www.huya.com/290429", t)
	testURL("https://www.huya.com/xiaohesb", t)
	t.Fail()
}

func testURL(u string, t *testing.T) {
	h := &Huya{
		URL: u,
	}
	if h.GetPlatformName() != "虎牙直播" {
		fmt.Println(h.GetPlatformName(), "!=", "虎牙直播")
		t.FailNow() // LOL
	}
	if info, e := h.GetInfo(nil); e != nil {
		fmt.Println(e)
		t.FailNow()
	} else {
		if info.IsLive {
			fmt.Println(info.Host.Name, "正在直播:", info.Title)
			// fmt.Printf("%#v\n", info)
			fmt.Println("Start:", info.StartTime)
			fmt.Println(info.Catagory, "->", info.ContentName)

			streams, e := h.GetStreams(nil)
			if e != nil {
				fmt.Println(e)
				t.FailNow()
			}
			for k, v := range streams {
				fmt.Println("  ", k, ":", v.URL)
			}
		} else {
			fmt.Println(info.Host.Name, "未在直播！")
		}
	}
}
