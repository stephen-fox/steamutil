package locations

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"
)

const (
	programFiles86 = programFiles + " (x86)"
	programFiles64 = programFiles
	programFiles   = "Program Files"
)

var (
	possibleDriveLetters = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
)

// DataDirPath returns the path to Steam's data directory.
func DataDirPath() (string, os.FileInfo, error) {
	bitmask, err := windows.GetLogicalDrives()
	if err != nil {
		return "", nil, err
	}

	driveLetters := bitsToDriveLetters(bitmask)

	for _, letter := range driveLetters {
		letter := letter + ":\\"

		target := filepath.Join(letter, programFiles86, "Steam")
		info, statErr := os.Stat(target)
		if statErr == nil {
			return target, info, nil
		}

		target = filepath.Join(letter, programFiles64, "Steam")
		info, statErr = os.Stat(target)
		if statErr == nil {
			return target, info, nil
		}
	}

	return "", nil, errors.New("Failed to locate Steam data directory on volumes " +
		strings.Join(driveLetters, " "))
}

// Based on work by "nemo": https://stackoverflow.com/a/23135463
func bitsToDriveLetters(bitMap uint32) []string {
	var drives []string

	for i := range possibleDriveLetters {
		if bitMap & 1 == 1 {
			drives = append(drives, possibleDriveLetters[i])
		}
		bitMap >>= 1
	}

	return drives
}
