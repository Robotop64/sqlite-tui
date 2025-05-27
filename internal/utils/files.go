package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func CleanPath(path string) string {
	path = strings.Trim(path, "\" ")
	path = filepath.Clean(path)
	return path
}

func CheckPath(path string) bool {
	path = CleanPath(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func FileFromPath(path string, extension bool) string {
	filename := strings.Replace(path, filepath.Dir(path), "", 1)

	if !extension {
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	}

	return filename
}
