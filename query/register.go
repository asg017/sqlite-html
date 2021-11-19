package query

import (
	"github.com/asg017/sqlite-html/internal"
	"github.com/augmentable-dev/vtab"
	"go.riyazali.net/sqlite"
)

var Funcs = []internal.Function{
	{
		Name:       "html_extract",
		Definition: &HtmlExtractFunc{},
		Documentation: map[string]string{
			"description": "Extract the first element from the given document at the given selector. Returns entire HTML representation.",
		},
	},
	{
		Name:       "html_text",
		Definition: &HtmlTextFunc{},
		Documentation: map[string]string{
			"description": "Extract the first element from the given document at the given selector. Returns the inner text representation of the element.",
		},
	},
	{
		Name:       "html_count",
		Definition: &HtmlCountFunc{},
		Documentation: map[string]string{
			"description": "Return the number of matching elements in the given document at the given selector.",
		},
	},
}
var Modules = []internal.Module{
	{
		Name:       "html_each",
		Definition: vtab.NewTableFunc("html_each", HtmlEachColumns, HtmlEachIterator),
		Documentation: map[string]string{
			"description": "Table function that returns one row for every matching element from the given document, using the given selector.",
		},
		ColumnDocumentation: map[string]map[string]string{
			"i": {
				"description": "Index number of the matched element, starting at 0.",
			},
			"html": {
				"description": "HTML representation of the matched element.",
			},
			"text": {
				"description": "Inner text of matched element.",
			},
			"length": {
				"description": "Length of the matched element, in number of characters.",
			},
		},
	},
}

func Register(api *sqlite.ExtensionApi) error {
	for _, f := range Funcs {
		f.Register(api)
	}
	for _, module := range Modules {
		if err := module.Register(api); err != nil {
			return err
		}
	}
	return nil
}
