package main

import (
	"log"

	"github.com/stephen-fox/steamutil/locations"
)

func main() {
	isInstalled := locations.IsInstalled()
	if isInstalled {
		log.Println("Steam is installed")
	} else {
		log.Println("Steam is *not* installed")
	}
}
