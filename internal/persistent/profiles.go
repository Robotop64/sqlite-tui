package persistent

import (
	"fmt"
	"path/filepath"

	utils "SQLite-GUI/internal/utils"
)

type Target struct {
	Name         string   `mapstructure:"Name" yaml:"Name"`
	DatabasePath string   `mapstructure:"Database_Path" yaml:"Database_Path"`
	ScriptPaths  []string `mapstructure:"Script_Paths" yaml:"Script_Paths"`
}

type Profile struct {
	Name    string   `mapstructure:"Name" yaml:"Name"`
	Targets []Target `mapstructure:"Targets" yaml:"Targets"`
	Note    string   `mapstructure:"Note" yaml:"Note"`
}

var Profiles []*Profile

func DefProfile() Profile {
	return Profile{
		Name:    "New Profile",
		Targets: []Target{},
		Note:    "",
	}
}

func LoadProfile(path string) (*Profile, error) {
	profile := DefProfile()

	if err := utils.LoadYamlFile(path, &profile); err != nil {
		return nil, fmt.Errorf("error loading profile YAML file: %v", err)
	}

	return &profile, nil
}

func SaveProfile(profile *Profile, path string) error {
	if err := utils.SaveYamlFile(path, profile); err != nil {
		return fmt.Errorf("error saving profile YAML file: %v", err)
	}

	return nil
}

func CreateProfile(path string) (*Profile, error) {
	path = utils.CleanPath(path)
	if !utils.EndsWith(path, "Profile.yaml") {
		path = filepath.Join(path, "Profile.yaml")
	}

	if utils.CheckPath(path) {
		return nil, fmt.Errorf("profile already exists at path: %s", path)
	}

	profile := DefProfile()

	if err := SaveProfile(&profile, path); err != nil {
		return nil, fmt.Errorf("error saving new profile: %v", err)
	}

	return &profile, nil
}

func LoadProfiles() {
	paths := Data.Profiles.Paths
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

	if Data.Profiles.LastProfileUsed < 0 || Data.Profiles.LastProfileUsed >= len(Profiles) {
		return nil
	}

	return Profiles[Data.Profiles.LastProfileUsed]
}

func ActiveProfilePath() string {
	prof := ActiveProfile()
	if prof == nil {
		return ""
	}
	return ProfilePath(prof)
}

func ProfilePath(profile *Profile) string {
	if profile == nil {
		return ""
	}

	for i := range len(Profiles) {
		iter_profile := Profiles[i]
		path := Data.Profiles.Paths[i]

		if iter_profile == profile {
			return path
		}
	}

	return ""
}
