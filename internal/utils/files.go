package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

const BIN_NAME = "sqlite-tui"

func DataLoc() string {
	var dirPrefix string

	switch runtime.GOOS {
	case "windows":
		dirPrefix = os.Getenv("LOCALAPPDATA")
	case "linux":
		dirPrefix, _ = os.UserHomeDir()
		dirPrefix = filepath.Join(dirPrefix, ".local", "share")
	default:
		log.Fatalf("Unsupported OS: %s", os.Getenv("OS"))
	}

	dataDir := filepath.Join(dirPrefix, BIN_NAME)
	if runtime.GOOS == "windows" {
		dataDir = filepath.Join(dataDir, "data")
	}

	return filepath.Join(dataDir, "userData.yaml")
}

func ConfigLoc() string {
	var configDirPrefix string

	switch runtime.GOOS {
	case "windows":
		configDirPrefix = os.Getenv("LOCALAPPDATA")
	case "linux":
		configDirPrefix, _ = os.UserConfigDir()
	default:
		log.Fatalf("Unsupported OS: %s", os.Getenv("OS"))
	}
	configDir := filepath.Join(configDirPrefix, BIN_NAME)

	return filepath.Join(configDir, "config.yaml")
}

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

func RelativeToAbsolutePath(root string, path string) string {
	if path[0] != '.' {
		return path
	}

	if !strings.HasSuffix(root, string(os.PathSeparator)) {
		root += string(os.PathSeparator)
	}
	if path[0:1] != ".." {
		path = path[2:]
	}
	path = CleanPath(path)
	return filepath.Join(root, path)
}

func SaveYamlFile(path string, data interface{}) error {
	path = CleanPath(path)
	if !CheckPath(filepath.Dir(path)) {
		fmt.Println("Creating directory:", filepath.Dir(path))
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}
	}

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data to YAML:", err)
		return err
	}
	if err := os.WriteFile(path, yamlData, 0644); err != nil {
		fmt.Println("Error writing YAML file:", err)
		return err
	}

	return nil
}

func LoadYamlFile(path string, data interface{}) error {
	path = CleanPath(path)
	if !CheckPath(path) {
		return os.ErrNotExist
	}

	var ydata []byte
	var err error
	if ydata, err = os.ReadFile(path); err != nil {
		return err
	}

	if err := yaml.Unmarshal(ydata, data); err != nil {
		return err
	}

	return nil
}

func OpenExternal(path string) error {
	path = CleanPath(path)
	if !CheckPath(path) {
		return fmt.Errorf("file does not exist: %s", path)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	return nil
}
