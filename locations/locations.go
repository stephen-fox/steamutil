package locations

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

const (
	userDataDirName   = "userdata"
	userConfigDirName = "config"
	shortcutsFileName = "shortcuts.vdf"
)

// DataVerifier gets and verifies file and directory paths to data-related
// Steam locations.
type DataVerifier interface {
	UserDataDirPath() (string, os.FileInfo, error)
	UserIdsToDataDirPaths() (map[string]string, error)
	ShortcutsFilePath(userId string) (string, os.FileInfo, error)
}

type defaultDataVerifier struct {
	dataDir string
}

// ShortcutsFilePath returns the path to the shortcuts file for a given Steam
// user ID.
func (o defaultDataVerifier) ShortcutsFilePath(userId string) (string, os.FileInfo, error) {
	filePath := ShortcutsFilePath(o.dataDir, userId)

	i, err := os.Stat(filePath)
	if err != nil {
		return "", nil, err
	}

	return filePath, i, nil
}

// UserIdsToDataDirPaths returns a map of local Steam user IDs to their data
// storage directories.
func (o defaultDataVerifier) UserIdsToDataDirPaths() (map[string]string, error) {
	idsToDirs := make(map[string]string)

	dir, _, err := o.UserDataDirPath()
	if err != nil {
		return idsToDirs, err
	}

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return idsToDirs, err
	}

	for _, in := range infos {
		idsToDirs[in.Name()] = UserIdDirPath(o.dataDir, in.Name())
	}

	return idsToDirs, nil
}

// UserDataDirPath returns the path to the user data directory.
func (o defaultDataVerifier) UserDataDirPath() (string, os.FileInfo, error) {
	dirPath := UserDataDirPath(o.dataDir)

	i, err := os.Stat(dirPath)
	if err != nil {
		return "", nil, err
	}

	return dirPath, i, nil
}

// ShortcutsFilePath generates a path to the shortcuts file for the specified
// data directory and Steam user ID.
func ShortcutsFilePath(dataDirPath string, userId string) string {
	return path.Join(dataDirPath, userDataDirName, userId, userConfigDirName, shortcutsFileName)
}

// UserIdDirPath generates a path to the specified Steam user
// ID's directory.
func UserIdDirPath(dataDirPath string, userId string) string {
	return path.Join(UserDataDirPath(dataDirPath), userId)
}

// UserDataDirPath generates a user data directory path based on
// the specified data directory path.
func UserDataDirPath(dataDirPath string) string {
	return path.Join(dataDirPath, userDataDirName)
}

// NewDataVerifier creates a DataVerifier for getting and verifying data
// related file and directory locations.
func NewDataVerifier() (DataVerifier, error) {
	dirPath, _, err := DataDirPath()
	if err != nil {
		return &defaultDataVerifier{}, err
	}

	return &defaultDataVerifier{
		dataDir: dirPath,
	}, nil
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
