package locations

import "path"

// DataDirPath returns the path to Steam's data directory.
func DataDirPath() (string, os.FileInfo, error) {
	dirPath := path.Join(os.Getenv("HOME"), "Library/Application Support/Steam")

	i, err := os.Stat(dirPath)
	if err != nil {
		return dirPath, nil, err
	}

	return dirPath, i, nil
}
