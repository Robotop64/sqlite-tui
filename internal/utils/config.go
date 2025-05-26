package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

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

	cfg.SetDefault("profiles.paths", []string{""})
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
	profile := cfg.New()

	temp := strings.Split(path, string(os.PathSeparator))
	filename := temp[len(temp)-1]
	profile.SetConfigName(filename)
	profile.SetConfigType("yaml")
	profile.AddConfigPath(path)

	profile.SetDefault("profile.name", "PROFILE_NAME")
	profile.SetDefault("database.path", "DATABASE_PATH")
	profile.SetDefault("database.type", "sqlite")
	profile.SetDefault("scripts.paths", []string{""})

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("profile path does not exist: %s", path)
	}
	if err := profile.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading profile config: %v", err)
	}

	return profile, nil
}

func GenProfile(path string) Profile {
	profile := cfg.New()

	temp := strings.Split(path, string(os.PathSeparator))
	err := os.MkdirAll(strings.Join(temp[:len(temp)-1], string(os.PathSeparator)), os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating path to profile config: %v", err)
	}

	filename := temp[len(temp)-1]

	profile.SetConfigName(filename)
	profile.SetConfigType("yaml")
	profile.AddConfigPath(path)

	profile.SetDefault("profile.name", "PROFILE_NAME")
	profile.SetDefault("database.path", "DATABASE_PATH")
	profile.SetDefault("database.type", "sqlite")
	profile.SetDefault("scripts.paths", []string{""})

	//check path and write config if the path does exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := profile.SafeWriteConfig(); err != nil {
			log.Fatalf("Error writing profile config: %v", err)
		}
	}

	return profile
}
