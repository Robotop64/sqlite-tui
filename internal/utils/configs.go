package utils

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	Profiles ProfileCat `mapstructure:"Profiles" yaml:"Profiles"`
}

var Configs Config

var configLocation string

func DefConfig() Config {
	return Config{
		Profiles: ProfileCat{
			Paths:    []string{},
			LastUsed: 0,
		},
	}
}

func LoadConfig() error {
	configDir := ConfigDir()

	os.MkdirAll(configDir, os.ModePerm)

	if !CheckPath(configLocation) {
		fmt.Printf("Config file not found at %s. \nCreating default.\n", configLocation)
		if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
			log.Fatalf("Error creating config directory: %v", err)
		}
		Configs = DefConfig()
		if err := SaveConfig(); err != nil {
			log.Fatalf("Error saving default config file: %v", err)
		}
		fmt.Println("Default config created at", configLocation)
		return nil
	} else {
		fmt.Printf("Loading config from %s\n", configLocation)
	}

	Configs = DefConfig()
	if err := LoadYamlFile(configLocation, &Configs); err != nil {
		fmt.Printf("Error loading config file: %v\n", err)
		overwriteDialog(configLocation)
		return nil
	}

	return nil
}

func overwriteDialog(path string) {
	fmt.Println("Overwrite existing config? (y/n)")
	for {
		var response string
		fmt.Scanln(&response)
		switch response {
		case "y", "Y":
			fmt.Println("Creating new config file...")
			Configs = DefConfig()
			if err := SaveConfig(); err != nil {
				log.Fatalf("Error saving new config file: %v", err)
			}
			fmt.Println("New config file created at", path)
			return
		case "n", "N":
			log.Fatalf("Exiting without changes to config file.")
		default:
			fmt.Println("Please enter 'y' or 'n'.")
			continue
		}
	}
}

func SaveConfig() error {
	if err := SaveYamlFile(configLocation, Configs); err != nil {
		return fmt.Errorf("error saving config YAML file: %v", err)
	}

	return nil
}
