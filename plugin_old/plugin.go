package plugin

import (
	"os"
	"path/filepath"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
)

// Generator loads all the bud plugins.
//
// If the import path is "gitlab.com/mnm/bud-tailwind/transform", then you'd
// load this plugin with "bud/plugin/tailwind/transform".
type Generator struct {
	Modfile mod.File
}

func (g *Generator) GenerateDir(f gen.F, dir *gen.Dir) error {
	plugins, err := g.Modfile.Plugins()
	if err != nil {
		return err
	}
	// Generate a directory of plugin names.
	// (e.g. bud/plugin/{tailwind,markdown})
	for _, plugin := range plugins {
		plugin := plugin
		dir.Entry(plugin.Name, gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
			fis, err := os.ReadDir(plugin.Dir)
			if err != nil {
				return err
			}
			for _, fi := range fis {
				name := fi.Name()
				if !fi.IsDir() || name[0] == '_' || name[0] == '.' {
					continue
				}
				// baseDir := filepath.Join("bud", "plugin", plugin.Name, name)
				pluginDir := filepath.Join(plugin.Dir, name)
				// Serve all inner files from ${plugin.Dir}/${name}/...
				// fmt.Println(filepath.Join("bud", "plugin", plugin.Name, name), "=>", filepath.Join(pluginDir))
				dir.Entry(name, gen.ServeFS(os.DirFS(pluginDir)))
				// dir.Entry(name, bfs.ServeDir(func(f bfs.FS, entry *bfs.Entry) error {
				// 	// Switch the base from the requested to the actual.
				// 	relPath, err := filepath.Rel(baseDir, entry.Path())
				// 	if err != nil {
				// 		return err
				// 	}
				// 	absPath := filepath.Join(pluginDir, relPath)
				// 	stat, err := os.Stat(absPath)
				// 	if err != nil {
				// 		return err
				// 	}
				// 	entry.Mode(stat.Mode())
				// 	// Serve Directories
				// 	if stat.IsDir() {
				// 		fis, err := os.ReadDir(absPath)
				// 		if err != nil {
				// 			return err
				// 		}
				// 		entry.Entry(fis...)
				// 		return nil
				// 	}
				// 	// Serve Files
				// 	data, err := os.ReadFile(absPath)
				// 	if err != nil {
				// 		return err
				// 	}
				// 	entry.Write(data)
				// 	return nil
				// }))
			}
			return nil
		}))
	}
	return nil
}