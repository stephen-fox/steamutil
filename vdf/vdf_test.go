package vdf

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
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

func TestV1Constructor_Read(t *testing.T) {
	o := Options{
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

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Error(err.Error())
	}

	b, _ := ioutil.ReadAll(f)

	fmt.Println("Exp:", string(b))

	s, err := v.String()
	if err != nil {
		t.Error()
	}

	fmt.Println("Got:", s)

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Error(err.Error())
	}

	originalHash, err := getHash(f, sha1.New())
	if err != nil {
		t.Error(err.Error())
	}

	r := strings.NewReader(s)

	newHash, err := getHash(r, sha1.New())
	if err != nil {
		t.Error(err.Error())
	}

	if newHash != originalHash {
		t.Error("hashes not equal")
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