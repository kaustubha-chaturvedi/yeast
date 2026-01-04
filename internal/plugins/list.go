package plugins

import (
	"sort"
	"strings"

	"github.com/kaustubha-chaturvedi/yeast/internal/plugin"
	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
)

type PluginWithPath struct {
	Metadata *plugin.PluginMetadata
	Path     string
}

func List() error {
	plugin.InvalidateIndex()
	if err := plugin.BuildIndex(); err != nil {
		return utils.HandleError("build index", err)
	}

	plugins := getAllPlugins()
	if len(plugins) == 0 {
		utils.Log("No plugins found in PATH")
		return nil
	}

	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Metadata.Alias < plugins[j].Metadata.Alias
	})
	printPlugins(plugins)
	return nil
}

func getAllPlugins() []PluginWithPath {
	var plugins []PluginWithPath
	for _, path := range plugin.GetAliasIndex() {
		if meta, err := plugin.GetMetadata(path); err == nil {
			plugins = append(plugins, PluginWithPath{Metadata: meta, Path: path})
		}
	}
	return plugins
}

func printPlugins(plugins []PluginWithPath) {
	utils.Logf("Found %d plugin(s)", len(plugins))
	for _, p := range plugins {
		meta := p.Metadata
		line := "  " + meta.Alias
		if meta.Domain != meta.Alias {
			line += " (domain: " + meta.Domain + ")"
		}
		utils.Log(line)
		if meta.Version != "" {
			utils.Logf("    Version: %s", meta.Version)
		}
		utils.Logf("    Path: %s", p.Path)
		if len(meta.Commands) > 0 {
			cmdNames := make([]string, len(meta.Commands))
			for i := range meta.Commands {
				cmdNames[i] = meta.Commands[i].Name
			}
			utils.Logf("    Commands: %s", strings.Join(cmdNames, ", "))
		}
		utils.Log("")
	}
}

