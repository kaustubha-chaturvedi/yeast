package plugin

import (
	"encoding/json"
	"fmt"
	"os/exec"
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
		return nil, fmt.Errorf("get metadata (%s): %w", pluginPath, err)
	}

	var metadata PluginMetadata
	if err := json.Unmarshal(output, &metadata); err != nil {
		return nil, fmt.Errorf("parse metadata (%s): %w", pluginPath, err)
	}

	return &metadata, nil
}

