package plugins

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type publishRequest struct {
	ModuleName  string `json:"moduleName"`
	Name        string `json:"name"`
	Domain      string `json:"domain"`
	Alias       string `json:"alias"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type validAliasResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

type publishResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

const minWords, maxWords int = 25, 256

func Publish() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("publish: get working directory: %w", err)
	}

	moduleName, err := readModuleName(cwd)
	if err != nil {
		return fmt.Errorf("publish: read go.mod: %w", err)
	}

	meta, err := getMetadataFromSource(cwd)
	if err != nil {
		return fmt.Errorf("publish: get metadata: %w", err)
	}

		if err := validateMetadata(meta); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	fmt.Printf("Checking if alias '%s' is available...\n", meta.Alias)
	if err := checkAliasAvailability(meta.Alias); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	description, err := promptDescription()
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	fmt.Printf("Publishing plugin '%s'...\n", meta.Name)
	req := publishRequest{
		ModuleName:  moduleName,
		Name:        meta.Name,
		Domain:      meta.Domain,
		Alias:       meta.Alias,
		Version:     meta.Version,
		Description: description,
	}

	if err := sendPublishRequest(req); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	fmt.Printf("Successfully published plugin '%s' (alias: %s)\n", meta.Name, meta.Alias)
	return nil
}

func readModuleName(dir string) (string, error) {
	goModPath := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

func getMetadataFromSource(dir string) (*pluginMetadata, error) {
	mainPath := filepath.Join(dir, "main.go")
	if _, err := os.Stat(mainPath); err != nil {
		return nil, fmt.Errorf("main.go not found in current directory")
	}

	tempBin := filepath.Join(dir, ".yst-temp-plugin")
	if runtime.GOOS == "windows" {
		tempBin += ".exe"
	}
	defer func() {
		os.Remove(tempBin)
	}()

	cmd := exec.Command("go", "build", "-o", tempBin, ".")
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("build failed: %s", stderr.String())
	}

	metaCmd := exec.Command(tempBin, "__yst_metadata")
	output, err := metaCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("get metadata from binary: %w", err)
	}

	var meta pluginMetadata
	if err := json.Unmarshal(output, &meta); err != nil {
		return nil, fmt.Errorf("parse metadata: %w", err)
	}

	return &meta, nil
}

type pluginMetadata struct {
	Name    string `json:"name"`
	Domain  string `json:"domain"`
	Alias   string `json:"alias"`
	Version string `json:"version,omitempty"`
}

func validateMetadata(meta *pluginMetadata) error {
	if meta.Name == "" {
		return fmt.Errorf("metadata missing required field: name")
	}
	if meta.Domain == "" {
		return fmt.Errorf("metadata missing required field: domain")
	}
	if meta.Alias == "" {
		return fmt.Errorf("metadata missing required field: alias")
	}
	return nil
}

func checkAliasAvailability(alias string) error {
	url := fmt.Sprintf("https://yeast.kaustubha.work/valid-alias?alias=%s", alias)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("check alias availability: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("check alias availability: server returned status %d", resp.StatusCode)
	}

	var result validAliasResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	if !result.Valid {
		if result.Error != "" {
			return fmt.Errorf("alias unavailable: %s", result.Error)
		}
		return fmt.Errorf("alias '%s' is already taken", alias)
	}

	return nil
}

func sendPublishRequest(req publishRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	resp, err := http.Post("https://yeast.kaustubha.work/publish", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("publish failed (status %d): %s", resp.StatusCode, string(body))
	}

	var result publishResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	if !result.Success {
		if result.Error != "" {
			return fmt.Errorf("publish failed: %s", result.Error)
		}
		return fmt.Errorf("publish failed: unknown error")
	}

	return nil
}

func promptDescription() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter plugin description (at least 25 words): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("read description: %w", err)
		}

		description := strings.TrimSpace(input)
		words := strings.Fields(description)
		wordCount := len(words)

		if wordCount < minWords || maxWords < wordCount {
			fmt.Printf("description has only %d words. Please enter at between %d and %d words.", wordCount, minWords, maxWords)
			continue
		}
		
		return description, nil
	}
}
