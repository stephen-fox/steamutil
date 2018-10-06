package shortcuts

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path"
	"testing"
)

const (
	testDataSubDir       = "/.testdata/"
	testDataOutputSubDir = testDataSubDir + "output/"
	shortcutsVdfSubDir   = testDataSubDir + "shortcuts-vdf-v1/"
	threeEntriesVdfName  = "3-entries.vdf"
)

func TestNewShortcuts(t *testing.T) {
	rp, err := shortcutsVdfV1TestPath()
	if err != nil {
		t.Error(err.Error())
	}

	f, err := os.Open(rp + threeEntriesVdfName)
	if err != nil {
		t.Error(err.Error())
	}
	defer f.Close()

	shortcuts, err := Shortcuts(f)
	if err != nil {
		t.Error(err.Error())
	}

	// TODO: Actually confirm the deserialized data is equal to the file's data.
	for _, s := range shortcuts {
		fmt.Println(s)
	}
}

func TestNewShortcut(t *testing.T) {
	// TODO: todo.
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
