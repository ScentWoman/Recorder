package api

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestApiParse(t *testing.T) {
	resp, e := http.Get("https://www.huya.com/290429")
	// resp, e := http.Get("https://www.huya.com/xiaohesb")

	if e != nil {
		fmt.Println(e)
		t.Fail()
	}
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}
	page := html.UnescapeString(string(body))
	meta := page

	// TT_ROOM_DATA
	roomData := regexp.MustCompile("TT_ROOM_DATA.*?};var ").FindString(meta)
	roomData = roomData[strings.Index(roomData, "{") : len(roomData)-5]
	var roomDataAPI TtRoomData
	if e = json.Unmarshal([]byte(roomData), &roomDataAPI); e != nil {
		fmt.Println(e)
		t.FailNow()
	}
	if x, e := json.MarshalIndent(roomDataAPI, "", "  "); e != nil {
		fmt.Println(e)
		t.FailNow()
	} else {
		fmt.Println(string(x))
	}

	// stream
	stream := regexp.MustCompile("\"stream\".*?};").FindString(meta)
	stream = stream[strings.Index(stream, "{") : len(stream)-2]
	stream = stream[:strings.LastIndex(stream, "}")+1]

	var streamAPI Stream
	if e = json.Unmarshal([]byte(stream), &streamAPI); e != nil {
		fmt.Println(e)
		t.FailNow()
	}
	if x, e := json.MarshalIndent(streamAPI, "", "  "); e != nil {
		fmt.Println(e)
		t.FailNow()
	} else {
		fmt.Println(string(x))
	}

	t.Fail()
}
