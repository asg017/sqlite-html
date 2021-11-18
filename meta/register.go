package meta

import (
	"github.com/asg017/sqlite-dom/internal"
	"github.com/augmentable-dev/vtab"
	"go.riyazali.net/sqlite"
)

type RegisterParams struct {
	Version string
	Commit  string
	Date    string
}

func buildFunctions(params RegisterParams) []internal.Function {
	return []internal.Function{
		{
			Name:       "http_version",
			Definition: &HtmlVersionFunc{version: params.Version},
			Documentation: map[string]string{
				"description": "",
			},
		},
		{
			Name:       "http_debug",
			Definition: &HtmlDebugFunc{version: params.Version, date: params.Date, commit: params.Commit},
			Documentation: map[string]string{
				"description": "",
			},
		},
	}
}

func buildModules(params RegisterParams) []internal.Module {
	return []internal.Module{
		{
			Name:       "html_docs",
			Definition: vtab.NewTableFunc("html_docs", HtmlDocsColumns, HtmlDocsBuildIterator(params)),
			Documentation: map[string]string{
				"description": "Table function that returns one row for every matching documentation item",
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

}

func Register(api *sqlite.ExtensionApi, params RegisterParams) error {
	functions := buildFunctions(params)
	modules := buildModules(params)

	for _, f := range functions {
		if err := f.Register(api); err != nil {
			return err
		}
	}

	for _, module := range modules {
		// skip html_docs for now, not needed
		if true {
			continue
		}
		if err := module.Register(api); err != nil {
			return err
		}
	}
	return nil
}
