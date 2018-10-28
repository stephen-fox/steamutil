package shortcuts

import (
	"bytes"
	"crypto/sha1"
	"os"
	"testing"
)

func TestWriteVdfV1(t *testing.T) {
	s := Shortcut{
		Id:                0,
		AppName:           "Chess",
		ExePath:           "/Applications/Chess.app",
		StartDir:          "/Applications",
		IconPath:          "",
		ShortcutPath:      "",
		LaunchOptions:     "-one -two \"-three and some\"",
		LastPlayTimeEpoch: 1538448950,
		Tags: []string{
			"cool",
			"story",
		},
	}

	b := bytes.NewBuffer(nil)

	err := WriteVdfV1([]Shortcut{s}, b)
	if err != nil {
		t.Fatal(err.Error())
	}

	scs, err := Shortcuts(b)
	if err != nil {
		t.Fatal(err.Error())
	}

	l := len(scs)
	if l != 1 {
		t.Fatal("Unexpected number of shortcuts -", l)
	}

	if !scs[0].Equals(s) {
		t.Fatal("Shortcuts are not equal")
	}
}

func TestWriteVdfV1MultipleEntries(t *testing.T) {
	scs := []Shortcut{
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

	buffer := bytes.NewBuffer(nil)
	err := WriteVdfV1(scs, buffer)
	if err != nil {
		t.Fatal(err.Error())
	}

	scsReadBack, err := Shortcuts(buffer)
	if err != nil {
		t.Fatal(err.Error())
	}

	l := len(scsReadBack)
	if l != len(scs) {
		t.Fatal("Read back unexpected number of shortcuts -", l)
	}

	for i := range scs {
		if !scs[i].Equals(scsReadBack[i]) {
			t.Fatal("Shortcut ", i, "is not equal")
		}
	}
}

func TestReadAndWriteFromFile(t *testing.T) {
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

	newFileBuffer := bytes.NewBuffer(nil)

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
