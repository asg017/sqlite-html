package main

import (
	"go.riyazali.net/sqlite"
)

// following linker flags are needed to suppress missing symbol warning in intermediate stages

// #cgo linux LDFLAGS: -Wl,--unresolved-symbols=ignore-in-object-files
// #cgo darwin LDFLAGS: -Wl,-undefined,dynamic_lookup
// #cgo windows LDFLAGS: -Wl,--allow-multiple-definition
import "C"

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
	sqlite.RegisterNamed("html", Register)
	sqlite.Register(Register)
}

func main() {}
