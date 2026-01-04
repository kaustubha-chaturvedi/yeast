package plugin

import (
	"encoding/json"
	"os/exec"

	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
)

type PluginMetadata struct {
	Name    string `json:"name"`
	Domain  string `json:"domain"`
	Alias   string `json:"alias"`
	Version string `json:"version,omitempty"`
}

func GetMetadata(pluginPath string) (*PluginMetadata, error) {
	cmd := exec.Command(pluginPath, "__yst_metadata")
	output, err := cmd.Output()
	if err != nil {
		return nil, utils.WrapError("get metadata", err, pluginPath)
	}

	var metadata PluginMetadata
	if err := json.Unmarshal(output, &metadata); err != nil {
		return nil, utils.WrapError("parse metadata", err, pluginPath)
	}

	return &metadata, nil
}

