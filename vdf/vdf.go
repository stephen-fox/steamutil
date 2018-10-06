package vdf

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

const (
	// null is the null ASCII character.
	null = "\x00"
	// soh is the 'Start of Header' ASCII character.
	soh = "\x01"
	// stx is the 'Start of Text' ASCII character.
	stx = "\x02"

	stringField = soh
	intField    = stx
	sliceField  = null
)

const (
	V1 FormatVersion = 1
)

type FormatVersion int

type Vdf interface {
	Append(object Object)
	Objects() []Object
	String() (string, error)
}

type defaultVdf struct {
	objects []Object
	options Options
}

func (o *defaultVdf) Append(object Object) {
	o.objects = append(o.objects, object)
}

func (o *defaultVdf) Objects() []Object {
	return o.objects
}

func (o *defaultVdf) String() (string, error) {
	sb := &strings.Builder{}

	switch o.options.Version {
	case V1:
		break
	default:
		return "", errors.New("The specified format is not supported - " + string(o.options.Version))
	}

	sb.Write(o.options.header)

	for i, object := range o.objects {
		for _, field := range object.Fields() {
			err := field.Append(sb, o.options.Version)
			if err != nil {
				return "", err
			}
		}

		if i < len(o.objects) - 1 {
			sb.Write(o.options.objectDelim)
		}
	}

	sb.Write(o.options.footer)

	return sb.String(), nil
}

type Options struct {
	Name        string
	Version     FormatVersion
	header      []byte
	footer      []byte
	objectDelim []byte
}

type Constructor interface {
	Build(objects []Object) Vdf
	Read(r io.Reader) (Vdf, error)
}

type v1Constructor struct {
	options Options
}

func (o *v1Constructor) Build(objects []Object) Vdf {
	return &defaultVdf{
		objects: objects,
		options: o.options,
	}
}

func (o *v1Constructor) Read(r io.Reader) (Vdf, error) {
	result := &defaultVdf{
		options: o.options,
	}

	s := bufio.NewScanner(r)
	s.Split(o.split)

	for s.Scan() {
		object, err := ParseRawObject(s.Text(), result.options.Version)
		if err != nil {
			return result, err
		}

		result.Append(object)
	}

	err := s.Err()
	if err != nil {
		return result, err
	}

	return result, nil
}

func (o *v1Constructor) split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.Index(data, o.options.header); i >= 0 {
		return i + len(o.options.header), nil, nil
	}

	if i := bytes.Index(data, o.options.objectDelim); i >= 0 {
		return i + len(o.options.objectDelim), data[0:i], nil
	}

	if i := bytes.Index(data, o.options.footer); i >= 0 {
		return i + len(o.options.footer), data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func NewConstructor(options Options) (Constructor, error) {
	if len(options.Name) == 0 {
		return &v1Constructor{}, errors.New("Please specify a name for the vdf")
	}

	switch options.Version {
	case V1:
		options.header = fileHeaderV1(options.Name)
		options.objectDelim = []byte{8, 8, 0}
		options.footer = []byte{8, 8, 8, 8}
		return &v1Constructor{
			options: options,
		}, nil
	}

	return &v1Constructor{}, errors.New("The specified format is not supported - " + string(options.Version))
}

func fileHeaderV1(name string) []byte {
	header := []byte{0}
	header = append(header, []byte(name)...)
	header = append(header, 0, 0)

	return header
}
