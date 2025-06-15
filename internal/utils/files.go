package utils

import (
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
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

func SaveYamlFile(path string, data interface{}) error {
	path = CleanPath(path)
	if !CheckPath(filepath.Dir(path)) {
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}
	}

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, yamlData, 0644); err != nil {
		return err
	}

	return nil
}
