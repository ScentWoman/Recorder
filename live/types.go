package live

import (
	"io/ioutil"
	"net/http"
	"time"
)

// Stream illustrates a live stream.
type Stream struct {
	Req         *http.Request
	URL         string
	Bitrate     int
	Suffix      string
	Description string
}

// Info contains info of the live room.
type Info struct {
	URL                   string
	Host                  Host
	Title                 string
	Catagory, ContentName string

	IsLive    bool
	StartTime time.Time
}

// Host illustrates an host.
type Host struct {
	Name, Nick string
	Avatar     string
	UID        string
}

var (
	client = http.Client{}
)

// PageBody fetches a page body.
func PageBody(req *http.Request) (body []byte, e error) {
	resp, e := client.Do(req)
	if e != nil {
		return
	}
	body, e = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}
