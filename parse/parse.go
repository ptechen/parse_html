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
	Attr     string                   `json:"attr"`
	Split    *Split                   `json:"split"`
	Contains []string                 `json:"contains"`
	Deletes  []string                 `json:"deletes"`
	Replaces []*Replace               `json:"replaces"`
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
		res[k] = content(dom.Selection, v)
	}

	return res, err
}

func content(dom *goquery.Selection, params *FilterParams) (ins interface{}) {
	s := dom
	text := ""
	if len(params.Finds) > 0 {
		for i := 0; i < len(params.Finds); i++ {
			find := params.Finds[i]
			s = s.Find(find)
		}
	}

	if params.Selector != "" {
		s.Find(params.Selector)
	}

	if params.Type == "list" {
		resList := make([]interface{}, 0, 10)
		s.Each(func(i int, ss *goquery.Selection) {
			res := make(map[string]interface{})
			if params.Keys != nil {
				for k, v := range params.Keys {
					res[k] = content(ss, v)
				}
				resList = append(resList, res)
			} else {
				r := content(ss, &FilterParams{
					Deletes: params.Deletes,
					Replaces: params.Replaces,
				})
				resList = append(resList, r)
			}

		})
		return resList
	}

	s = lastFirstEq(s, params)

	text = getText(s, params)

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
			if curDelete == "\\n" {
				text = strings.ReplaceAll(text, "\n", "")
			} else if curDelete == "\\t" {
				text = strings.ReplaceAll(text, "\t", "")
			}  else {
				text = strings.ReplaceAll(text, curDelete, "")
			}

		}
	}

	if len(params.Replaces) > 0 {
		for i := 0; i < len(params.Replaces); i++ {
			rep := params.Replaces[i]
			text = strings.ReplaceAll(text, rep.Before, rep.After)
		}
	}
	return text
}

func lastFirstEq(s *goquery.Selection, params *FilterParams) *goquery.Selection {
	if params.Last {
		s = s.Last()
	}

	if params.First {
		s = s.First()
	}

	if params.Eq !=0 {
		s = s.Eq(params.Eq)
	}
	return s
}

func getText(s *goquery.Selection, params *FilterParams) (text string)  {
	if params.Attr != "" {
		ok := false
		text, ok = s.Attr(params.Attr)
		if !ok {
			return ""
		}
	} else {
		text = s.Text()
	}
	return text
}