package platform

import "github.com/ScentWoman/Recorder/live"

// Room interface.
type Room interface {
	GetPlatformName() string
	GetInfo(*Opt) (live.Info, error)
	GetStreams(*Opt) ([]live.Stream, error)
}

// Opt is the interface of additional option.
type Opt interface{}
