package parse

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type FilterParams struct {
	Selector string                   `json:"selector"`
	Finds    []string                 `json:"finds"`
	Type     string                   `json:"type"`
	Keys     map[string]*FilterParams `json:"keys"`
	Last     bool                     `json:"last"`
	First    bool                     `json:"first"`
	Eq       int                      `json:"eq"`
	HasClass string                   `json:"has_class"`
	Attr     string                   `json:"attr"`
	Split    *Split                   `json:"split"`
	Contains []string                 `json:"contains"`
	Deletes  []string                 `json:"deletes"`
	Replaces []*Replace               `json:"replaces"`
}

type Tag struct {
	Index   int    `json:"index"`
	Contain string `json:"contain"`
}

type Split struct {
	Keys []string
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
		if k == "manual" {

		} else {
			res[k] = content(dom.Selection, v)
		}

	}

	return res, err
}

func finds(s *goquery.Selection, finds []string) *goquery.Selection {
	if len(finds) > 0 {
		for i := 0; i < len(finds); i++ {
			find := finds[i]
			s = s.Find(find)
		}
	}
	return s
}

func selection(s *goquery.Selection, selector string) *goquery.Selection {
	if selector != "" {
		s.Find(selector)
	}
	return s
}

func content(dom *goquery.Selection, params *FilterParams, index... int) (ins interface{}) {
	s := dom
	text := ""

	s = selection(s, params.Selector)

	s = finds(s, params.Finds)

	if params.Type == "list" {
		resList := make([]interface{}, 0, 10)
		s.Each(func(i int, ss *goquery.Selection) {
			res := make(map[string]interface{})
			if params.Keys != nil {
				for k, v := range params.Keys {
					res[k] = content(ss, v)
				}
				resList = append(resList, res, i)
			} else {

				r := content(ss, &FilterParams{
					Deletes:  params.Deletes,
					Replaces: params.Replaces,
				})

				resList = append(resList, r, i)

			}
		})
		return resList
	}

	s = lastFirstEq(s, params)

	if !s.HasClass(params.HasClass) {
		return ""
	}

	text = getText(s, params)

	text = clear(text, params)
	return text
}
