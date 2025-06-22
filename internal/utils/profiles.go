package utils

import (
	"fmt"
	"path/filepath"
)

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

var Profiles []*Profile

func DefProfile() Profile {
	return Profile{
		Name:    "New Profile",
		Path:    "",
		Targets: []Target{},
		Note:    "",
	}
}

func LoadProfile(path string) (*Profile, error) {
	profile := DefProfile()

	if err := LoadYamlFile(path, &profile); err != nil {
		return nil, fmt.Errorf("error loading profile YAML file: %v", err)
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

	profile := DefProfile()

	profile.Path = path

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
