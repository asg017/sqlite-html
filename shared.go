package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	funcs "github.com/asg017/sqlite-dom/funcs"
	"github.com/augmentable-dev/vtab"
	"go.riyazali.net/sqlite"
)


type DomCursor struct {
	current int

	document *goquery.Document
	children *goquery.Selection
}

func (cur *DomCursor) Column(ctx *sqlite.Context, c int) error {

	col := cols[c].Name
	switch(col) {
		case "document":
		ctx.ResultText("xxx")
		case "selector":
		ctx.ResultText("sel")

		case "i":
			ctx.ResultInt(cur.current)
		case "html":
		html, _ := goquery.OuterHtml(cur.children.Eq(cur.current))
		ctx.ResultText(html)
		case "text":
		ctx.ResultText(cur.children.Eq(cur.current).Text())
		case "length": // length
		ctx.ResultInt(cur.children.Eq(cur.current).Length())
	
	}
	return nil
} 

func (cur *DomCursor) Next() (vtab.Row, error) {
	cur.current += 1
	if cur.current >= cur.children.Size() {
		return nil, io.EOF
	}
	return cur, nil
}

var cols = []vtab.Column{
	{Name: "document", Type: sqlite.SQLITE_TEXT.String(), NotNull: true, Hidden: true, Filters: []*vtab.ColumnFilter{{Op: sqlite.INDEX_CONSTRAINT_EQ, Required: true, OmitCheck: true}}},
	{Name: "selector", Type: sqlite.SQLITE_TEXT.String(), NotNull: true, Hidden: true, Filters: []*vtab.ColumnFilter{{Op: sqlite.INDEX_CONSTRAINT_EQ, Required: true, OmitCheck: true}}},
	
	{Name: "i", Type: sqlite.SQLITE_INTEGER.String()},
	{Name: "html", Type: sqlite.SQLITE_TEXT.String()},
	{Name: "text", Type: sqlite.SQLITE_TEXT.String()},
	{Name: "length", Type: sqlite.SQLITE_INTEGER.String()},
}

func NewDomModule() sqlite.Module {
	m:= vtab.NewTableFunc("dom_querySelectorAll", cols, func(constraints []*vtab.Constraint, order []*sqlite.OrderBy) (vtab.Iterator, error) {
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
			fmt.Println(err)
			return nil, sqlite.SQLITE_ABORT
		}

		children := doc.Find(selector)
		current := -1


		return &DomCursor{
			current:  current,
			document: doc,
			children: children,
		}, nil
	})
	
	return m
}

var functions = map[string]sqlite.Function{
	"dom_$": &funcs.DomQuerySelectorHtmlFunc{},
	"dom_$text": &funcs.DomQuerySelectorTextFunc{},
	"dom_$json": &funcs.DomQuerySelectorJsonFunc{},
	"dom_count": &funcs.DomCount{},
	"dom_trim": &funcs.DomTrimFunc{},
	"dom_table": &funcs.DomTableFunc{},
	"dom_attr_get": &funcs.DomAttributeGet{},
	"assert": &funcs.Assert{},
}

var modules = map[string]sqlite.Module{
	"dom_$$": NewDomModule(),
}

func init() {

	sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
		for name, function := range functions {
			if err := api.CreateFunction(name, function); err != nil {
				return sqlite.SQLITE_ERROR, err
			}
		}
		for name, module := range modules {
			if err := api.CreateModule(name, module); err != nil {
				return sqlite.SQLITE_ERROR, err
			}	
		}
		return sqlite.SQLITE_OK, nil
	})
}

func main() {}
