package plugin

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
)

func ScanPathForPlugins() ([]string, error) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return nil, nil
	}

	var plugins []string
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, dir := range paths {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasPrefix(entry.Name(), "yst-") {
				continue
			}

			fullPath := filepath.Join(dir, entry.Name())
			if utils.IsExecutable(fullPath) {
				plugins = append(plugins, fullPath)
			}
		}
	}

	return plugins, nil
}

