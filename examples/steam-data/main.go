package main

import (
	"log"

	"github.com/stephen-fox/steamutil/locations"
)

func main() {
	dv, err := locations.NewDataVerifier()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get the root data directory path:
	log.Println("Root data path:", dv.RootDirPath())

	// Get the user data path:
	userDataPath, info, err := dv.UserDataDirPath()
	if err != nil {
		log.Fatal("Failed to get user data directory path - " + err.Error())
	}
	log.Println("User data directory path:", userDataPath,
		"- Mod time", info.ModTime())

	// Get Steam user IDs and their directory paths:
	userIdsToDataDirPaths, err := dv.UserIdsToDataDirPaths()
	if err != nil {
		log.Fatal("Failed to get user IDs - " + err.Error())
	}
	for id, dirPath := range userIdsToDataDirPaths {
		log.Println("ID:", id, "- Location:", dirPath)

		// Get the custom shortcuts file for a given user ID:
		shortcutPath, info, err := dv.ShortcutsFilePath(id)
		if err != nil {
			log.Println("Failed to get shortcut file path for " +
				id + " - " + err.Error())
		} else {
			log.Println("Shortcut file path for ID", id + ":" +
				shortcutPath, "- Size:", info.Size())
		}
	}
}
