package plugins

import (
	"fmt"
	"sort"

	"github.com/kaustubha-chaturvedi/yeast/internal/plugin"
)

func List() error {
	plugin.InvalidateIndex()
	if err := plugin.BuildIndex(); err != nil {
		fmt.Printf("[ERROR] Failed to build index: %v\n", err)
		return err
	}

	type pluginInfo struct {
		alias string
		path  string
	}
	var plugins []pluginInfo

	plugin.ForEachPlugin(func(alias, path string) {
		plugins = append(plugins, pluginInfo{alias, path})
	})

	if len(plugins) == 0 {
		fmt.Printf("No plugins found in PATH\n")
		return nil
	}

	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].alias < plugins[j].alias
	})

	fmt.Printf("Found %d plugin(s)\n", len(plugins))
	for _, p := range plugins {
		if meta, err := plugin.GetMetadata(p.path); err == nil {
			line := "  " + p.alias
			if meta.Domain != "" {
				line += " (domain: " + meta.Domain + ")"
			}
			fmt.Printf("%s\n", line)
			fmt.Printf("Path: %s\n", p.path)
		} else {
			fmt.Printf("  %s - %s\n", p.alias, p.path)
		}
	}
	return nil
}

