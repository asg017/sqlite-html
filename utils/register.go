package utils

import (
	"github.com/asg017/sqlite-html/internal"
	"go.riyazali.net/sqlite"
)

var Funcs = []internal.Function{
	{
		Name:       "html_table",
		Definition: &HtmlTableFunc{},
		Documentation: map[string]string{
			"description": "",
		},
	},
	{
		Name:       "html_escape",
		Definition: &HtmlEscapeFunc{},
		Documentation: map[string]string{
			"description": "",
		},
	},
	{
		Name:       "html_unescape",
		Definition: &HtmlUnescapeFunc{},
		Documentation: map[string]string{
			"description": "",
		},
	},
	{
		Name:       "html_trim",
		Definition: &HtmlTrimFunc{},
		Documentation: map[string]string{
			"description": "HtmlTrimFunc",
		},
	},
}

var Modules = []internal.Module{}

func Register(api *sqlite.ExtensionApi) error {
	for _, f := range Funcs {
		if err := f.Register(api); err != nil {
			return err
		}
	}
	for _, module := range Modules {
		if err := module.Register(api); err != nil {
			return err
		}
	}

	return nil
}
