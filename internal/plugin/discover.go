package plugin

import (
	"os"
	"path/filepath"
	"strings"
)

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular() && (info.Mode().Perm()&0111 != 0)
}

func ScanPathForPlugins() ([]string, error) {
	var plugins []string
	var dirsToScan []string

	if cacheDir != "" && cacheDir != "." {
		dirsToScan = append(dirsToScan, cacheDir)
	}

	pathEnv := os.Getenv("PATH")
	if pathEnv != "" {
		pathDirs := strings.Split(pathEnv, string(os.PathListSeparator))
		dirsToScan = append(dirsToScan, pathDirs...)
	}

	seen := make(map[string]bool)
	for _, dir := range dirsToScan {
		if seen[dir] {
			continue
		}
		seen[dir] = true

		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasPrefix(entry.Name(), "yst-") {
				continue
			}

			fullPath := filepath.Join(dir, entry.Name())
			if isExecutable(fullPath) {
				plugins = append(plugins, fullPath)
			}
		}
	}

	return plugins, nil
}

