package main

import (
	attrs "github.com/asg017/sqlite-html/attrs"
	elements "github.com/asg017/sqlite-html/elements"
	meta "github.com/asg017/sqlite-html/meta"
	query "github.com/asg017/sqlite-html/query"
	utils "github.com/asg017/sqlite-html/utils"
	"go.riyazali.net/sqlite"
)

// Set in Makefile
var (
	Commit  string
	Date    string
	Version string
)

func Register(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {

	if err := meta.Register(api, meta.RegisterParams{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}); err != nil {
		return sqlite.SQLITE_ERROR, err
	}

	if err := attrs.Register(api); err != nil {
		return sqlite.SQLITE_ERROR, err
	}
	if err := query.Register(api); err != nil {
		return sqlite.SQLITE_ERROR, err
	}
	if err := utils.Register(api); err != nil {
		return sqlite.SQLITE_ERROR, err
	}
	if err := elements.Register(api); err != nil {
		return sqlite.SQLITE_ERROR, err
	}
	return sqlite.SQLITE_OK, nil
}

func init() {
	sqlite.Register(Register)
}

func main() {}
