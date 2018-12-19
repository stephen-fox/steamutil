package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/stephen-fox/steamutil/shortcuts"
)

func main() {
	filePath := flag.String("f", "", "The path to the shortcuts file to modify")
	gameName := flag.String("n", "", "The name of the game to modify")
	flag.Parse()

	if len(os.Args) == 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if len(strings.TrimSpace(*filePath)) == 0 {
		log.Fatal("Please specify a shortcuts file")
	}

	if len(strings.TrimSpace(*gameName)) == 0 {
		log.Fatal("Please specify a game name to modify")
	}

	onMatch := func(name string, match *shortcuts.Shortcut) {
		now := time.Now()

		log.Println("Setting last played time for", name, "to", now)

		match.LastPlayTimeEpoch = int32(now.Unix())
	}

	noMatch := func(s string) (shortcuts.Shortcut, bool) {
		log.Println("No match found for:", *gameName, "- Creating empty shortcut")

		return shortcuts.Shortcut{
			AppName: *gameName,
		}, false
	}

	config := shortcuts.CreateOrUpdateConfig{
		Path:      *filePath,
		MatchName: *gameName,
		OnMatch:   onMatch,
		NoMatch:   noMatch,
	}

	result, err := shortcuts.CreateOrUpdateFile(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Result:", result)
}
