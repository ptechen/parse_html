package parse

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseHtml(t *testing.T) {
	params := make(map[string]*FilterParams)
	params["val"] = &FilterParams{
		Selector: "#job-view-enterprise > div.wrap.clearfix > div.clearfix > div.main > div.about-position > div:nth-child(2) > div.clearfix > div.job-title-left > p.job-item-title",
		Split: &Split{
			Key:   "\n",
			Index: 0,
		},
	}
	dataBytes, err := ioutil.ReadFile("index.html")
	if err != nil {
		t.Error(err)
	}
	dataStr := string(dataBytes)
	res, err := ParseHtml(dataStr, params)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("%#v", res))
}