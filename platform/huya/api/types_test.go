package api

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
)

func TestApiParse(t *testing.T) {
	// resp, e := http.Get("https://www.huya.com/290429")
	resp, e := http.Get("https://www.huya.com/xiaohesb")

	if e != nil {
		fmt.Println(e)
		t.Fail()
	}
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		fmt.Println(e)
		t.Fail()
	}

	// TT_ROOM_DATA
	roomData := regexp.MustCompile("TT_ROOM_DATA = {.*?};var").Find(body)
	for roomData[0] != '{' {
		roomData = roomData[1:]
	}
	roomData = roomData[:len(roomData)-4]
	var roomDataAPI TtRoomData
	if e = json.Unmarshal(roomData, &roomDataAPI); e != nil {
		fmt.Println(e)
		t.FailNow()
	}
	if x, e := json.MarshalIndent(roomDataAPI, "", "  "); e != nil {
		fmt.Println(e)
		t.FailNow()
	} else {
		fmt.Println(string(x))
	}

	// TtProfileInfo
	profileData := regexp.MustCompile("TT_PROFILE_INFO = {.*?};var").Find(body)
	for profileData[0] != '{' {
		profileData = profileData[1:]
	}
	profileData = profileData[:len(profileData)-4]
	// fmt.Println(string(profileData))
	var profile TtProfileInfo
	if e = json.Unmarshal(profileData, &profile); e != nil {
		fmt.Println(e)
		t.FailNow()
	}
	if x, e := json.MarshalIndent(profile, "", "  "); e != nil {
		fmt.Println(e)
		t.FailNow()
	} else {
		fmt.Println(string(x))
	}

	// stream
	stream := regexp.MustCompile("\"stream\".*?};").Find(body)
	if len(stream) < 30 {
		fmt.Println("NO STREAM")
		t.FailNow()
		return
	}
	stream = stream[:len(stream)-2]
	for stream[0] != '{' {
		stream = stream[1:]
	}
	for stream[len(stream)-1] != '}' {
		stream = stream[:len(stream)-1]
	}

	var streamAPI Stream
	if e = json.Unmarshal(stream, &streamAPI); e != nil {
		fmt.Println(e)
		t.FailNow()
	}
	if x, e := json.MarshalIndent(streamAPI, "", "  "); e != nil {
		fmt.Println(e)
		t.FailNow()
	} else {
		fmt.Println(string(x))
	}
	fmt.Println(html.UnescapeString(streamAPI.Data[0].GameStreamInfoList[0].SFlvAntiCode))

	t.Fail()
}
