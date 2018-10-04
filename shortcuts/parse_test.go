package shortcuts

import (
	"fmt"
	"os"
	"path"
	"testing"
)

const (
	testDataSubDir       = "/.testdata/"
	testDataOutputSubDir = testDataSubDir + "output/"
	shortcutsVdfSubDir   = testDataSubDir + "shortcuts-vdf-v1/"
)

func TestNewShortcuts(t *testing.T) {
	rp, err := repoPath()
	if err != nil {
		t.Error(err.Error())
	}

	f, err := os.Open(rp + shortcutsVdfSubDir + "3-entries.vdf")
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

func repoPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Dir(wd), nil
}
