package generator

import "github.com/pushthat/bud/internal/imports"

type State struct {
	Imports    []*imports.Import
	Generators []*Gen
}

type Gen struct {
	Import *imports.Import
	Path   string
	Pascal string
}
