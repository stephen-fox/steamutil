package shortcuts

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestShortcut_VdfString(t *testing.T) {
	s := Shortcut{
		Id:                0,
		AppName:           "My Cool app",
		ExePath:           "/path/to/my cool app",
		StartDir:          "/path/to",
		IconPath:          "",
		ShortcutPath:      "",
		LaunchOptions:     "-one -two \"-three and some\"",
		LastPlayTimeEpoch: 1538448950,
		Tags: []string{
			"cool",
			"story",
		},
	}

	rp, err := repoPath()
	if err != nil {
		t.Error(err.Error())
	}

	fullPath := rp + testDataOutputSubDir

	err = os.MkdirAll(fullPath, 0700)
	if err != nil {
		t.Error(err.Error())
	}

	err = ioutil.WriteFile(fullPath + "shortcut.vdf", []byte(s.VdfString()), 0600)
	if err != nil {
		t.Error(err.Error())
	}
}
