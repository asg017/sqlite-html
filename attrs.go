package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.riyazali.net/sqlite"
)

// html_attribute_get(document, selector, name)
// Get the value of the "name" attribute from the element found in document, using selector
type HtmlAttributeGetFunc struct{}

func (*HtmlAttributeGetFunc) Deterministic() bool { return true }
func (*HtmlAttributeGetFunc) Args() int           { return 3 }
func (*HtmlAttributeGetFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	html := values[0].Text()
	selector := values[1].Text()
	attribute := values[2].Text()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		c.ResultError(err)
		return
	}

	attr, exists := doc.FindMatcher(goquery.Single(selector)).Attr(attribute)

	if !exists {
		c.ResultNull()
	} else {
		c.ResultText(attr)
	}

}

// html_attribute_has(document, selector, name)
// 1 or 0, if the "name" attribute from the element found in document, using selector, exists
type HtmlAttributeHasFunc struct{}

func (*HtmlAttributeHasFunc) Deterministic() bool { return true }
func (*HtmlAttributeHasFunc) Args() int           { return 3 }
func (*HtmlAttributeHasFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	html := values[0].Text()
	selector := values[1].Text()
	attribute := values[2].Text()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		c.ResultError(err)
		return
	}

	_, exists := doc.FindMatcher(goquery.Single(selector)).Attr(attribute)

	if !exists {
		c.ResultInt(0)
	} else {
		c.ResultInt(1)
	}

}

func RegisterAttrs(api *sqlite.ExtensionApi) error {
	var err error
	if err = api.CreateFunction("html_attribute_get", &HtmlAttributeGetFunc{}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_attr_get", &HtmlAttributeGetFunc{}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_attribute_has", &HtmlAttributeHasFunc{}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_attr_has", &HtmlAttributeHasFunc{}); err != nil {
		return err
	}
	return nil
}
