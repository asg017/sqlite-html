package meta

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/asg017/sqlite-dom/attrs"
	"github.com/asg017/sqlite-dom/elements"
	"github.com/asg017/sqlite-dom/internal"
	"github.com/asg017/sqlite-dom/query"
	"github.com/asg017/sqlite-dom/utils"
	"github.com/augmentable-dev/vtab"
	"go.riyazali.net/sqlite"
)

var allFuncs = [][]internal.Function{
	attrs.Funcs,
	elements.Funcs,
	query.Funcs,
	utils.Funcs,
}

var allModules = [][]internal.Module{
	attrs.Modules,
	elements.Modules,
	query.Modules,
	// Can't document itself yet :(
	//Modules,
}

// html_version()
// return the version string of the current sqlite-html module.
// passed in from top-level Makefile and build args
type HtmlVersionFunc struct {
	version string
}

func (*HtmlVersionFunc) Deterministic() bool { return true }
func (*HtmlVersionFunc) Args() int           { return 0 }
func (f *HtmlVersionFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	c.ResultText(f.version)
}

// html_debug()
// Returns more information for the current html module, including build date + comment hash.
type HtmlDebugFunc struct {
	version string
	commit  string
	date    string
}

func (*HtmlDebugFunc) Deterministic() bool { return true }
func (*HtmlDebugFunc) Args() int           { return 0 }
func (f *HtmlDebugFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	c.ResultText(fmt.Sprintf("Version: %s\nCommit: %s\nRuntime: %s %s/%s\nDate: %s\n",
		f.version,
		f.commit,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		f.date,
	))
}

// bro idk if this works
var HtmlDocsColumns = []vtab.Column{
	{Name: "_name", Type: sqlite.SQLITE_TEXT.String(), NotNull: true, Hidden: true, Filters: []*vtab.ColumnFilter{{Op: sqlite.INDEX_CONSTRAINT_EQ, Required: false, OmitCheck: true}}},
	{Name: "_type", Type: sqlite.SQLITE_TEXT.String(), NotNull: true, Hidden: true, Filters: []*vtab.ColumnFilter{{Op: sqlite.INDEX_CONSTRAINT_EQ, Required: false, OmitCheck: true}}},
	{Name: "source", Type: sqlite.SQLITE_TEXT.String()},
	{Name: "type", Type: sqlite.SQLITE_TEXT.String()},
	{Name: "name", Type: sqlite.SQLITE_TEXT.String()},
	{Name: "column", Type: sqlite.SQLITE_TEXT.String()},
	{Name: "tag", Type: sqlite.SQLITE_TEXT.String()},
	{Name: "value", Type: sqlite.SQLITE_TEXT.String()},
}

type HtmlDocsRow struct {
	rType  string
	name   string
	column string
	tag    string
	value  string
}

// html_docs(name, [type])
type HtmlDocsCursor struct {
	current int
	rows    []*HtmlDocsRow
	version string
}

func (cur *HtmlDocsCursor) Column(ctx *sqlite.Context, c int) error {

	col := HtmlDocsColumns[c].Name
	switch col {
	case "source":
		ctx.ResultText(fmt.Sprintf("sqlite-html http_docs() %s", cur.version))
	case "type":
		ctx.ResultText(cur.rows[cur.current].rType)
	case "name":
		ctx.ResultText(cur.rows[cur.current].name)
	case "column":
		ctx.ResultText(cur.rows[cur.current].column)
	case "tag":
		ctx.ResultText(cur.rows[cur.current].tag)
	case "value":
		ctx.ResultText(cur.rows[cur.current].value)

	}
	return nil
}

func (cur *HtmlDocsCursor) Next() (vtab.Row, error) {
	cur.current += 1
	if cur.current >= len(cur.rows) {
		return nil, io.EOF
	}
	return cur, nil
}

func HtmlDocsBuildIterator(params RegisterParams) func(constraints []*vtab.Constraint, order []*sqlite.OrderBy) (vtab.Iterator, error) {
	return func(constraints []*vtab.Constraint, order []*sqlite.OrderBy) (vtab.Iterator, error) {
		argName := ""
		argType := "table"

		for _, constraint := range constraints {
			if constraint.Op == sqlite.INDEX_CONSTRAINT_EQ {
				colName := HtmlDocsColumns[constraint.ColIndex].Name
				switch colName {
				case "_name":
					argName = constraint.Value.Text()
				case "_type":
					argType = constraint.Value.Text()
				}
			}
		}

		rows := []*HtmlDocsRow{}

		if argType == "" || strings.ToLower(argType) == "table" {
			for _, modules := range allModules {
				for _, m := range modules {
					if argName == "" || m.Match(argName) {
						for tag, value := range m.Documentation {
							rows = append(rows, &HtmlDocsRow{
								rType:  "function",
								name:   m.Name,
								column: "",
								tag:    tag,
								value:  value,
							})
						}
						for column, docs := range m.ColumnDocumentation {
							for tag, value := range docs {
								rows = append(rows, &HtmlDocsRow{
									rType:  "table",
									name:   m.Name,
									column: column,
									tag:    tag,
									value:  value,
								})
							}
						}
					}

				}
			}
		} else if argType == "" || strings.ToLower(argType) == "function" {
			for _, funcs := range allFuncs {
				for _, f := range funcs {
					if argName == "" || f.Match(argName) {
						for tag, value := range f.Documentation {
							rows = append(rows, &HtmlDocsRow{
								rType:  "function",
								name:   f.Name,
								column: "",
								tag:    tag,
								value:  value,
							})
						}
					}
				}
			}
		} else {
			return nil, (fmt.Errorf("not a valid type parameter: '%s'", argType))
		}

		return &HtmlDocsCursor{
			current: -1,
			rows:    rows,
			version: params.Version,
		}, nil
	}
}
