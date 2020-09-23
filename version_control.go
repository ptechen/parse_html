package parse_html

import (
	"regexp"
	"strings"
)

type VersionControl struct {
	Rule   string                   `json:"rule" yaml:"rule"`
	Fields map[string]*FilterParams `json:"fields" yaml:"fields"`
}

type RuleFields struct {
	Configs []*VersionControl `json:"configs" yaml:"configs"`
}

func ParseHtmlVersion(html string, params []*VersionControl) (res map[string]interface{}, err error) {
	html = strings.ReplaceAll(html, "\n", "")
	for i := 0; i < len(params); i++ {
		reg := regexp.MustCompile(params[i].Rule)
		if reg == nil {
			continue
		}
		result := reg.FindAllStringSubmatch(html, -1)
		if result == nil {
			continue
		}
		res, err = ParseHtml(html, params[i].Fields)
		return res, err
	}
	return res, err
}