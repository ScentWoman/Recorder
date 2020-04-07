package live

import (
	"net/http"
	"time"
)

// Stream illustrates a live stream.
type Stream struct {
	Req         *http.Request
	URL         string
	Bitrate     string
	Suffix      string
	Description string
}

// Info contains info of the live room.
type Info struct {
	URL                   string
	Name                  string
	Host                  Host
	Title                 string
	Catagory, ContentName string

	IsLive    bool
	StartTime time.Time
	Duration  time.Duration
}

// Host illustrates an host.
type Host struct {
	Name, Nick string
	Avatar     string
	UID        string
}
