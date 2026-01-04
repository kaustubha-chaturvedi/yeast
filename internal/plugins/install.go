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
)

type installResponse struct {
	DownloadURL string `json:"downloadUrl"`
	Error       string `json:"error,omitempty"`
}

func Install(pluginTag string) error {
	parts := strings.Split(pluginTag, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("install: invalid format. Use 'author:plugin-alias'")
	}

	author, pluginAlias := parts[0], parts[1]
	apiURL := fmt.Sprintf("https://yeast.kaustubha.work/api/download/%s:%s?os=%s&arch=%s", author, pluginAlias, runtime.GOOS, runtime.GOARCH)

	fmt.Printf("Fetching plugin info for %s...\n", pluginTag)
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("install: failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("install: plugin not found (status: %d)", resp.StatusCode)
	}

	var apiResp installResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("install: invalid API response: %w", err)
	}

	if apiResp.Error != "" {
		return fmt.Errorf("install: API error: %s", apiResp.Error)
	}
	if apiResp.DownloadURL == "" {
		return fmt.Errorf("install: no download URL")
	}
	urlParts := strings.Split(apiResp.DownloadURL, "/")
	fileName := urlParts[len(urlParts)-1]
	if fileName == "" {
		return fmt.Errorf("install: invalid download URL")
	}

	pluginPath := filepath.Join(plugin.GetBinaryDir(), fileName)
	fmt.Printf("Downloading to %s...\n", pluginPath)

	if err := downloadFile(apiResp.DownloadURL, pluginPath); err != nil {
		return fmt.Errorf("install: download failed: %w", err)
	}

	if err := os.Chmod(pluginPath, 0755); err != nil {
		return fmt.Errorf("install: chmod failed: %w", err)
	}

	fmt.Printf("Installed: %s\n", pluginPath)
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
