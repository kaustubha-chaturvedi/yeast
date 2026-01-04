package cli

import (
	"os"

	"github.com/kaustubha-chaturvedi/yeast/internal/plugin"
	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
)

func HandlePluginCommand(args []string) {
	if len(args) == 0 {
		utils.Printf("[ERROR] No alias specified\n")
		os.Exit(1)
	}

	alias, rest := args[0], args[1:]
	if len(rest) == 0 {
		pluginPath, err := plugin.FindPlugin(alias)
		if err != nil {
			utils.Printf("[ERROR] Plugin not found: %s\n", alias)
			os.Exit(1)
		}
		plugin.Execute(pluginPath, []string{"-h"})
		return
	}

	if err := Route(alias, rest); err != nil {
		os.Exit(1)
	}
}

func Route(alias string, args []string) error {
	pluginPath, err := plugin.FindPlugin(alias)
	if err != nil {
		utils.Printf("[ERROR] Plugin not found: %s\n", alias)
		return err
	}
	return utils.HandleError("execute plugin", plugin.Execute(pluginPath, args))
}
