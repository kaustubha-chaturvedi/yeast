package cli

import (
	"encoding/json"

	"github.com/kaustubha-chaturvedi/yeast/internal/plugin"
	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
)

func ShowMetadata(domain string) error {
	pluginPath, err := plugin.FindPlugin(domain)
	if err != nil {
		return utils.HandleErrorf("show metadata", "plugin not found: %s", domain)
	}

	meta, err := plugin.GetMetadata(pluginPath)
	if err != nil {
		return utils.HandleError("get metadata", err)
	}

	output, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return utils.HandleError("marshal metadata", err)
	}

	utils.Log(string(output))
	return nil
}
