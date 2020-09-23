package parse_html

import (
	"errors"
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
	curHtml := strings.ReplaceAll(html, "\n", "")
	for i := 0; i < len(params); i++ {
		reg := regexp.MustCompile(params[i].Rule)
		if reg == nil {
			continue
		}
		result := reg.FindAllStringSubmatch(curHtml, -1)
		if result == nil {
			if i == len(params) -1 {
				return res, errors.New("all rules were failed")
			}
			continue
		}
		res, err = ParseHtml(html, params[i].Fields)
		return res, err
	}
	return res, err
}