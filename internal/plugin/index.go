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

	if os.Getenv("PATH") == "" {
		indexBuilt = true
		return nil
	}

	cached, _ := loadCached()
	if cached == nil {
		cached = make(map[string]string)
	}
	
	for alias, path := range cached {
		if _, err := os.Stat(path); err == nil {
			if meta, err := GetMetadata(path); err == nil && validatePlugin(meta) {
				aliasIndex[alias] = path
			}
		}
	}

	pluginPaths, err := ScanPathForPlugins()
	if err != nil {
		return err
	}

	needsSave := false
	for _, fullPath := range pluginPaths {
		meta, err := GetMetadata(fullPath)
		if err != nil || !validatePlugin(meta) {
			continue
		}

		cachedPath, inCache := cached[meta.Alias]
		if !inCache || cachedPath != fullPath {
			needsSave = true
		}

		aliasIndex[meta.Alias] = fullPath
	}

	if needsSave {
		saveCached(aliasIndex)
	}

	indexBuilt = true
	return nil
}

func validatePlugin(meta *PluginMetadata) bool {
	return meta.Alias != "" && meta.Domain != "" && meta.Name != ""
}

func InvalidateIndex() {
	indexMutex.Lock()
	defer indexMutex.Unlock()
	indexBuilt = false
	aliasIndex = make(map[string]string)
	if err := os.Remove(getCachePath()); err != nil {
		utils.Printf("[ERROR] Failed to invalidate index: %v\n", err)
	}
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

	path, ok := aliasIndex[alias]
	if !ok {
		return "", &PluginNotFoundError{Alias: alias}
	}

	if absPath, err := filepath.Abs(path); err == nil {
		return absPath, nil
	}
	return path, nil
}

type PluginNotFoundError struct {
	Alias string
}

func (e *PluginNotFoundError) Error() string {
	return "plugin not found: " + e.Alias
}

func GetAliasIndex() map[string]string {
	indexMutex.RLock()
	defer indexMutex.RUnlock()
	
	result := make(map[string]string)
	for alias, path := range aliasIndex {
		result[alias] = path
	}
	return result
}

