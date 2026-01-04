package plugin

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
)

var (
	aliasIndex     = make(map[string]string)
	indexBuilt     bool
	indexMutex     sync.RWMutex
	cacheDir       string
)

func init() {
	execPath, err := os.Executable()
	if err != nil {
		cacheDir = "."
		return
	}
	cacheDir = utils.GetDir(execPath)
}

func BuildIndex() error {
	indexMutex.Lock()
	defer indexMutex.Unlock()

	aliasIndex = make(map[string]string)
	pluginPaths, err := ScanPathForPlugins()
	if err != nil {
		return err
	}

	for _, path := range pluginPaths {
		if meta, err := GetMetadata(path); err == nil && meta.Alias != "" && meta.Domain != "" && meta.Name != "" {
			aliasIndex[meta.Alias] = path
		}
	}

	indexBuilt = true
	return nil
}

func InvalidateIndex() {
	indexMutex.Lock()
	defer indexMutex.Unlock()
	indexBuilt = false
	aliasIndex = make(map[string]string)
}

func FindPlugin(alias string) (string, error) {
	indexMutex.RLock()
	if !indexBuilt {
		indexMutex.RUnlock()
		if err := BuildIndex(); err != nil {
			return "", utils.HandleError("build index", err)
		}
		indexMutex.RLock()
	}
	defer indexMutex.RUnlock()

	if path, ok := aliasIndex[alias]; ok {
		if absPath, err := filepath.Abs(path); err == nil {
			return absPath, nil
		}
		return path, nil
	}
	return "", &PluginNotFoundError{Alias: alias}
}

type PluginNotFoundError struct {
	Alias string
}

func (e *PluginNotFoundError) Error() string {
	return "plugin not found: " + e.Alias
}

func ForEachPlugin(fn func(alias, path string)) {
	indexMutex.RLock()
	defer indexMutex.RUnlock()
	for alias, path := range aliasIndex {
		fn(alias, path)
	}
}

func GetBinaryDir() string {
	return cacheDir
}

