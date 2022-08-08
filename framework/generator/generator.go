package generator

import (
	"context"
	_ "embed"

	"github.com/pushthat/bud/internal/gotemplate"
	"github.com/pushthat/bud/package/gomod"
	"github.com/pushthat/bud/package/overlay"
	"github.com/pushthat/bud/package/parser"
)

//go:embed generator.gotext
var template string

var generator = gotemplate.MustParse("framework/generator/generator.gotext", template)

func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func New(module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{module, parser}
}

type Generator struct {
	module *gomod.Module
	parser *parser.Parser
}

func (g *Generator) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	state, err := Load(fsys, g.module, g.parser)
	if err != nil {
		return err
	}
	code, err := Generate(state)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
