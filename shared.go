package main

import (
	"go.riyazali.net/sqlite"
)

// Set in Makefile
var (
	Commit  string
	Date    string
	Version string
)

func Register(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {

	if err := RegisterMeta(api); err != nil {
		return sqlite.SQLITE_ERROR, err
	}

	if err := RegisterAttrs(api); err != nil {
		return sqlite.SQLITE_ERROR, err
	}
	if err := RegisterElements(api); err != nil {
		return sqlite.SQLITE_ERROR, err
	}
	if err := RegisterQuery(api); err != nil {
		return sqlite.SQLITE_ERROR, err
	}
	if err := RegisterUtils(api); err != nil {
		return sqlite.SQLITE_ERROR, err
	}
	return sqlite.SQLITE_OK, nil
}

func init() {
	sqlite.Register(Register)
}

func main() {}
