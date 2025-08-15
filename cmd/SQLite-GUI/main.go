package main

import (
	"fmt"
	"os"

	Content "SQLite-GUI/internal/content"
	Persistent "SQLite-GUI/internal/persistent"
)

func main() {
	setup()

	Content.Init()

	shutdown()
}

func setup() {
	// load config
	if err := Persistent.LoadConfig(); err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	// load user data
	if err := Persistent.LoadData(); err != nil {
		fmt.Println("Error loading data:", err)
		os.Exit(1)
	}
	// load profiles
	Persistent.LoadProfiles()
}

func shutdown() {

}
