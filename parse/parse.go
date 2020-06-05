package parse

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type FilterParams struct {
	Selector    string                   `json:"selector" yaml:"selector"`
	Finds       []string                 `json:"finds" yaml:"finds"`
	Type        string                   `json:"type" yaml:"type"`
	SubFinds    []string                 `json:"sub_finds" yaml:"sub_finds"`
	Keys        map[string]*FilterParams `json:"keys" yaml:"keys"`
	Last        bool                     `json:"last" yaml:"last"`
	First       bool                     `json:"first" yaml:"first"`
	Html        bool                     `json:"html" yaml:"html"`
	Eq          int                      `json:"eq" yaml:"eq"`
	HasClass    string                   `json:"has_class" yaml:"has_class"`
	Attr        string                   `json:"attr" yaml:"attr"`
	Split       *Split                   `json:"split" yaml:"split"`
	Contains    *Contain                 `json:"contains" yaml:"contains"`
	NotContains []string                 `json:"not_contains" yaml:"not_contains"`
	Deletes     []string                 `json:"deletes" yaml:"deletes"`
	Replaces    []*Replace               `json:"replaces" yaml:"replaces"`
}

type HasAttr struct {
	Key string `json:"key" yaml:"key"`
	Val string `json:"val" yaml:"val"`
}

type Contain struct {
	Key      string   `json:"key" yaml:"key"`
	HasClass string   `json:"has_class" yaml:"has_class"`
	Lable    *Lable   `json:"lable" yaml:"lable"`
	Finds    []string `json:"finds" yaml:"finds"`
	Eq       int      `json:"eq" yaml:"eq"`
	HasAttr  *HasAttr `json:"has_attr" yaml:"has_attr"`
}

type Lable struct {
	Finds    []string `json:"finds" yaml:"finds"`
	HasClass string   `json:"has_class" yaml:"has_class"`
	Contains []string `json:"contains" yaml:"contains"`
}

func ParseHtml(html string, params map[string]*FilterParams) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return res, err
	}
	for k, v := range params {
		if v.Type == "contains_list" {
			res[k] = v.containsList(dom.Selection)
		} else {
			res[k] = v.content(dom.Selection)
		}
	}
	return res, err
}

func (params *FilterParams) containsList(dom *goquery.Selection) (ins interface{}) {
	s := dom.Clone()
	text := ""

	if params.Selector != "" {
		s.Find(params.Selector)
	}

	s = finds(params.Finds, s)
	s.Each(func(i int, ss *goquery.Selection) {
		lableSelector := ss.Clone()
		if len(params.Contains.Finds) > 0 {
			ss = finds(params.Contains.Finds, ss)
		}
		if params.Contains.Eq != 0 {
			ss = ss.Eq(params.Contains.Eq - 1)
		}
		ok := lableHasClass(lableSelector, params.Contains.Lable)
		if params.Contains.HasClass != "" && params.Contains.Key == "" && params.Contains.HasAttr == nil {
			flag := ss.HasClass(params.Contains.HasClass)
			if flag {
				if ok {
					text = ss.Text()
				}
			}
		} else if params.Contains.Key != "" && params.Contains.HasClass == "" && params.Contains.HasAttr == nil {
			if strings.Contains(ss.Text(), params.Contains.Key) {
				if ok {
					text = ss.Text()
				}
			}
		} else if params.Contains.HasAttr != nil && params.Contains.HasClass == "" && params.Contains.Key == "" {
			val, ok := ss.Attr(params.Contains.HasAttr.Key)
			if ok {
				if val == params.Contains.HasAttr.Val {
					if ok {
						text = ss.Text()
					}
				}
			}
		} else if params.Contains.Key != "" && params.Contains.HasClass != "" && params.Contains.HasAttr == nil {
			flag := ss.HasClass(params.Contains.HasClass)
			if flag {
				if strings.Contains(ss.Text(), params.Contains.Key) {
					if ok {
						text = ss.Text()
					}
				}
			}
		} else if params.Contains.Key == "" && params.Contains.HasClass == "" && params.Contains.HasAttr == nil {
			if ok {
				text = ss.Text()
			}
		}
	})

	text = params.splitDeletesReplace(text)

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

func (params *FilterParams) content(dom *goquery.Selection) (ins interface{}) {
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
							res[k] = v.content(ss)
						}
					} else {
						res[k] = v.content(ss)
					}
				}
				resList = append(resList, res)
			} else {

				if params.HasClass != "" {
					hasClass := ss.HasClass(params.HasClass)
					if hasClass {
						cur := &FilterParams{
							Deletes:  params.Deletes,
							Replaces: params.Replaces,
						}
						r := cur.content(ss)

						resList = append(resList, r)
					}
				} else {
					cur := &FilterParams{
						Deletes:  params.Deletes,
						Replaces: params.Replaces,
					}
					r := cur.content(ss)

					resList = append(resList, r)
				}

			}

		})
		return resList
	}

	s = params.lastFirstEq(s)
	text = params.getText(s)
	text = params.notContains(text)
	text = params.splitDeletesReplace(text)
	return text
}

func lableHasClass(s *goquery.Selection, params *Lable) bool {
	s = s.Clone()
	flag := true
	if params != nil {
		s = finds(params.Finds, s)
		if params.HasClass != "" {
			if finds(params.Finds, s).HasClass(params.HasClass) {
				flag = true
			} else {
				return false
			}
		}
		if len(params.Contains) > 0{
			for i := 0; i < len(params.Contains); i ++ {
				if !strings.Contains(s.Text(), params.Contains[i]) {
					return false
				}
			}
		}
	}
	return flag
}

func (params *FilterParams) notContains(text string) string {
	if len(params.NotContains) > 0 {
		length := len(params.NotContains)
		for i := 0; i < length; i++ {
			cur := params.NotContains[i]
			if !strings.Contains(text, cur) {
				return ""
			}
		}
	}
	return text
}

func (params *FilterParams) lastFirstEq(s *goquery.Selection) *goquery.Selection {
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

func (params *FilterParams) getText(s *goquery.Selection) (text string) {
	if params.Attr != "" {
		ok := false
		text, ok = s.Attr(params.Attr)
		if !ok {
			return ""
		}
	} else if params.Html {
		text, _ = s.Html()
	} else {
		text = s.Text()
	}
	return text
}

func (params *FilterParams) splitDeletesReplace(text string) string {
	text = params.Split.split(text)

	text = deletes(params.Deletes, text)

	text = replaces(params.Replaces, text)
	return text
}
