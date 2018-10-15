package locations

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

const (
	userdataDirName   = "userdata"
	shortcutsFileName = "shortcuts.vdf"
)

// ShortcutsFilePath returns the path to the shortcuts file for a given Steam
// user ID.
func ShortcutsFilePath(userId string) (string, os.FileInfo, error) {
	userIdsToDirPaths, err := UserIdsToDataDirPaths()
	if err != nil {
		return "", nil, err
	}

	dirPath, ok := userIdsToDirPaths[userId]
	if !ok {
		return "", nil, errors.New("The specified user ID does not exist")
	}

	filePath := path.Join(dirPath, shortcutsFileName)

	i, err := os.Stat(filePath)
	if err != nil {
		return "", nil, err
	}

	return filePath, i, nil
}

// UserIdsToDataDirPaths returns a map of local Steam user IDs to their data
// storage directories.
func UserIdsToDataDirPaths() (map[string]string, error) {
	m := make(map[string]string)

	dir, _, err := UserDataDirPath()
	if err != nil {
		return m, err
	}

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return m, err
	}

	for _, in := range infos {
		m[in.Name()] = path.Join(dir, in.Name())
	}

	return m, nil
}

// UserDataDirPath returns the path to the user data directory.
func UserDataDirPath() (string, os.FileInfo, error) {
	data, _, err := DataDirPath()
	if err != nil {
		return "", nil, err
	}

	dirPath := path.Join(data, userdataDirName)

	i, err := os.Stat(dirPath)
	if err != nil {
		return "", nil, err
	}

	return dirPath, i, nil
}

// IsInstalled returns true if Steam is installed.
func IsInstalled() bool {
	_, _, err := DataDirPath()
	if err != nil {
		return false
	}

	return true
}

// DataDirPath returns the path to Steam's data directory.
func DataDirPath() (string, os.FileInfo, error) {
	dirPath := ""

	switch operatingSystem := runtime.GOOS; operatingSystem {
	case "darwin":
		dirPath = os.Getenv("HOME") + "/Library/Application Support/Steam"
	case "linux":
		return "", nil, errors.New("Linux is not currently supported :(")
	case "windows":
		return "", nil, errors.New("Windows is not currently supported :(")

	}

	i, err := os.Stat(dirPath)
	if err != nil {
		return dirPath, nil, err
	}

	return dirPath, i, nil
}
