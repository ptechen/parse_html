package parse_html

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseHtml(t *testing.T) {
	params := make(map[string]*FilterParams)
	//params["val"] = &FilterParams{
	//	Selector: "#job-view-enterprise > div.wrap.clearfix > div.clearfix > div.main > div.about-position > div:nth-child(2) > div.clearfix > div.job-title-left > p.job-item-title",
	//	Split: &Split{
	//		Key:   "\n",
	//		Index: 0,
	//	},
	//	Finds: []string{".job-title-left", ".job-item-title"},
	//}

	//mapKey := make(map[string]*FilterParams)
	params["job_name"] = &FilterParams{
		Finds:    []string{".comp-summary-tag", "a"},
		Type:     "contains_list",
		Keys:     nil,
		Last:     false,
		First:    false,
		Eq:       0,
		Attr:     "title",
		Splits:    nil,
		Contains: &Contain{
			HasAttr: &HasAttr{
				Key: "data-selector",
				Val: "comp-industry",
			},
		},
		Deletes:  nil,
		Replaces: nil,
	}

	//params["position_list"] = &FilterParams{
	//	Finds:    []string{".sojob-list", "li"},
	//	Type:     "list",
	//	Keys:     mapKey,
	//}
	dataBytes, err := ioutil.ReadFile("index.html")
	if err != nil {
		t.Error(err)
	}
	dataStr := string(dataBytes)
	res, err := ParseHtml(dataStr, params)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("%#v", res))
}