package parse

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type FilterParams struct {
	Selector string                   `json:"selector"`
	Finds    []string                 `json:"finds"`
	Type     string                   `json:"type"`
	SubFinds []string                 `json:"sub_finds"`
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

func finds(finds []string, s *goquery.Selection) *goquery.Selection {
	if len(finds) > 0 {
		for i := 0; i < len(finds); i++ {
			find := finds[i]
			s = s.Find(find)
		}
	}
	return s
}

func content(dom *goquery.Selection, params *FilterParams) (ins interface{}) {
	s := dom
	text := ""

	if params.Selector != "" {
		s.Find(params.Selector)
	}

	s = finds(params.Finds, s)

	if params.Type == "list" {
		resList := make([]interface{}, 0, 10)
		s.Each(func(i int, ss *goquery.Selection) {
			ss = finds(params.SubFinds, ss)
			res := make(map[string]interface{})
			if params.Keys != nil {
				for k, v := range params.Keys {
					if params.HasClass != "" {
						hasClass := ss.HasClass(params.HasClass)
						if hasClass {
							res[k] = content(ss, v)
						}
					} else {
						res[k] = content(ss, v)
					}
				}
				resList = append(resList, res)
			} else {

				if params.HasClass != "" {
					hasClass := ss.HasClass(params.HasClass)
					if hasClass {
						r := content(ss, &FilterParams{
							Deletes:  params.Deletes,
							Replaces: params.Replaces,
						})

						resList = append(resList, r)
					}
				} else {
					r := content(ss, &FilterParams{
						Deletes:  params.Deletes,
						Replaces: params.Replaces,
					})

					resList = append(resList, r)
				}

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
			} else {
				text = strings.ReplaceAll(text, curDelete, "")
			}
		}
	}

	if len(params.Replaces) > 0 {
		for i := 0; i < len(params.Replaces); i++ {
			rep := params.Replaces[i]
			if rep.Before == "\\n" {
				text = strings.ReplaceAll(text, "\n", rep.After)
			} else if rep.Before == "\\t" {
				text = strings.ReplaceAll(text, "\t", rep.After)
			} else {
				text = strings.ReplaceAll(text, rep.Before, rep.After)
			}
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

	if params.Eq != 0 {
		s = s.Eq(params.Eq)
	}
	return s
}

func getText(s *goquery.Selection, params *FilterParams) (text string) {
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
