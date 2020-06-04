package parse

import (
	"encoding/json"
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
	Contains *Contain                 `json:"contains"`
	Deletes  []string                 `json:"deletes"`
	Replaces []*Replace               `json:"replaces"`
}

type HasAttr struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type Contain struct {
	Key      string   `json:"key"`
	HasClass string   `json:"has_class"`
	Finds    []string `json:"finds"`
	HasAttr  *HasAttr `json:"has_attr"`
}

type Split struct {
	Key    string `json:"key"`
	Index  int    `json:"index"`
	Enable bool   `json:"enable"`
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
		if v.Type == "contains_list" {
			res[k] = containsList(dom.Selection, v)
		} else {
			res[k] = content(dom.Selection, v)
		}

	}

	return res, err
}

func containsList(dom *goquery.Selection, params *FilterParams) (ins interface{}) {
	s := dom.Clone()
	text := ""

	if params.Selector != "" {
		s.Find(params.Selector)
	}

	s = finds(params.Finds, s)
	s.Each(func(i int, ss *goquery.Selection) {
		if len(params.Contains.Finds) > 0 {
			ss = finds(params.Contains.Finds, ss)
		}
		if params.Contains.HasClass != "" && params.Contains.Key == "" && params.Contains.HasAttr == nil{
			flag := ss.HasClass(params.Contains.HasClass)
			if flag {
				text = ss.Text()
			}
		} else if params.Contains.Key != "" && params.Contains.HasClass == "" && params.Contains.HasAttr == nil {
			if strings.Contains(ss.Text(), params.Contains.Key) {
				text = ss.Text()
			}
		} else if params.Contains.HasAttr != nil && params.Contains.HasClass == "" && params.Contains.Key == "" {
			val, ok := s.Attr(params.Contains.HasAttr.Key)
			if ok {
				if val == params.Contains.HasAttr.Val {
					text = s.Text()
				}
			}
		} else if params.Contains.Key != "" && params.Contains.HasClass != "" && params.Contains.HasAttr == nil {
			flag := ss.HasClass(params.Contains.HasClass)
			if flag {
				if strings.Contains(ss.Text(), params.Contains.Key) {
					text = ss.Text()
				}
			}
		}
	})

	text = splitDeletesReplace(text, params)

	return text
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
	s := dom.Clone()
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

	text = splitDeletesReplace(text, params)
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
		s = s.Eq(params.Eq - 1)
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

func splitDeletesReplace(text string, params *FilterParams) string {
	if params.Split != nil {
		if params.Split.Enable {
			if params.Split.Key == "\\n" {
				text = strings.Split(text, "\n")[params.Split.Index]
			} else if params.Split.Key == "\\t" {
				text = strings.Split(text, "\t")[params.Split.Index]
			} else {
				text = strings.Split(text, params.Split.Key)[params.Split.Index]
			}
		} else {
			dataBytes := make([]byte, 0)
			if params.Split.Key == "\\n" {
				dataBytes, _ = json.Marshal(strings.Split(text, "\n"))
				text = string(dataBytes)
			} else if params.Split.Key == "\\t" {
				dataBytes, _ = json.Marshal(strings.Split(text, "\t"))

			} else {
				dataBytes, _ = json.Marshal(strings.Split(text, params.Split.Key))
			}
			text = string(dataBytes)
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
