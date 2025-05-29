package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	cfg "github.com/spf13/viper"
)

func LoadConfig() error {
	cfg.SetConfigName("config")
	cfg.SetConfigType("yaml")

	cfg.AddConfigPath(".")

	// switch os.Getenv("OS") {
	// case "windows":
	// 	cfg.AddConfigPath("C:\\ProgramData\\MyApp")
	// case "linux":
	// 	cfg.AddConfigPath("/etc/myapp")
	// case "darwin":
	// 	cfg.AddConfigPath("/usr/local/etc/myapp")
	// default:
	// 	fmt.Println("Unsupported OS, using current directory for config.")
	// 	cfg.AddConfigPath(".")
	// }

	cfg.SetDefault("profiles.paths", []string{})
	cfg.SetDefault("profiles.last_used", 0)

	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
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

func LoadProfile(path string) (Profile, error) {
	path = CleanPath(path)
	if state := CheckPath(path); !state {
		return nil, fmt.Errorf("profile path does not exist: %s", path)
	}

	profile := cfg.New()

	profile.SetConfigName(FileFromPath(path, false))
	profile.SetConfigType("yaml")
	profile.AddConfigPath(filepath.Dir(path))

	profile.SetDefault("profile.name", "PROFILE_NAME")
	profile.SetDefault("profile.path", "PATH")
	profile.SetDefault("database.path", "DATABASE_PATH")
	profile.SetDefault("database.type", "sqlite")
	profile.SetDefault("scripts.paths", []string{})
	profile.SetDefault("note", "")

	if err := profile.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading profile config: %v", err)
	}

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

func GenProfile(path string) (Profile, error) {
	path = CleanPath(path)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error creating profile directory: %v", err)
	}

	profile := cfg.New()

	profile.SetConfigName("Profile")
	profile.SetConfigType("yaml")
	profile.AddConfigPath(path)

	profile.SetDefault("profile.name", "New Profile")
	profile.SetDefault("profile.path", "PATH")
	profile.SetDefault("database.path", "DATABASE_PATH")
	profile.SetDefault("database.type", "sqlite")
	profile.SetDefault("scripts.paths", []string{""})
	profile.SetDefault("note", "")

	if err := profile.SafeWriteConfig(); err != nil {
		return nil, fmt.Errorf("error writing profile config: %v", err)
	}

	return profile, nil
}
