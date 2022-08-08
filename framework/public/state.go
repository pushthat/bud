package public

import (
	"github.com/pushthat/bud/framework"
	"github.com/pushthat/bud/internal/embed"
	"github.com/pushthat/bud/internal/imports"
)

type State struct {
	Imports []*imports.Import
	Embeds  []*embed.File
	Flag    *framework.Flag
}
