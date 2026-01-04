package plugin

import (
	"encoding/json"
	"os/exec"

	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
)

type PluginMetadata struct {
	Name     string            `json:"name"`
	Domain   string            `json:"domain"`
	Alias    string            `json:"alias,omitempty"`
	Version  string            `json:"version,omitempty"`
	Commands []CommandMetadata `json:"commands"`
}

type CommandMetadata struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Flags       []FlagMetadata    `json:"flags,omitempty"`
	Examples    []string          `json:"examples,omitempty"`
}

type FlagMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type,omitempty"`
	Default     string `json:"default,omitempty"`
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

