package persistent

import (
	"fmt"
	"os"
	"path/filepath"

	utils "SQLite-GUI/internal/utils"
)

type ProfileCat struct {
	Paths           []string `mapstructure:"Paths" yaml:"Paths"`
	LastProfileUsed int      `mapstructure:"Last_Profile_Used" yaml:"Last_Profile_Used"`
	LastTargetUsed  int      `mapstructure:"Last_Target_Used" yaml:"Last_Target_Used"`
}

type DataType struct {
	Profiles ProfileCat `mapstructure:"Profiles" yaml:"Profiles"`
}

var Data DataType

func DefData() DataType {
	return DataType{
		Profiles: ProfileCat{
			Paths:           []string{},
			LastProfileUsed: 0,
		},
	}
}

func LoadData() error {
	dataLocation := utils.DataLoc()
	if !utils.CheckPath(dataLocation) {
		if err := os.MkdirAll(filepath.Dir(dataLocation), os.ModePerm); err != nil {
			return fmt.Errorf("error creating data directory: %v", err)
		}
		Data = DefData()
		if err := SaveData(); err != nil {
			return fmt.Errorf("error saving default data file: %v", err)
		}
		fmt.Println("Default data created at", dataLocation)
		return nil
	} else {
		fmt.Printf("Loading data from %s\n", dataLocation)
	}

	Data = DefData()
	if err := utils.LoadYamlFile(dataLocation, &Data); err != nil {
		return fmt.Errorf("error loading data file: %v", err)
	}

	return nil
}

func SaveData() error {
	if err := utils.SaveYamlFile(utils.DataLoc(), &Data); err != nil {
		return fmt.Errorf("error saving data file: %v", err)
	}
	return nil
}
