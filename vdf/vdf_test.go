package vdf

import (
	"encoding/binary"
	"encoding/hex"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

const (
	testDataSubDir       = "/.testdata/"
	testDataOutputSubDir = testDataSubDir + "output/"
	shortcutsVdfSubDir   = testDataSubDir + "shortcuts-vdf-v1/"
	threeEntriesVdfName  = "3-entries.vdf"
)

func TestNewConstructorV1(t *testing.T) {
	o := Config{
		Name: "junk",
		Version: V1,
	}

	_, err := NewConstructor(o)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestNewConstructorUnknown(t *testing.T) {
	o := Config{
		Name: "junk",
		Version: FormatVersion(3333),
	}

	_, err := NewConstructor(o)
	if err == nil {
		t.Error("New constructor with bad version did not error")
	}
}

func TestV1Constructor_Read(t *testing.T) {
	o := Config{
		Name:    "shortcuts",
		Version: V1,
	}

	c, err := NewConstructor(o)
	if err != nil {
		t.Error(err.Error())
	}

	vdfv1Path, err := shortcutsVdfV1TestPath()
	if err != nil {
		t.Error(err.Error())
	}

	f, err := os.Open(vdfv1Path + threeEntriesVdfName)
	if err != nil {
		t.Error(err.Error())
	}
	defer f.Close()

	v, err := c.Read(f)
	if err != nil {
		t.Error(err.Error())
	}

	s, err := v.String()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Error(err.Error())
	}

	raw, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err.Error())
	}

	original := string(raw)

	if s != original {
		t.Error("New string does not match current file\nExp: '" + original + "'\nGot: '" + s + "'")
	}
}

func TestV1Constructor_Build(t *testing.T) {
	o := Config{
		Name:    "test",
		Version: V1,
	}

	c, err := NewConstructor(o)
	if err != nil {
		t.Error(err.Error())
	}

	o0 := NewObject([]Field{
		NewIdField(0),
		NewStringField("MyField", "One"),
		NewBoolField("IsAwesome", true),
		NewInt32Field("Blerg", 666),
		NewSliceField("Stuff", []string{"one thing", "two thing"}),
	})

	o1 := NewObject([]Field{
		NewIdField(1),
		NewStringField("MyField", "Woah"),
		NewBoolField("IsAwesome", false),
		NewInt32Field("Blerg", 0),
		NewSliceField("Stuff", []string{}),
	})

	v := c.Build([]Object{o0, o1})

	result, err := v.String()
	if err != nil {
		t.Error(err.Error())
	}

	sb := &strings.Builder{}
	sb.Write([]byte{0, 't', 'e', 's', 't', 0, 0, '0', 0, 1})
	sb.WriteString("MyField")
	sb.WriteString(null)
	sb.WriteString("One")
	sb.WriteString(null)
	sb.WriteString(intField)
	sb.WriteString("IsAwesome")
	sb.WriteString(null)
	sb.WriteString(soh)
	sb.WriteString(strings.Repeat(null, 3))
	sb.WriteString(intField)
	sb.WriteString("Blerg")
	sb.WriteString(null)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, 666)
	sb.Write(b)
	sb.WriteString(null)
	sb.WriteString("Stuff")
	sb.WriteString(null)
	sb.WriteString(stringField)
	sb.WriteString("0")
	sb.WriteString(null)
	sb.WriteString("one thing")
	sb.WriteString(null)
	sb.WriteString(stringField)
	sb.WriteString("1")
	sb.WriteString(null)
	sb.WriteString("two thing")

	sb.Write([]byte{0, 8, 8, 0, '1', 0, 1})
	sb.WriteString("MyField")
	sb.WriteString(null)
	sb.WriteString("Woah")
	sb.WriteString(null)
	sb.WriteString(intField)
	sb.WriteString("IsAwesome")
	sb.WriteString(null)
	sb.WriteString(strings.Repeat(null, 4))
	sb.WriteString(intField)
	sb.WriteString("Blerg")
	sb.WriteString(null)
	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, 0)
	sb.Write(b)
	sb.WriteString(null)
	sb.WriteString("Stuff")

	sb.Write([]byte{0, 8, 8, 8, 8})

	expect := sb.String()

	if result != expect {
		t.Error("Generater vdf string is not equal to expected value.\nExp: '" +
			expect + "'\nGot: '" + result + "'")
	}
}

func shortcutsVdfV1TestPath() (string, error) {
	p, err := repoPath()
	if err != nil {
		return "", err
	}

	return p + shortcutsVdfSubDir, nil
}

func repoPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Dir(wd), nil
}

func getHash(r io.Reader, hash hash.Hash) (string, error) {
	_, err := io.Copy(hash, r)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
