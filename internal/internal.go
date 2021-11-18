package internal

import (
	"go.riyazali.net/sqlite"
)

// a small abstraction around `sqlite.Function` to include documentation and aliases (if any)
type Function struct {
	Name          string
	Aliases       []string
	Definition    sqlite.Function
	Documentation map[string]string
}

// Register the funcion with the given sqlite.ExtensionApi, including any aliases.
func (f *Function) Register(api *sqlite.ExtensionApi) error {
	err := api.CreateFunction(f.Name, f.Definition)
	if err != nil {
		return err
	}
	for _, alias := range f.Aliases {
		err := api.CreateFunction(alias, f.Definition)
		if err != nil {
			return err
		}
	}
	return nil
}

// Check if a given name matches the function's name, or any of its aliases
func (f *Function) Match(s string) bool {
	if s == f.Name {
		return true
	}
	for _, a := range f.Aliases {
		if a == s {
			return true
		}
	}
	return false
}

// a small abstraction around sqlite.Module to include documentation
type Module struct {
	Name                string
	Aliases             []string
	Definition          sqlite.Module
	opts                []func(*sqlite.ModuleOptions)
	Documentation       map[string]string
	ColumnDocumentation map[string]map[string]string
}

// Register the given module, and any aliases, to the given sqlite.ExtensionApi
func (f *Module) Register(api *sqlite.ExtensionApi) error {
	err := api.CreateModule(f.Name, f.Definition, f.opts...)
	if err != nil {
		return err
	}
	for _, alias := range f.Aliases {
		err := api.CreateModule(alias, f.Definition, f.opts...)
		if err != nil {
			return err
		}
	}
	return nil
}

// return true if the given name matches the modules name, or any of its aliases
func (f *Module) Match(s string) bool {
	if s == f.Name {
		return true
	}
	for _, a := range f.Aliases {
		if a == s {
			return true
		}
	}
	return false
}
