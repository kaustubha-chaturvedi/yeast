package plugins

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)


func copyTemplates(targetDir, name, domain, alias string) error {
	_, filename, _, _ := runtime.Caller(0)
	templatePath := filepath.Join(filepath.Dir(filename), "templates")

	return filepath.WalkDir(templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read template: %w", err)
		}

		content := string(data)
		content = strings.ReplaceAll(content, "{{NAME}}", name)
		content = strings.ReplaceAll(content, "{{DOMAIN}}", domain)
		content = strings.ReplaceAll(content, "{{ALIAS}}", alias)

		relPath, _ := filepath.Rel(templatePath, path)
		relPath = strings.TrimSuffix(relPath, ".template")
		targetPath := filepath.Join(targetDir, relPath)

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("create directory: %w", err)
		}

		if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("write template: %w", err)
		}

		return nil
	})
}


func CreateNew(name, alias, domain, targetDir string) error {
	if name == "" {
		return fmt.Errorf("create: name is required")
	}
	if alias == "" {
		return fmt.Errorf("create: alias is required")
	}
	if domain == "" {
		return fmt.Errorf("create: domain is required")
	}

	if targetDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}
		targetDir = filepath.Join(cwd, "yst-"+alias)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	return copyTemplates(targetDir, name, domain, alias)
}

