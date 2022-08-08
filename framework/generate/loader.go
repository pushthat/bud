package generate

import (
	"fmt"
	"io/fs"

	"github.com/pushthat/bud/internal/bail"
	"github.com/pushthat/bud/internal/imports"
	"github.com/pushthat/bud/package/di"
	"github.com/pushthat/bud/package/gomod"
	"github.com/pushthat/bud/package/vfs"
)

func Load(fsys fs.FS, injector *di.Injector, module *gomod.Module) (*State, error) {
	loader := &loader{
		imports:  imports.New(),
		injector: injector,
		module:   module,
	}
	return loader.Load(fsys)
}

type loader struct {
	bail.Struct
	imports  *imports.Set
	injector *di.Injector
	module   *gomod.Module
}

// Load the command state
func (l *loader) Load(fsys fs.FS) (state *State, err error) {
	defer l.Recover2(&err, "generate")
	if err := vfs.Exist(fsys, "bud/internal/generate/generator/generator.go"); err != nil {
		return nil, err
	}
	state = new(State)
	state.Provider = l.loadProvider()
	state.Imports = l.loadImports()
	return state, nil
}

func (l *loader) loadProvider() *di.Provider {
	provider, err := l.injector.Wire(&di.Function{
		Name:    "loadGenerator",
		Imports: l.imports,
		Params: []*di.Param{
			{Import: "github.com/pushthat/bud/package/log", Type: "Interface"},
			{Import: "github.com/pushthat/bud/package/gomod", Type: "*Module"},
			{Import: "context", Type: "Context"},
		},
		Results: []di.Dependency{
			di.ToType(l.module.Import("bud/internal/generate/generator"), "*FileSystem"),
			&di.Error{},
		},
	})
	if err != nil {
		l.Bail(fmt.Errorf("unable to load provider: %s", err))
	}
	return provider
}

func (l *loader) loadImports() []*imports.Import {
	l.imports.AddStd("os", "context", "errors")
	l.imports.AddNamed("commander", "github.com/pushthat/bud/package/commander")
	l.imports.AddNamed("console", "github.com/pushthat/bud/package/log/console")
	l.imports.AddNamed("log", "github.com/pushthat/bud/package/log")
	l.imports.AddNamed("filter", "github.com/pushthat/bud/package/log/filter")
	l.imports.AddNamed("remotefs", "github.com/pushthat/bud/package/remotefs")
	return l.imports.List()
}
