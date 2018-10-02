package shortcuts

import (
	"fmt"
	"os"
	"testing"
)

func TestNewShortcuts(t *testing.T) {
	f, err := os.Open("/Users/sfox/Library/Application Support/Steam/userdata/34161670/config/shortcuts.vdf")
	if err != nil {
		t.Error(err.Error())
	}
	defer f.Close()

	shortcuts, err := Shortcuts(f)
	if err != nil {
		t.Error(err.Error())
	}

	for _, s := range shortcuts {
		fmt.Println(s)
	}
}

func TestNewShortcut(t *testing.T) {
	// TODO: todo.
}
