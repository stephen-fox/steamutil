package locations

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	userDataDirName   = "userdata"
	userConfigDirName = "config"
	shortcutsFileName = "shortcuts.vdf"
	gridDirName       = "grid"
)

// DataVerifier gets and verifies file and directory paths to data-related
// Steam locations.
type DataVerifier interface {
	// RootDirPath returns the path to the root Steam data directory.
	RootDirPath() string

	// UserDataDirPath returns the path to the user data directory.
	UserDataDirPath() (string, os.FileInfo, error)

	// UserIdsToDataDirPaths returns a map of local Steam user IDs
	// to their data storage directories.
	UserIdsToDataDirPaths() (map[string]string, error)

	// ShortcutsFilePath returns the path to the shortcuts file for a
	// given Steam user ID.
	ShortcutsFilePath(userId string) (string, os.FileInfo, error)

	// GridDirPath returns the path to the grid images directory
	// for a given Steam user ID.
	GridDirPath(userId string) (string, os.FileInfo, error)
}

type defaultDataVerifier struct {
	dataDir string
}

func (o defaultDataVerifier) ShortcutsFilePath(userId string) (string, os.FileInfo, error) {
	filePath := ShortcutsFilePath(o.dataDir, userId)

	i, err := os.Stat(filePath)
	if err != nil {
		return "", nil, err
	}

	return filePath, i, nil
}

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

func (o defaultDataVerifier) UserDataDirPath() (string, os.FileInfo, error) {
	dirPath := UserDataDirPath(o.dataDir)

	i, err := os.Stat(dirPath)
	if err != nil {
		return "", nil, err
	}

	return dirPath, i, nil
}

func (o defaultDataVerifier) RootDirPath() string {
	return o.dataDir
}

func (o defaultDataVerifier) GridDirPath(userId string) (string, os.FileInfo, error) {
	dirPath := GridDirPath(o.dataDir, userId)

	i, err := os.Stat(dirPath)
	if err != nil {
		return "", nil, err
	}

	return dirPath, i, nil
}

// GridDirPath generates a path to the grid images directory for the specified
// data directory and Steam user ID.
func GridDirPath(dataDirPath string, userId string) string {
	return path.Join(dataDirPath, userDataDirName, userId, userConfigDirName, gridDirName)
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

func homePath() (string, error) {
	homePath := os.Getenv("HOME")

	if len(strings.TrimSpace(homePath)) == 0 {
		return "", errors.New("The HOME environment variable is not set")
	}

	return homePath, nil
}
