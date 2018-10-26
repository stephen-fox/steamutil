package locations

import (
	"os"
	"path"
)

// DataDirPath returns the path to Steam's data directory.
func DataDirPath() (string, os.FileInfo, error) {
	homePath, err := homePath()
	if err != nil {
		return "", nil, err
	}

	dirPath := path.Join(homePath, ".steam", "root")

	i, err := os.Stat(dirPath)
	if err != nil {
		return "", nil, err
	}

	return dirPath, i, nil
}
