package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	cfg "github.com/spf13/viper"
	yaml "gopkg.in/yaml.v3"
)

const BIN_NAME = "sqlite-tui"

type ProfileCat struct {
	Paths    []string `mapstructure:"Paths" yaml:"Paths"`
	LastUsed int      `mapstructure:"Last_used" yaml:"Last_used"`
}

type Config struct {
	Profiles ProfileCat `mapstructure:"Profiles" yaml:"Profiles"`
}

type Target struct {
	Name         string   `mapstructure:"Name" yaml:"Name"`
	DatabasePath string   `mapstructure:"Database_Path" yaml:"Database_Path"`
	ScriptPaths  []string `mapstructure:"Script_Paths" yaml:"Script_Paths"`
}

type Profile struct {
	Name    string   `mapstructure:"Name" yaml:"Name"`
	Path    string   `mapstructure:"Path" yaml:"Path"`
	Targets []Target `mapstructure:"Targets" yaml:"Targets"`
	Note    string   `mapstructure:"Note" yaml:"Note"`
}

var Configs Config
var Profiles []*Profile

func DefConfig() Config {
	return Config{
		Profiles: ProfileCat{
			Paths:    []string{},
			LastUsed: 0,
		},
	}
}

func DefProfile() Profile {
	return Profile{
		Name:    "New Profile",
		Path:    "",
		Targets: []Target{},
		Note:    "",
	}
}

func structToCfg(c *cfg.Viper, data interface{}) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling struct to YAML: %v", err)
	}

	var m map[string]interface{}
	if err := yaml.Unmarshal(yamlData, &m); err != nil {
		return fmt.Errorf("error unmarshalling YAML to map: %v", err)
	}

	c.MergeConfigMap(m)

	return nil
}

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
	cfg.AddConfigPath(configDir)
	configLocation := filepath.Join(configDir, "config.yaml")

	if err := cfg.ReadInConfig(); err != nil {
		Configs = DefConfig()
		if _, ok := err.(cfg.ConfigFileNotFoundError); !ok {
			fmt.Printf("Config file not found at %s. \nCreating default.\n", configLocation)
			structToCfg(cfg.GetViper(), Configs)
			if err := cfg.SafeWriteConfig(); err != nil {
				log.Fatalf("Error writing default config file: %v", err)
			}
		} else {
			fmt.Printf("Error reading config file: %v\n", err)
		}
		fmt.Println("Using defaults.")
		return nil
	}

	if err := cfg.Unmarshal(&Configs); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	return nil
}

func SaveConfig() error {
	if err := structToCfg(cfg.GetViper(), Configs); err != nil {
		return fmt.Errorf("error converting config to viper: %v", err)
	}

	if err := cfg.WriteConfig(); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

func LoadProfile(path string) (*Profile, error) {
	path = CleanPath(path)
	if state := CheckPath(path); !state {
		return nil, fmt.Errorf("profile path does not exist: %s", path)
	}

	profileLoader := cfg.New()
	profile := DefProfile()

	profileLoader.SetConfigName(FileFromPath(path, false))
	profileLoader.AddConfigPath(filepath.Dir(path))
	profileLoader.SetConfigType("yaml")

	if err := profileLoader.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading profile config: %v", err)
	}
	if err := profileLoader.Unmarshal(&profile); err != nil {
		return nil, fmt.Errorf("error unmarshalling profile config: %v", err)
	}

	return &profile, nil
}

func SaveProfile(profile *Profile, path string) error {
	if err := SaveYamlFile(path, profile); err != nil {
		return fmt.Errorf("error saving profile YAML file: %v", err)
	}

	return nil
}

func CreateProfile(path string) (*Profile, error) {
	path = CleanPath(path)
	if !EndsWith(path, "Profile.yaml") {
		path = filepath.Join(path, "Profile.yaml")
	}

	if CheckPath(path) {
		return nil, fmt.Errorf("profile already exists at path: %s", path)
	}

	Configs.Profiles.Paths = append(Configs.Profiles.Paths, path)

	profile := DefProfile()
	profile.Name = FileFromPath(path, false)

	if err := SaveProfile(&profile, path); err != nil {
		return nil, fmt.Errorf("error saving new profile: %v", err)
	}

	return &profile, nil
}

func LoadProfiles() {
	paths := Configs.Profiles.Paths
	Profiles = make([]*Profile, len(paths))

	for i, path := range paths {
		profile, err := LoadProfile(path)
		if err != nil {
			Profiles[i] = nil
		} else {
			Profiles[i] = profile
		}
	}
}

func ActiveProfile() *Profile {

	if Configs.Profiles.LastUsed < 0 || Configs.Profiles.LastUsed >= len(Profiles) {
		return nil
	}

	return Profiles[Configs.Profiles.LastUsed]
}
