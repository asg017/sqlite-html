package main

import (
	"fmt"
	"runtime"

	"go.riyazali.net/sqlite"
)

/**		html_version()
 *	Returns the semver version of the current sqlite-html module.
 **/
type HtmlVersionFunc struct{}

func (*HtmlVersionFunc) Deterministic() bool { return true }
func (*HtmlVersionFunc) Args() int           { return 0 }
func (f *HtmlVersionFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	c.ResultText(Version)
}

/**		html_debug()
 *	Returns more information for the current html module, 
 * 	including build date + commit hash.
 **/
type HtmlDebugFunc struct{}

func (*HtmlDebugFunc) Deterministic() bool { return true }
func (*HtmlDebugFunc) Args() int           { return 0 }
func (f *HtmlDebugFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	c.ResultText(fmt.Sprintf("Version: %s\nCommit: %s\nRuntime: %s %s/%s\nDate: %s\n",
		Version,
		Commit,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		Date,
	))
}

func RegisterMeta(api *sqlite.ExtensionApi) error {
	var err error
	if err = api.CreateFunction("html_version", &HtmlVersionFunc{}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_debug", &HtmlDebugFunc{}); err != nil {
		return err
	}
	return nil
}
