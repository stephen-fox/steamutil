package shortcuts

import (
	"io"
	"os"

	"github.com/stephen-fox/steamutil/vdf"
)

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
