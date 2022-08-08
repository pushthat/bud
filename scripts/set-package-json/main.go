package main

import (
	"os"
	"path/filepath"

	"github.com/pushthat/bud/internal/npm"
	"github.com/pushthat/bud/internal/versions"
	"github.com/pushthat/bud/package/gomod"
	"github.com/pushthat/bud/package/log/console"
)

func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	dir, err := gomod.Absolute(".")
	if err != nil {
		return err
	}
	// Update the dependencies in ./livebud/package.json
	if err := npm.Set(filepath.Join(dir, "livebud"), map[string]string{
		"dependencies.svelte":              versions.Svelte,
		"dependencies.react":               versions.React,
		"dependencies.react-dom":           versions.React,
		"devDependencies.@types/react":     versions.React,
		"devDependencies.@types/react-dom": versions.React,
	}); err != nil {
		return err
	}
	// Update the dependencies in .
	if err := npm.Set(dir, map[string]string{
		"devDependencies.svelte":           versions.Svelte,
		"devDependencies.react":            versions.React,
		"devDependencies.react-dom":        versions.React,
		"devDependencies.@types/react":     versions.React,
		"devDependencies.@types/react-dom": versions.React,
	}); err != nil {
		return err
	}
	return nil
}
