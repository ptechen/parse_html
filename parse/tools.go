package parse

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

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

func clear(text string, params *FilterParams) string {
	if params.Split != nil {
		if params.Split.Keys != nil && len(params.Split.Keys) > 0 {
			for i := 0; i < len(params.Split.Keys); i++{
				key := params.Split.Keys[i]
				if key == "\\n" {
					text = strings.Split(text, "\n")[params.Split.Index]
				} else {
					text = strings.Split(text, params.Split.Key)[params.Split.Index]
				}
			}

		}

	}

	if params.Contains != "" {
		if !strings.Contains(text, params.Contains) {
			return ""
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