package version_control_test

import (
	"fmt"
	"github.com/ptechen/config"
	parse "github.com/ptechen/parse_html"
	"io/ioutil"
	"testing"
)

func TestParseHtmlVersion(t *testing.T) {
	yml := &parse.RuleFields{}
	conf := &config.Config{}
	conf.YAML("test.yml", &yml)
	htmlBytes, err := ioutil.ReadFile("test.html")
	if err != nil {
		t.Error(err)
	}
	html := string(htmlBytes)
	res, err := parse.ParseHtmlVersion(html, yml.Configs)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}