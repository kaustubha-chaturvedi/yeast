package plugin

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
)

func getCachePath() string {
	return filepath.Join(cacheDir, "yst.index.enc")
}

func saveCached(aliasToPath map[string]string) error {
	var lines []string
	for alias, path := range aliasToPath {
		lines = append(lines, alias+":"+path)
	}
	data := []byte(strings.Join(lines, "\n"))

	encrypted, err := utils.Encrypt(data)
	if err != nil {
		return utils.HandleError("encrypt cache", err)
	}

	return os.WriteFile(getCachePath(), encrypted, 0600)
}

func loadCached() (map[string]string, error) {
	cachePath := getCachePath()
	
	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, utils.HandleError("read cache", err)
	}

	decrypted, err := utils.Decrypt(data)
	if err != nil {
		os.Remove(cachePath)
		return nil, utils.HandleError("decrypt cache", err)
	}

	result := make(map[string]string)
	lines := strings.Split(string(decrypted), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}

	return result, nil
}
