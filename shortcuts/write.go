package shortcuts

import (
	"errors"
	"io"
	"os"

	"github.com/stephen-fox/steamutil/vdf"
)

const (
	Unchanged      UpdateResult = "No changes were made to the file"
	CreatedNewFile UpdateResult = "Created new file"
	UpdatedEntry   UpdateResult = "Updated existing entry in the file"
	AddedNewEntry  UpdateResult = "Added new entry to the file"
)

type UpdateResult string

type CreateOrUpdateConfig struct {
	Path      string
	Mode      os.FileMode
	MatchName string
	OnMatch   func(name string, match *Shortcut)
	NoMatch   func(name string) Shortcut
}

func (o CreateOrUpdateConfig) IsValid() error {
	if len(o.MatchName) == 0 {
		return errors.New("The shortcut name to match cannot be empty")
	}

	if len(o.Path) == 0 {
		return errors.New("The shortcut file path cannot be empty")
	}

	if o.OnMatch == nil {
		return errors.New("The shortcut match function cannot be nil")
	}

	if o.NoMatch == nil {
		return errors.New("The shortcut not matched function cannot be nil")
	}

	return nil
}

func CreateOrUpdateVdfV1File(config CreateOrUpdateConfig) (UpdateResult, error) {
	err := config.IsValid()
	if err != nil {
		return Unchanged, err
	}

	alreadyExists := false

	_, statErr := os.Stat(config.Path)
	if statErr == nil {
		alreadyExists = true
	}

	f, err := os.OpenFile(config.Path, os.O_RDWR|os.O_CREATE, config.Mode)
	if err != nil {
		return Unchanged, err
	}
	defer f.Close()

	var currentScs []Shortcut

	if alreadyExists {
		currentScs, err = ReadVdfV1File(f)
		if err != nil {
			return Unchanged, err
		}
	}

	result := Unchanged

	for i := range currentScs {
		if currentScs[i].AppName == config.MatchName {
			result = UpdatedEntry

			config.OnMatch(config.MatchName, &currentScs[i])

			break
		}
	}

	if result == Unchanged {
		s := config.NoMatch(config.MatchName)
		s.Id = len(currentScs)
		currentScs = append(currentScs, s)
		result = AddedNewEntry
	}

	err = OverwriteVdfV1File(f, currentScs)
	if err != nil {
		return Unchanged, err
	}

	if alreadyExists {
		return result, nil
	}

	return CreatedNewFile, nil
}

func OverwriteVdfV1File(f *os.File, scs []Shortcut) error {
	_, err := f.Seek(0, 0)
	if err != nil {
		return err
	}

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	err = WriteVdfV1(scs, f)
	if err != nil {
		return err
	}

	return nil
}

func WriteVdfV1(shortcuts []Shortcut, w io.Writer) error {
	config := vdf.Config{
		Name:    header,
		Version: vdf.V1,
	}

	c, err := vdf.NewConstructor(config)
	if err != nil {
		return err
	}

	var objects []vdf.Object

	for _, s := range shortcuts {
		objects = append(objects, s.object())
	}

	v := c.Build(objects)

	s, err := v.String()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(s))
	if err != nil {
		return err
	}

	return nil
}
