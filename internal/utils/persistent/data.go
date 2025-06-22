package persistent

import (
	"fmt"
	"os"

	utils "github.com/Robotop64/sqlite-tui/internal/utils"
)

type ProfileCat struct {
	Paths    []string `mapstructure:"Paths" yaml:"Paths"`
	LastUsed int      `mapstructure:"Last_used" yaml:"Last_used"`
}

type DataType struct {
	Profiles ProfileCat `mapstructure:"Profiles" yaml:"Profiles"`
}

var Data DataType

func DefData() DataType {
	return DataType{
		Profiles: ProfileCat{
			Paths:    []string{},
			LastUsed: 0,
		},
	}
}

func LoadData() error {
	dataLocation := utils.DataLoc()
	if !utils.CheckPath(dataLocation) {
		if err := os.MkdirAll(dataLocation, os.ModePerm); err != nil {
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
