package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	cfg "github.com/spf13/viper"
)

const BIN_NAME = "sqlite-tui"

func LoadConfig() error {
	cfg.SetConfigName("config")
	cfg.SetConfigType("yaml")

	var configDirPrefix string

	switch runtime.GOOS {
	case "windows":
		configDirPrefix = os.Getenv("LOCALAPPDATA")
	case "linux":
		configDirPrefix = "~/.config"
	// case "darwin":
	// 	cfg.AddConfigPath("/usr/local/etc/myapp")
	default:
		log.Fatalf("Unsupported OS: %s", os.Getenv("OS"))
	}
	configDir := filepath.Join(configDirPrefix, BIN_NAME)
	os.MkdirAll(configDir, os.ModePerm)
	configLocation := filepath.Join(configDir, "config.yaml")

	cfg.AddConfigPath(configDir)

	cfg.SetDefault("profiles.paths", []string{})
	cfg.SetDefault("profiles.last_used", 0)

	if _, err := os.Stat(configLocation); os.IsNotExist(err) {
		fmt.Println("Config file not found. \nCreating default.")

		if err := cfg.SafeWriteConfig(); err != nil {
			log.Fatalf("Error writing config file: %v", err)
		}
	} else {
		if err := cfg.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}
	}

	return nil
}

type Profile = *cfg.Viper

func GenProfile() Profile {
	profile := cfg.New()

	profile.SetConfigName("Profile")
	profile.SetConfigType("yaml")

	profile.SetDefault("profile.name", "New Profile")
	profile.SetDefault("database.paths", []string{})
	profile.SetDefault("scripts.paths", []string{})
	profile.SetDefault("note", "")

	return profile
}

func LoadProfile(path string) (Profile, error) {
	path = CleanPath(path)
	if state := CheckPath(path); !state {
		return nil, fmt.Errorf("profile path does not exist: %s", path)
	}

	profile := GenProfile()

	profile.SetConfigName(FileFromPath(path, false))
	profile.AddConfigPath(filepath.Dir(path))

	if err := profile.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading profile config: %v", err)
	}

	profile.Set("profile.path", path)

	return profile, nil
}

func LoadProfiles() []Profile {
	paths := cfg.GetStringSlice("profiles.paths")
	profiles := make([]Profile, len(paths))

	for i, path := range paths {
		profile, err := LoadProfile(path)
		if err != nil {
			profiles[i] = nil
		} else {
			profiles[i] = profile
		}
	}

	return profiles
}

func CreateProfile(path string) (Profile, error) {
	profile := GenProfile()

	path = CleanPath(path)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error creating profile directory: %v", err)
	}

	profile.AddConfigPath(path)

	if err := profile.SafeWriteConfig(); err != nil {
		return nil, fmt.Errorf("error writing profile config: %v", err)
	}

	profile.Set("profile.path", path)

	return profile, nil
}
