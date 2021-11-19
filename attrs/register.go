package attrs

import (
	"github.com/asg017/sqlite-html/internal"
	"go.riyazali.net/sqlite"
)

var Funcs = []internal.Function{
	{
		Name:       "html_attribute_get",
		Aliases:    []string{"html_attr_get"},
		Definition: &HtmlAttributeGetFunc{},
		Documentation: map[string]string{
			"description": "",
		},
	},
	{
		Name:       "html_attribute_has",
		Aliases:    []string{"html_attr_has"},
		Definition: &HtmlAttributeHasFunc{},
		Documentation: map[string]string{
			"description": "",
		},
	},
}

// TODO html_attrs_each
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
