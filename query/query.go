package query

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/augmentable-dev/vtab"
	"go.riyazali.net/sqlite"
)

// html_text(document, selector)
// Returns the combined text contents of the selected element from document, using selector.
// Raises an error if document is not proper HTML.
type HtmlTextFunc struct{}

func (*HtmlTextFunc) Deterministic() bool { return true }
func (*HtmlTextFunc) Args() int           { return 2 }
func (*HtmlTextFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	html := values[0].Text()
	selector := values[1].Text()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		c.ResultError(err)
		return
	}
	c.ResultText(doc.FindMatcher(goquery.Single(selector)).Text())
}

// html_extract(document, selector)
// Returns the entire HTML representation of the selected element from document, using selector.
// Raises an error if document is not proper HTML.
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
}

// html_count(document, selector)
// Count the number of matching selected elements in the given document.
// Raises an error if document is not proper HTML.
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

// html_each(document, selector)
// A table value function returned a row for every matching element inside document using selector.
// Raises an error if document is not proper HTML.
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

	case "i":
		ctx.ResultInt(cur.current)
	case "html":
		html, err := goquery.OuterHtml(cur.children.Eq(cur.current))
		if err != nil {
			ctx.ResultError(err)
		} else {
			ctx.ResultText(html)
		}
	case "text":
		ctx.ResultText(cur.children.Eq(cur.current).Text())
	case "length":
		ctx.ResultInt(cur.children.Eq(cur.current).Length())

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

var HtmlEachColumns = []vtab.Column{
	{Name: "document", Type: sqlite.SQLITE_TEXT.String(), NotNull: true, Hidden: true, Filters: []*vtab.ColumnFilter{{Op: sqlite.INDEX_CONSTRAINT_EQ, Required: true, OmitCheck: true}}},
	{Name: "selector", Type: sqlite.SQLITE_TEXT.String(), NotNull: true, Hidden: true, Filters: []*vtab.ColumnFilter{{Op: sqlite.INDEX_CONSTRAINT_EQ, Required: true, OmitCheck: true}}},

	{Name: "i", Type: sqlite.SQLITE_INTEGER.String()},
	{Name: "html", Type: sqlite.SQLITE_TEXT.String()},
	{Name: "text", Type: sqlite.SQLITE_TEXT.String()},
	{Name: "length", Type: sqlite.SQLITE_INTEGER.String()},
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
