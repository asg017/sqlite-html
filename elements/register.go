package elements

import (
	"github.com/asg017/sqlite-dom/internal"
	"go.riyazali.net/sqlite"
)

var Funcs = []internal.Function{
	{
		Name:       "html",
		Definition: &HtmlFunc{},
		Documentation: map[string]string{
			"description": "",
		},
	},
	{
		Name:       "html_element",
		Definition: &HtmlElementFunc{},
		Documentation: map[string]string{
			"description": "",
		},
	},
}

var Modules = []internal.Module{}

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
