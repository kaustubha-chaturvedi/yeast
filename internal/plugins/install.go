package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kaustubha-chaturvedi/yeast/internal/plugin"
	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
)

type installResponse struct {
	DownloadURL string `json:"downloadUrl"`
	Error       string `json:"error,omitempty"`
}

func Install(pluginTag string) error {
	parts := strings.Split(pluginTag, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return utils.HandleErrorf("install", "invalid format. Use 'author:plugin-name'")
	}

	author, pluginName := parts[0], parts[1]
	apiURL := fmt.Sprintf("https://yeast.kaustubha.work/plugins/%s:%s?os=%s", author, pluginName, runtime.GOOS)

	utils.Printf("Fetching plugin info for %s...\n", pluginTag)
	resp, err := http.Get(apiURL)
	if err != nil {
		return utils.HandleErrorf("install", "failed to fetch: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return utils.HandleErrorf("install", "plugin not found (status: %d)", resp.StatusCode)
	}

	var apiResp installResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return utils.HandleErrorf("install", "invalid API response: %v", err)
	}

	if apiResp.Error != "" {
		return utils.HandleErrorf("install", "API error: %s", apiResp.Error)
	}
	if apiResp.DownloadURL == "" {
		return utils.HandleErrorf("install", "no download URL")
	}

	pluginPath := filepath.Join(plugin.GetBinaryDir(), fmt.Sprintf("yst-%s", pluginName))
	utils.Printf("Downloading to %s...\n", pluginPath)

	if err := downloadFile(apiResp.DownloadURL, pluginPath); err != nil {
		return utils.HandleErrorf("install", "download failed: %v", err)
	}

	if err := os.Chmod(pluginPath, 0755); err != nil {
		return utils.HandleErrorf("install", "chmod failed: %v", err)
	}

	utils.Printf("Installed: %s\n", pluginPath)
	plugin.InvalidateIndex()
	plugin.BuildIndex()
	return nil
}

func downloadFile(url, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}


