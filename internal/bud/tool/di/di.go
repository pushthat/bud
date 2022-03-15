package di

import (
	"context"
	"fmt"
	"os"
	"strings"

	"gitlab.com/mnm/bud/internal/bud"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

type Command struct {
	Bud          *bud.Command
	Target       string
	Map          map[string]string
	Dependencies []string
	Externals    []string
	Hoist        bool
	Verbose      bool
}

func (c *Command) Run(ctx context.Context) error {
	module, err := gomod.Find(c.Bud.Dir)
	if err != nil {
		return err
	}
	overlay, err := overlay.Load(module)
	if err != nil {
		return err
	}
	parser := parser.New(overlay, module)
	fn := &di.Function{
		Hoist: c.Hoist,
	}
	fn.Target, err = c.toImportPath(module, c.Target)
	if err != nil {
		return err
	}
	// Add the type mapping
	for from, to := range c.Map {
		fromDep, err := c.toDependency(module, from)
		if err != nil {
			return err
		}
		toDep, err := c.toDependency(module, to)
		if err != nil {
			return err
		}
		fn.Aliases[fromDep] = toDep
	}
	// Add the dependencies
	for _, dependency := range c.Dependencies {
		dep, err := c.toDependency(module, dependency)
		if err != nil {
			return err
		}
		fn.Results = append(fn.Results, dep)
	}
	// Add the externals
	for _, external := range c.Externals {
		ext, err := c.toDependency(module, external)
		if err != nil {
			return err
		}
		fn.Params = append(fn.Params, ext)
	}
	injector := di.New(overlay, module, parser)
	node, err := injector.Load(fn)
	if err != nil {
		return err
	}
	if c.Verbose {
		fmt.Println(node.Print())
	}
	provider := node.Generate(imports.New(), "Load", fn.Target)
	fmt.Fprintln(os.Stdout, provider.File())
	return nil
}

// This should handle both stdlib (e.g. "net/http"), directories (e.g. "web"),
// and dependencies
func (c *Command) toImportPath(module *gomod.Module, importPath string) (string, error) {
	importPath = strings.Trim(importPath, "\"")
	maybeDir := module.Directory(importPath)
	if _, err := os.Stat(maybeDir); err == nil {
		importPath, err = module.ResolveImport(maybeDir)
		if err != nil {
			return "", fmt.Errorf("di: unable to resolve import %s because %+s", importPath, err)
		}
	}
	return importPath, nil
}

func (c *Command) toDependency(module *gomod.Module, dependency string) (di.Dependency, error) {
	i := strings.LastIndex(dependency, ".")
	if i < 0 {
		return nil, fmt.Errorf("di: external must have form '<import>.<type>'. got %q ", dependency)
	}
	importPath, err := c.toImportPath(module, dependency[0:i])
	if err != nil {
		return nil, err
	}
	dataType := dependency[i+1:]
	// Create the dependency
	return &di.Type{
		Import: importPath,
		Type:   dataType,
	}, nil
}