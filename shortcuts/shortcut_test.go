package shortcuts

import (
	"bytes"
	"crypto/sha1"
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

	dir, err := setupTestDataOutputDir()
	if err != nil {
		t.Error(err.Error())
	}

	err = ioutil.WriteFile(dir + "shortcut.vdf", []byte(s.VdfV1String([]byte{})), 0600)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestWriteVdfV1(t *testing.T) {
	shortcuts := []Shortcut{
		{
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
		},
		{
			Id:                1,
			AppName:           "woah",
			ExePath:           "/path/to/something",
			StartDir:          "/path/to",
			IconPath:          "",
			ShortcutPath:      "",
			LaunchOptions:     "-one -two \"-three and some\"",
			LastPlayTimeEpoch: 1538448950,
			Tags: []string{
				"cool",
				"story",
			},
		},
		{
			Id:                2,
			AppName:           "Another Cool app",
			ExePath:           "/path/to/my other cool app",
			StartDir:          "/path/to",
			IconPath:          "",
			ShortcutPath:      "",
			LaunchOptions:     "-one -two \"-three and some\"",
			LastPlayTimeEpoch: 1538448950,
			Tags: []string{
				"cool",
				"story",
			},
		},
	}

	dir, err := setupTestDataOutputDir()
	if err != nil {
		t.Error(err.Error())
	}

	f, err := os.OpenFile(dir + "/shortcuts.vdf", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		t.Error(err.Error())
	}

	err = WriteVdfV1(shortcuts, f)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestReadAndWrite(t *testing.T) {
	p, err := shortcutsVdfV1TestPath()
	if err != nil {
		t.Error(err.Error())
	}

	f, err := os.Open(p + threeEntriesVdfName)
	if err != nil {
		t.Error(err.Error())
	}
	defer f.Close()

	shortcuts, err := Shortcuts(f)
	if err != nil {
		t.Error(err.Error())
	}

	newFileBuffer := bytes.NewBuffer([]byte{})

	err = WriteVdfV1(shortcuts, newFileBuffer)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Error(err.Error())
	}

	originalHash, err := getHash(f, sha1.New())
	if err != nil {
		t.Error(err.Error())
	}

	newHash, err := getHash(newFileBuffer, sha1.New())
	if err != nil {
		t.Error(err.Error())
	}

	if newHash != originalHash {
		t.Error("Hashes do not match")
	}
}

func setupTestDataOutputDir() (string, error) {
	rp, err := repoPath()
	if err != nil {
		return "", err
	}

	fullPath := rp + testDataOutputSubDir

	err = os.MkdirAll(fullPath, 0700)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}
