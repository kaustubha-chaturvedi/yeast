package cli

import (
	"github.com/kaustubha-chaturvedi/yeast/internal/plugin"
	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
)

func HasHelpCommand(alias string) bool {
	pluginPath, err := plugin.FindPlugin(alias)
	if err != nil {
		return false
	}
	
	meta, err := plugin.GetMetadata(pluginPath)
	if err != nil {
		return false
	}
	
	for _, cmd := range meta.Commands {
		if cmd.Name == "help" {
			return true
		}
	}
	return false
}

func Route(domain string, args []string) error {
	pluginPath, err := plugin.FindPlugin(domain)
	if err != nil {
		return utils.HandleErrorf("route", "plugin not found: %s", domain)
	}

	if len(args) == 0 {
		return utils.HandleErrorf("route", "no command specified")
	}

	if args[0] == "metadata" {
		return nil
	}

	return utils.HandleError("execute plugin", plugin.Execute(pluginPath, args))
}
