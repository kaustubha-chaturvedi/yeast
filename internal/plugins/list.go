package plugins

import (
	"sort"

	"github.com/kaustubha-chaturvedi/yeast/internal/plugin"
	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
)

func List() error {
	plugin.InvalidateIndex()
	if err := plugin.BuildIndex(); err != nil {
		utils.Printf("[ERROR] Failed to build index: %v\n", err)
		return err
	}

	index := plugin.GetAliasIndex()
	if len(index) == 0 {
		utils.Printf("No plugins found in PATH\n")
		return nil
	}

	aliases := make([]string, 0, len(index))
	for alias := range index {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)

	utils.Printf("Found %d plugin(s)\n", len(aliases))
	for _, alias := range aliases {
		path := index[alias]
		if meta, err := plugin.GetMetadata(path); err == nil {
			line := "  " + alias
			if meta.Domain != "" {
				line += " (domain: " + meta.Domain + ")"
			}
			utils.Printf("%s\n", line)
			utils.Printf("    Path: %s\n", path)
		} else {
			utils.Printf("  %s - %s\n", alias, path)
		}
	}
	return nil
}

