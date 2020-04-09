package recorder

import (
	"fmt"
	"strings"
	"testing"
)

var (
	config = `Huya-Example:
  url: https://www.huya.com/290429
  interval: 2
  split: 3600
  save: /home/Record/save
  filename_format: "${Date:2006-01-02-15-04-05}_${Title}_${Plt}"
  on_finish: "script dir\\ with\\ space ${1} ${Date:01-02}"`
)

func TestParse(t *testing.T) {
	cm, e := Parse(strings.NewReader(config))
	if e != nil {
		fmt.Println(e)
		t.FailNow()
	}
	for k, v := range cm {
		fmt.Println(k)
		fmt.Printf("%#v\n", v)
		for _, v := range strings.Split(v.OnFinish, " ") {
			fmt.Println("--", v)
		}
	}

	t.Fail()
}
