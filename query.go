package main

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/augmentable-dev/vtab"
	"go.riyazali.net/sqlite"
)

/** html_text(document [, selector])
 * Returns the combined text contents of the selected element. similar to .innerText
 * Raises an error if document is not proper HTML.
 * @param document {text | html} - HTML document to read from.
 * @param selector {text} - CSS-style selector of which element in document to read.
 */
 type HtmlTextFunc struct{
	nArgs int
 }

 func (*HtmlTextFunc) Deterministic() bool { return true }
 func (h *HtmlTextFunc) Args() int           { return h.nArgs }
 func (*HtmlTextFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	 html := values[0].Text()
	 doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
 
	 if err != nil {
		 c.ResultError(err)
		 return
	 }
	 if len(values) > 1 {
		selector := values[1].Text()
		c.ResultText(doc.FindMatcher(goquery.Single(selector)).Text())
	 }else {
		c.ResultText(doc.Text())
	 } 
 }

/** html_extract(document, selector)
 * Returns the entire HTML representation of the selected element from document, using selector.
 * Raises an error if document is not proper HTML.
 * @param document {text | html} - HTML document to read from.
 * @param selector {text} - CSS-style selector of which element in document to read.
 */
type HtmlExtractFunc struct{}

func (*HtmlExtractFunc) Deterministic() bool { return true }
func (*HtmlExtractFunc) Args() int           { return 2 }
func (*HtmlExtractFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	html := values[0].Text()
	selector := values[1].Text()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

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
	c.ResultSubType(HTML_SUBTYPE)
}

/** html_count(document, selector)
 * Count the number of matching selected elements in the given document.
 * Raises an error if document is not proper HTML.
 * @param document {text | html} - HTML document to read from.
 * @param selector {text} - CSS-style selector of which element in document to read.
 */
type HtmlCountFunc struct{}

func (*HtmlCountFunc) Deterministic() bool { return true }
func (*HtmlCountFunc) Args() int           { return 2 }
func (*HtmlCountFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	html := values[0].Text()
	selector := values[1].Text()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		c.ResultError(err)
		return
	}

	count := doc.Find(selector).Length()

	c.ResultInt(count)
}

/** html_each(document, selector)
 * A table value function returned a row for every matching element inside document using selector.
 * Raises an error if document is not proper HTML.
 * @param document {text | html} - HTML document to read from.
 * @param selector {text} - CSS-style selector of which element in document to read.
 */
 var HtmlEachColumns = []vtab.Column{
	{Name: "document", Type: sqlite.SQLITE_TEXT.String(), NotNull: true, Hidden: true, Filters: []*vtab.ColumnFilter{{Op: sqlite.INDEX_CONSTRAINT_EQ, Required: true, OmitCheck: true}}},
	{Name: "selector", Type: sqlite.SQLITE_TEXT.String(), NotNull: true, Hidden: true, Filters: []*vtab.ColumnFilter{{Op: sqlite.INDEX_CONSTRAINT_EQ, Required: true, OmitCheck: true}}},

	{Name: "html", Type: sqlite.SQLITE_TEXT.String()},
	{Name: "text", Type: sqlite.SQLITE_TEXT.String()},
}

 type HtmlEachCursor struct {
	current int

	document *goquery.Document
	children *goquery.Selection
}

func (cur *HtmlEachCursor) Column(ctx *sqlite.Context, c int) error {

	col := HtmlEachColumns[c].Name
	switch col {
	case "document":
		ctx.ResultText("")
	case "selector":
		ctx.ResultText("")

	case "html":
		html, err := goquery.OuterHtml(cur.children.Eq(cur.current))
		if err != nil {
			ctx.ResultError(err)
		} else {
			ctx.ResultText(html)
			ctx.ResultSubType(HTML_SUBTYPE)
		}
	case "text":
		ctx.ResultText(cur.children.Eq(cur.current).Text())
	}
	return nil
}

func (cur *HtmlEachCursor) Next() (vtab.Row, error) {
	cur.current += 1
	if cur.current >= cur.children.Size() {
		return nil, io.EOF
	}
	return cur, nil
}

func HtmlEachIterator(constraints []*vtab.Constraint, order []*sqlite.OrderBy) (vtab.Iterator, error) {
	document := ""
	selector := ""

	for _, constraint := range constraints {
		if constraint.Op == sqlite.INDEX_CONSTRAINT_EQ {
			switch constraint.ColIndex {
			case 0:
				document = constraint.Value.Text()
			case 1:
				selector = constraint.Value.Text()
			}
		}
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return nil, sqlite.SQLITE_ABORT
	}

	children := doc.Find(selector)
	current := -1

	return &HtmlEachCursor{
		current:  current,
		document: doc,
		children: children,
	}, nil
}

func RegisterQuery(api *sqlite.ExtensionApi) error {
	var err error
	if err = api.CreateFunction("html_extract", &HtmlExtractFunc{}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_text", &HtmlTextFunc{nArgs: 1}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_text", &HtmlTextFunc{nArgs: 2}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_count", &HtmlCountFunc{}); err != nil {
		return err
	}
	if err = api.CreateModule("html_each", vtab.NewTableFunc("html_each", HtmlEachColumns, HtmlEachIterator)); err != nil {
		return err
	}
	return nil
}
