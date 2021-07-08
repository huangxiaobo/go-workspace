package task

import "github.com/PuerkitoBio/goquery"

type CallbackFunc func(selection *goquery.Selection)

type HtmlCallback struct {
	Selector string
	Function CallbackFunc
}

type FetchTask struct {
	Url           string
	HtmlCallbacks []*HtmlCallback
}

func (t *FetchTask) OnSuccess(selector string, f CallbackFunc) {
	t.HtmlCallbacks = append(
		t.HtmlCallbacks,
		&HtmlCallback{
			Selector: selector,
			Function: f,
		},
	)
}
