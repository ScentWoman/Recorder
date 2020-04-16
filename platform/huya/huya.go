package huya

import (
	"encoding/json"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ScentWoman/Recorder/live"
	"github.com/ScentWoman/Recorder/platform/huya/api"
)

var (
	ttRoomDataRegexp    = regexp.MustCompile("TT_ROOM_DATA = {.*?};var")
	streamRegexp        = regexp.MustCompile("\"stream\".*?};")
	ttProfileInfoRegexp = regexp.MustCompile("TT_PROFILE_INFO = {.*?};var")

	_UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36 Edg/80.0.361.111"
)

// Huya struct
type Huya struct {
	URL string
}

// GetPlatformName implements (Room).GetPlatformName
func (h *Huya) GetPlatformName() string {
	return "虎牙直播"
}

// GetInfo implements (Room).GetInfo
func (h *Huya) GetInfo(opt interface{}) (l *live.Info, e error) {
	defer func() {
		if p := recover(); p != nil {
			e = p.(error)
		}
	}()

	req, _ := http.NewRequest("GET", h.URL, nil)
	req.Header.Add("User-Agent", _UA)
	body, e := live.PageBody(req)
	if e != nil {
		return
	}
	roomData := ttRoomDataRegexp.Find(body)
	var room api.TtRoomData
	if e = json.Unmarshal(shave(roomData), &room); e != nil {
		return
	}
	profileData := ttProfileInfoRegexp.Find(body)
	var profile api.TtProfileInfo
	if e = json.Unmarshal(shave(profileData), &profile); e != nil {
		return
	}

	l = &live.Info{
		URL: h.URL,
		Host: live.Host{
			Name:   profile.Nick,
			Nick:   profile.Nick,
			Avatar: profile.Avatar,
			// UID:
		},
		Title:       room.Introduction,
		Catagory:    room.GameHostName,
		ContentName: room.GameFullName,

		IsLive:    room.IsOn,
		StartTime: itoTime(room.StartTime),
	}

	return
}

// GetStreams implements (Room).GetStreams
func (h *Huya) GetStreams(opt interface{}) (s []live.Stream, e error) {
	defer func() {
		if p := recover(); p != nil {
			e = p.(error)
		}
	}()

	req, _ := http.NewRequest("GET", h.URL, nil)
	req.Header.Add("User-Agent", _UA)
	body, e := live.PageBody(req)
	if e != nil {
		return
	}
	streamData := streamRegexp.Find(body)
	// dangerous
	streamData = shave(streamData[:len(streamData)-3])
	var stream api.Stream
	if e = json.Unmarshal(streamData, &stream); e != nil {
		return
	}

	s = make([]live.Stream, len(stream.Data[0].GameStreamInfoList))
	for k, v := range stream.Data[0].GameStreamInfoList {
		sb := strings.Builder{}
		_, _ = sb.WriteString(html.UnescapeString(v.SFlvURL))
		_, _ = sb.WriteString("/")
		_, _ = sb.WriteString(html.UnescapeString(v.SStreamName))
		_, _ = sb.WriteString(".")
		_, _ = sb.WriteString(html.UnescapeString(v.SFlvURLSuffix))
		_, _ = sb.WriteString("?t=100&sv=1910112100&")
		_, _ = sb.WriteString(html.UnescapeString(v.SFlvAntiCode))
		surl := sb.String()
		req, e := http.NewRequest("GET", surl, nil)
		if e != nil {
			return nil, e
		}
		s[k] = live.Stream{
			Req:         req,
			URL:         surl,
			Bitrate:     0,
			Suffix:      v.SFlvURLSuffix,
			Description: v.SCdnType,
		}
	}
	return
}

func shave(b []byte) []byte {
	for b[0] != '{' {
		b = b[1:]
	}
	for b[len(b)-1] != '}' {
		b = b[:len(b)-1]
	}
	return b
}

func itoTime(i interface{}) time.Time {
	var ts int64
	switch it := i.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		ts = toint64(it)
	case string:
		ts, _ = strconv.ParseInt(it, 10, 64)
	}

	return time.Unix(ts, 0)
}

func toint64(i interface{}) int64 {
	switch ii := i.(type) {
	case int:
		return int64(ii)
	case int8:
		return int64(ii)
	case int16:
		return int64(ii)
	case int32:
		return int64(ii)
	case int64:
		return int64(ii)
	case uint:
		return int64(ii)
	case uint8:
		return int64(ii)
	case uint16:
		return int64(ii)
	case uint32:
		return int64(ii)
	case uint64:
		return int64(ii)
	default:
		return 0
	}
}
