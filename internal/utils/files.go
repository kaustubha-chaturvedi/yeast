package utils

import (
	"os"
	"path/filepath"
)

func IsExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular() && (info.Mode().Perm()&0111 != 0)
}

func GetDir(path string) string {
	return filepath.Dir(path)
}

