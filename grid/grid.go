package grid

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/stephen-fox/steamutil/locations"
	"github.com/stephen-fox/steamutil/naming"
)

const (
	defaultImageMode = 0644
)

// ImageDetails stores important details about the grid image, such as
// the owner, and location.
type ImageDetails struct {
	// DataVerifier is used to get the grid images directory path.
	DataVerifier locations.DataVerifier

	// OwnerUserId is the Steam user ID for the image being operated on.
	OwnerUserId string

	// GameName is the name of the game for which the grid image is for.
	GameName string

	// GameExecutablePath is the full executable path (including any
	// quotation marks or other characters) for the grid image's game.
	GameExecutablePath string
}

// Validate returns a non-nil error if the ImageDetails is invalid.
func (o *ImageDetails) Validate() error {
	if o.DataVerifier == nil {
		return errors.New("the DataVerifier cannot be nil")
	}

	if len(strings.TrimSpace(o.OwnerUserId)) == 0 {
		return errors.New("please specify a Steam user ID")
	}

	return nil
}

// FilePath generates a file path for the grid image. It does not test if the
// image exists - however, it does check if the grid image directory for the
// given user does exist. An optional file extension can be provided, which
// will be appended to the end of the file path.
func (o *ImageDetails) FilePath(optionalExtension string) (string, error) {
	err := o.Validate()
	if err != nil {
		return "", err
	}

	gridDirPath, _, err := o.DataVerifier.GridDirPath(o.OwnerUserId)
	if err != nil {
		return "", err
	}

	gameId := naming.LegacyNonSteamGameId(o.GameName, o.GameExecutablePath)

	return path.Join(gridDirPath, gameId) + optionalExtension, nil
}

// AddConfig configures the grid image addition operation.
type AddConfig struct {
	// ResultDetails specifies details about the resulting image.
	ResultDetails ImageDetails

	// ImageSourcePath is the source path of the image being operated on.
	ImageSourcePath string

	// OverwriteExisting specifies whether or not an existing grid image
	// should be overwritten.
	OverwriteExisting bool

	// Mode specifies the os.FileMode for the resulting grid image file.
	// If not specified, defaultImageMode will be used.
	Mode os.FileMode
}

// Validate returns a non-nil error if the AddConfig is invalid.
func (o *AddConfig) Validate() error {
	err := o.ResultDetails.Validate()
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(o.ImageSourcePath)) == 0 {
		return errors.New("please specify a tile image source path")
	}

	if o.Mode == 0 {
		o.Mode = defaultImageMode
	}

	return nil
}

// RemoveConfig configures the grid image removal operation.
type RemoveConfig struct {
	// TargetDetails specifies details about the image to be removed.
	TargetDetails ImageDetails

	// FileExtension is the file extension to target. If not set,
	// the remove operation will target all files matching the
	// specified TargetDetails.
	FileExtension string
}

// Validate returns a non-nil error if the RemoveConfig is invalid.
func (o *RemoveConfig) Validate() error {
	err := o.TargetDetails.Validate()
	if err != nil {
		return err
	}

	return nil
}

// AddImage adds an image as a Steam grid image.
func AddImage(config AddConfig) error {
	err := config.Validate()
	if err != nil {
		return err
	}

	gridDirPath, _, err := config.ResultDetails.DataVerifier.GridDirPath(config.ResultDetails.OwnerUserId)
	if err != nil {
		return err
	}

	gameId := naming.LegacyNonSteamGameId(config.ResultDetails.GameName, config.ResultDetails.GameExecutablePath)

	resultingFilePath := path.Join(gridDirPath, gameId)

	extensionIndex := strings.LastIndex(config.ImageSourcePath, ".")
	if extensionIndex > 0 {
		resultingFilePath = resultingFilePath + config.ImageSourcePath[extensionIndex:]
	}

	if !config.OverwriteExisting {
		_, statErr := os.Stat(resultingFilePath)
		if statErr == nil {
			return nil
		}
	}

	source, err := os.Open(config.ImageSourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.OpenFile(resultingFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, config.Mode)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		return err
	}

	return nil
}

// RemoveImage removes a Steam grid image.
func RemoveImage(config RemoveConfig) error {
	err := config.Validate()
	if err != nil {
		return err
	}

	filePath, err := config.TargetDetails.FilePath(config.FileExtension)
	if err != nil {
		return err
	}

	if len(config.FileExtension) != 0 {
		os.Remove(filePath)

		return nil
	}

	dirPath := path.Dir(filePath)

	infos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	filenamePrefix := path.Base(filePath)

	for _, info := range infos {
		if info.IsDir() {
			continue
		}

		if strings.HasPrefix(info.Name(), filenamePrefix) {
			os.Remove(path.Join(dirPath, info.Name()))
		}
	}

	return nil
}
