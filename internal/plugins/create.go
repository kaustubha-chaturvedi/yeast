package plugins

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
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
			return utils.HandleError("read template", err)
		}

		content := string(data)
		content = strings.ReplaceAll(content, "{{NAME}}", name)
		content = strings.ReplaceAll(content, "{{DOMAIN}}", domain)
		content = strings.ReplaceAll(content, "{{ALIAS}}", alias)

		relPath, _ := filepath.Rel(templatePath, path)
		relPath = strings.TrimSuffix(relPath, ".template")
		targetPath := filepath.Join(targetDir, relPath)

		if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
			return utils.HandleError("write template", err)
		}

		return nil
	})
}


func CreateNew(name, alias, domain, targetDir string) error {
	if name == "" {
		return utils.HandleErrorf("create", "name is required")
	}
	if alias == "" {
		return utils.HandleErrorf("create", "alias is required")
	}
	if domain == "" {
		return utils.HandleErrorf("create", "domain is required")
	}

	if targetDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return utils.HandleError("get working directory", err)
		}
		targetDir = filepath.Join(cwd, "yst-"+name)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return utils.HandleError("create directory", err)
	}

	if err := copyTemplates(targetDir, name, domain, alias); err != nil {
		return err
	}

	utils.Logf("Created plugin skeleton: %s", targetDir)
	return nil
}

