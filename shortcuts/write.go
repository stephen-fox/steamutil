package shortcuts

import (
	"io"

	"github.com/stephen-fox/steamutil/vdf"
)

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
