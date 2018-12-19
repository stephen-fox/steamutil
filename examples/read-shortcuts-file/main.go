package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/stephen-fox/steamutil/locations"
	"github.com/stephen-fox/steamutil/shortcuts"
)

func main() {
	dv, err := locations.NewDataVerifier()
	if err != nil {
		log.Fatal(err.Error())
	}

	userIdsToDirs, err := dv.UserIdsToDataDirPaths()
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(userIdsToDirs) == 0 {
		log.Fatal("No user data directories exist")
	}

	var shortcutFilePath string

	for id := range userIdsToDirs {
		var err error
		shortcutFilePath, _, err = dv.ShortcutsFilePath(id)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if shortcutFilePath == "" {
		log.Fatal("Could not locate a shortcuts file to read")
	}

	f, err := os.Open(shortcutFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	scs, err := shortcuts.ReadFile(f)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, s := range scs {
		fmt.Println("Application name:", s.AppName)
		fmt.Println("Start dir:", s.StartDir)
		fmt.Println("Application name:", s.ExePath)
		fmt.Println("Application name:", s.LaunchOptions)
		fmt.Println("Application name:", s.LastPlayTimeEpoch)
		fmt.Println(strings.Repeat("-", 80))
	}
}
