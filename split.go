package parse_html

import (
	"encoding/json"
	"strings"
)

type Split struct {
	Key    string `json:"key" yaml:"key"`
	Index  int    `json:"index" yaml:"index"`
	Enable bool   `json:"enable" yaml:"enable"`
}

func (params *Split) split(text string) string {
	if params != nil {
		if params.Enable {
			if params.Key == "\\n" {
				texts := strings.Split(text, "\n")
				if len(texts) > params.Index {
					text = texts[params.Index]
				} else {
					text = strings.Join(texts, ",")
				}
			} else if params.Key == "\\t" {
				texts := strings.Split(text, "\t")
				if len(texts) > params.Index {
					text = texts[params.Index]
				} else {
					text = strings.Join(texts, ",")
				}
			} else {
				texts := strings.Split(text, params.Key)
				if len(texts) > params.Index {
					text = texts[params.Index]
				} else {
					text = strings.Join(texts, ",")
				}
			}
		} else {
			dataBytes := make([]byte, 0)
			if params.Key == "\\n" {
				dataBytes, _ = json.Marshal(strings.Split(text, "\n"))
				text = string(dataBytes)
			} else if params.Key == "\\t" {
				dataBytes, _ = json.Marshal(strings.Split(text, "\t"))

			} else {
				dataBytes, _ = json.Marshal(strings.Split(text, params.Key))
			}
			text = string(dataBytes)
		}
	}
	return text
}