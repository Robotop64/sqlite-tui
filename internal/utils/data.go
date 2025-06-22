package utils

import (
	"fmt"
	"os"
)

type ProfileCat struct {
	Paths    []string `mapstructure:"Paths" yaml:"Paths"`
	LastUsed int      `mapstructure:"Last_used" yaml:"Last_used"`
}

type DataType struct {
	Profiles ProfileCat `mapstructure:"Profiles" yaml:"Profiles"`
}

var Data DataType
var dataLocation string

func DefData() DataType {
	return DataType{
		Profiles: ProfileCat{
			Paths:    []string{},
			LastUsed: 0,
		},
	}
}

func LoadData() error {
	dataLocation = DataLoc()
	if !CheckPath(dataLocation) {
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
	if err := LoadYamlFile(dataLocation, &Data); err != nil {
		return fmt.Errorf("error loading data file: %v", err)
	}

	return nil
}

func SaveData() error {
	if err := SaveYamlFile(DataLoc(), &Data); err != nil {
		return fmt.Errorf("error saving data file: %v", err)
	}
	return nil
}
