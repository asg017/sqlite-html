package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.riyazali.net/sqlite"
	"golang.org/x/net/html"
)

// A random "magic" number to use for sqlite subtypes.
// only the lower 8 bits are used. chose at random tbh
// https://www.sqlite.org/c3ref/result_subtype.html
var HTML_SUBTYPE = 0xdd

func subtypeIsHtml(value sqlite.Value) bool {
	return value.SubType() == HTML_SUBTYPE
}

// https://github.com/sqlite/sqlite/blob/8b554e2a1ea4de0cb30a49357684836710f44905/ext/misc/json1.c#L159
const JSON_SUBTYPE = 74

func subtypeIsJson(v sqlite.Value) bool {
	return v.SubType() == JSON_SUBTYPE
}

/**		html(document)
 *	Verifies and "cleans" (quotes attributes) the given document as HTML.
 *  Also sets the return subtype to the HTML magic number, for
 *  use in other funcs like html_element to designate something as "HTML"
 **/
 type HtmlFunc struct{}

 func (*HtmlFunc) Deterministic() bool { return true }
 func (*HtmlFunc) Args() int           { return 1 }
 func (*HtmlFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	 html := values[0].Text()
	 doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
 
	 if err != nil {
		 c.ResultError(err)
		 return
	 }
	 // goquery seems to wrap with "<html><body>" stuff, so only get original
	 outHtml, err := doc.Find("body").Html()
 
	 if err != nil {
		 c.ResultError(err)
		 return
	 }
 
	 c.ResultText(outHtml)
	 c.ResultSubType(HTML_SUBTYPE)
 }

/** 	html_element(tag, attributes, child1, ...)
 * Create an HTML element with the given tag, attributes, and children.
 * Modeled after React.createElement https://reactjs.org/docs/react-without-jsx.html
 * @param tag {text} - required top-level tag name for the returned root element.
 * @param attributes {json} - should be a json object with string keys/values, for the attributes of the element.
 * @param children {text | html} - are either strings or html elements.
 * 	If 'children" is a string, then it will be rendered as a TextNode in the top-level element
 * 	If 'children" is an html element, from "html()" or "html_element()" then it will be rendered as a RawNode in the top-level element
 */
type HtmlElementFunc struct{}

func (*HtmlElementFunc) Deterministic() bool { return true }
func (*HtmlElementFunc) Args() int           { return -1 }
func (*HtmlElementFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	var children []sqlite.Value
	if len(values) < 1 {
		c.ResultError(errors.New("html_element requires tag"))
	}
	tag := values[0].Text()
	if len(values) >= 3 {
		children = values[2:]
	}
	var attr []html.Attribute

	if len(values) > 1 && values[1].Type() != sqlite.SQLITE_NULL {
		rawAttrs := values[1].Text()

		var attrs map[string]string
		if err := json.Unmarshal([]byte(rawAttrs), &attrs); err != nil {
			c.ResultError(errors.New("attributes is not a JSON object"))
		}

		for k, v := range attrs {
			attr = append(attr, html.Attribute{
				Key: k,
				Val: v,
			})
		}

	}

	root := &html.Node{
		Type: html.ElementNode,
		Data: tag,
		Attr: attr,
	}

	for _, v := range children {
		var child *html.Node
		childData := v.Text()

		if subtypeIsHtml(v) {
			child = &html.Node{
				Type: html.RawNode,
				Data: childData,
			}
		} else {
			child = &html.Node{
				Type: html.TextNode,
				Data: childData,
			}
		}
		root.AppendChild(child)
	}

	var buf bytes.Buffer
	html.Render(&buf, root)
	c.ResultText(buf.String())
	c.ResultSubType(HTML_SUBTYPE)
}

type HtmlGroupElementFunc struct{
	parent string
}

func (h *HtmlGroupElementFunc) Args() int           { return -1 }
func (h *HtmlGroupElementFunc) Deterministic() bool { return true }

type HtmlGroupElementContext struct {
	root *html.Node
}

func (s *HtmlGroupElementFunc) Step(ctx *sqlite.AggregateContext, values ...sqlite.Value) {
	if len(values) < 1 {
		ctx.ResultError(errors.New("html_element requires tag"))
	}
	
	if ctx.Data() == nil {
		root := &html.Node{
			Type: html.ElementNode,
			Data: s.parent,
			Attr: nil,
		}
		ctx.SetData(&HtmlGroupElementContext{root:root,})
	}

	var gCtx = ctx.Data().(*HtmlGroupElementContext)
	var children []sqlite.Value

	if len(values) >= 3 {
		children = values[2:]
	}
	var attr []html.Attribute

	if len(values) > 1 && values[1].Type() != sqlite.SQLITE_NULL {
		rawAttrs := values[1].Text()

		var attrs map[string]string
		if err := json.Unmarshal([]byte(rawAttrs), &attrs); err != nil {
			ctx.ResultError(errors.New("attributes is not a JSON object"))
		}

		for k, v := range attrs {
			attr = append(attr, html.Attribute{
				Key: k,
				Val: v,
			})
		}

	}


	node := &html.Node{
		Type: html.ElementNode,
		Data: values[0].Text(),
		Attr: attr,
	}
	for _, v := range children {
		var child *html.Node
		childData := v.Text()

		if subtypeIsHtml(v) {
			child = &html.Node{
				Type: html.RawNode,
				Data: childData,
			}
		} else {
			child = &html.Node{
				Type: html.TextNode,
				Data: childData,
			}
		}
		node.AppendChild(child)
	}
	gCtx.root.AppendChild(node)
}

func (s *HtmlGroupElementFunc) Final(ctx *sqlite.AggregateContext) {
	if ctx.Data() != nil {
		var buf bytes.Buffer
		var gCtx = ctx.Data().(*HtmlGroupElementContext)
		html.Render(&buf, gCtx.root)
		//gCtx.root.FirstChild.N
		ctx.ResultText(buf.String())
		ctx.ResultSubType(HTML_SUBTYPE)
	}
}

// select html_group_elements('tr', null, value) from json_each('[1,2,3,4]')
func RegisterElements(api *sqlite.ExtensionApi) error {
	var err error
	if err = api.CreateFunction("html", &HtmlFunc{}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_element", &HtmlElementFunc{}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_group_element_div", &HtmlGroupElementFunc{parent: "div"}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_group_element_span", &HtmlGroupElementFunc{parent: "span"}); err != nil {
		return err
	}
	return nil
}
