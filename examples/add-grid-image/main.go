package main

import (
	"flag"
	"log"
	"os"

	"github.com/stephen-fox/steamutil/grid"
	"github.com/stephen-fox/steamutil/locations"
)

func main() {
	imagePath := flag.String("i", "", "The path to the image to add")
	gameName := flag.String("g", "", "The name of the game that the image is for")
	gameExePath := flag.String("e", "", "The game's executable path")

	flag.Parse()

	if len(os.Args) == 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	dv, err := locations.NewDataVerifier()
	if err != nil {
		log.Fatal(err.Error())
	}

	resultDetails := grid.ImageDetails{
		DataVerifier:       dv,
		GameName:           *gameName,
		GameExecutablePath: *gameExePath,
	}

	userIdsToDirPaths, err := dv.UserIdsToDataDirPaths()
	if err != nil {
		log.Fatal(err.Error())
	}

	for userId := range userIdsToDirPaths {
		resultDetails.OwnerUserId = userId

		addConfig := grid.AddConfig{
			ResultDetails:   resultDetails,
			ImageSourcePath: *imagePath,
		}

		err := grid.AddImage(addConfig)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
