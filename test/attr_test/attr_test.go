package attr_test

import (
	"fmt"
	"github.com/ptechen/config"
	parse "github.com/ptechen/parse_html"
	"io/ioutil"
	"testing"
)

func TestAttr(t *testing.T) {
	yml := map[string]*parse.FilterParams{}
	conf := &config.Config{}
	conf.YAML("test.yml", &yml)
	htmlBytes, err := ioutil.ReadFile("test.html")
	if err != nil {
		t.Error(err)
	}
	html := string(htmlBytes)
	res, err := parse.ParseHtml(html, yml)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
