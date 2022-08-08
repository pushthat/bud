package view

import (
	"github.com/pushthat/bud/framework"
	"github.com/pushthat/bud/internal/embed"
	"github.com/pushthat/bud/internal/imports"
)

type State struct {
	Imports []*imports.Import
	Flag    *framework.Flag
	Embeds  []*embed.File
}
