package main

import (
	"fmt"
	"runtime"

	"go.riyazali.net/sqlite"
)

// html_version()
// return the version string of the current sqlite-html module.
// passed in from top-level Makefile and build args
type HtmlVersionFunc struct{}

func (*HtmlVersionFunc) Deterministic() bool { return true }
func (*HtmlVersionFunc) Args() int           { return 0 }
func (f *HtmlVersionFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	c.ResultText(Version)
}

// html_debug()
// Returns more information for the current html module, including build date + comment hash.
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
