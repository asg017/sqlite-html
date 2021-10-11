package funcs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.riyazali.net/sqlite"
)



type DomQuerySelectorTextFunc struct{}

func (*DomQuerySelectorTextFunc) Deterministic() bool { return true }
func (*DomQuerySelectorTextFunc) Args() int           { return 2 }
func (*DomQuerySelectorTextFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	html := values[0].Text()
	selector := values[1].Text()
	
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html));

	if err != nil {
		c.ResultError(err)
		return
	}
	c.ResultText(doc.FindMatcher(goquery.Single(selector)).Text())
}

type DomQuerySelectorHtmlFunc struct{}

func (*DomQuerySelectorHtmlFunc) Deterministic() bool { return true }
func (*DomQuerySelectorHtmlFunc) Args() int           { return 2 }
func (*DomQuerySelectorHtmlFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	html := values[0].Text()
	selector := values[1].Text()
	
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html));

	if err != nil {
		c.ResultError(err)
		return
	}
	
	sub, err := goquery.OuterHtml(doc.FindMatcher(goquery.Single(selector)))
	if err != nil {
		c.ResultError(err)
		return
	}

	c.ResultText(sub)
}

type DomElementJson struct {
	TextContent string `json: "textContent"`
	Attributes []string `json "attributes"`
	TagName string `json: "tagName"`
}
type DomQuerySelectorJsonFunc struct{}
func (*DomQuerySelectorJsonFunc) Deterministic() bool { return true }
func (*DomQuerySelectorJsonFunc) Args() int           { return 2 }
func (*DomQuerySelectorJsonFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	html := values[0].Text()
	selector := values[1].Text()
	
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html));

	if err != nil {
		c.ResultError(err)
		return
	}
	el := doc.FindMatcher(goquery.Single(selector))

		txt := el.Text()

		attrs := make([]string, len(el.Get(0).Attr))

		for _, attr := range el.Get(0).Attr {
			attrs = append(attrs, attr.Key, attr.Val)
		}

		bytes, err := json.Marshal(DomElementJson{
			TextContent: txt,
			Attributes: attrs,
			TagName:  goquery.NodeName(el),
	})
		if err != nil {
			panic(err)
		}

	c.ResultText(string(bytes))
}

type DomTrimFunc struct{}
func (*DomTrimFunc) Deterministic() bool { return true }
func (*DomTrimFunc) Args() int           { return 1 }
func (*DomTrimFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	s := values[0].Text()
	c.ResultText(strings.TrimSpace(s))
}

type DomTableFunc struct{}
func (*DomTableFunc) Deterministic() bool { return true }
func (*DomTableFunc) Args() int           { return 1 }
func (*DomTableFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	s := values[0].Text()
	c.ResultText(fmt.Sprintf("<table>%s", s))
}



type DomAttributeGet struct{}
func (*DomAttributeGet) Deterministic() bool {return true}
func (*DomAttributeGet) Args() int           { return 3 }
func (*DomAttributeGet) Apply(c *sqlite.Context, values ...sqlite.Value) {
	html := values[0].Text()
	selector := values[1].Text()
	attribute := values[2].Text()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html));

	if err != nil {
		c.ResultError(err)
		return
	}

	attr, exists := doc.FindMatcher(goquery.Single(selector)).Attr(attribute)

	if !exists {
		c.ResultNull()
	}else {
		c.ResultText(attr)
	}
	
}

type DomCount struct{}
func (*DomCount) Deterministic() bool {return true}
func (*DomCount) Args() int           { return 2 }
func (*DomCount) Apply(c *sqlite.Context, values ...sqlite.Value) {
	html := values[0].Text()
	selector := values[1].Text()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html));

	if err != nil {
		c.ResultError(err)
		return
	}

	count := doc.Find(selector).Length()

	c.ResultInt(count)
}

type Assert struct{}
func (*Assert) Deterministic() bool {return true}
func (*Assert) Args() int           { return -1 }
func (*Assert) Apply(c *sqlite.Context, values ...sqlite.Value) {
	result := values[0].Int()

	if result == 0 {
		c.ResultError(errors.New("ASSERTION FAILED"))
	}else {
		c.ResultText("✅")
	}
}