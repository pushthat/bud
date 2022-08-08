package app

import (
	"github.com/pushthat/bud/framework"
	"github.com/pushthat/bud/internal/imports"
	"github.com/pushthat/bud/package/di"
)

type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
	Flag     *framework.Flag
}
