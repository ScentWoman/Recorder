package recorder

import (
	"io"

	"gopkg.in/yaml.v2"
)

// Config is config struct.
type Config struct {
	URL      string `yaml:"url"`
	Interval int    `yaml:"interval"`
	Split    int    `yaml:"split"`
	Save     string `yaml:"save"`
	NFormat  string `yaml:"filename_format"`
	OnFinish string `yaml:"on_finish"`
}

// Parse parses config via io.Reader.
func Parse(r io.Reader) (c map[string]Config, e error) {
	e = yaml.NewDecoder(r).Decode(&c)
	return
}
