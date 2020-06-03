package parse_html

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type FilterParams struct {
	Selector string     `json:"selector"`
	Attr     string     `json:"attr"`
	Split    *Split      `json:"split"`
	Contains []string   `json:"contains"`
	Deletes  []string   `json:"deletes"`
	Replaces []*Replace `json:"replaces"`
}

type Split struct {
	Key   string `json:"key"`
	Index int    `json:"index"`
}

type Replace struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

func ParseHtml(html string, params map[string]*FilterParams) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return res, err
	}
	for k, v := range params {
		res[k] = content(dom, v)
	}

	return res, err
}

func content(dom *goquery.Document, params *FilterParams) (text string) {
	if params.Attr == "" {
		text = dom.Find(params.Selector).Text()
	} else {
		ok := false
		text, ok = dom.Find(params.Selector).Attr(params.Attr)
		if !ok {
			return ""
		}
	}
	if params.Split != nil {
		text = strings.Split(text, params.Split.Key)[params.Split.Index]
	}

	if len(params.Contains) > 0 {
		for i := 0; i < len(params.Contains); i++ {
			contain := params.Contains[i]
			if !strings.Contains(text, contain) {
				return ""
			}
		}
	}

	if len(params.Deletes) > 0 {
		for i := 0; i < len(params.Deletes); i++ {
			curDelete := params.Deletes[i]
			strings.ReplaceAll(text, curDelete, "")
		}
	}

	if len(params.Replaces) > 0 {
		for i := 0; i < len(params.Replaces); i++ {
			rep := params.Replaces[i]
			strings.ReplaceAll(text, rep.Before, rep.After)
		}
	}
	return text
}