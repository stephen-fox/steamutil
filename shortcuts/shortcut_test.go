package shortcuts

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestShortcut_VdfV1String(t *testing.T) {
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

	data := handJamShortcut(s)
	result := s.VdfV1String([]byte{})

	if result != data {
		t.Error("Shortcut string is not equal to data:\nExp: '" + data + "'\nGot: '" + result + "'")
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

func handJamShortcut(s Shortcut) string {
	sb := &strings.Builder{}

	sb.WriteString("0")
	sb.WriteString(null)
	sb.WriteString(stringField)
	sb.WriteString(appNameField)
	sb.WriteString(null)
	sb.WriteString(s.AppName)
	sb.WriteString(null)

	sb.WriteString(stringField)
	sb.WriteString(exePathField)
	sb.WriteString(null)
	sb.WriteString("\"")
	sb.WriteString(s.ExePath)
	sb.WriteString("\"")
	sb.WriteString(null)

	sb.WriteString(stringField)
	sb.WriteString(startDirField)
	sb.WriteString(null)
	sb.WriteString("\"")
	sb.WriteString(s.StartDir)
	sb.WriteString("\"")
	sb.WriteString(null)

	sb.WriteString(stringField)
	sb.WriteString(iconPathField)
	sb.WriteString(null)
	if len(s.IconPath) > 0 {
		sb.WriteString(s.IconPath)
	}
	sb.WriteString(null)

	sb.WriteString(stringField)
	sb.WriteString(shortcutPathField)
	sb.WriteString(null)
	if len(s.ShortcutPath) > 0 {
		sb.WriteString(s.ShortcutPath)
	}
	sb.WriteString(null)

	sb.WriteString(stringField)
	sb.WriteString(launchOptionsField)
	sb.WriteString(null)
	if len(s.LaunchOptions) > 0 {
		sb.WriteString(s.LaunchOptions)
	}
	sb.WriteString(null)

	sb.WriteString(intField)
	sb.WriteString(isHiddenField)
	sb.WriteString(null)
	if s.IsHidden {
		sb.WriteString(soh)
	} else {
		sb.WriteString(null)
	}
	sb.WriteString(strings.Repeat(null, 3))

	sb.WriteString(intField)
	sb.WriteString(allowDesktopConfigField)
	sb.WriteString(null)
	if s.AllowDesktopConfig {
		sb.WriteString(soh)
	} else {
		sb.WriteString(null)
	}
	sb.WriteString(strings.Repeat(null, 3))

	sb.WriteString(intField)
	sb.WriteString(allowOverlayField)
	sb.WriteString(null)
	if s.AllowOverlay {
		sb.WriteString(soh)
	} else {
		sb.WriteString(null)
	}
	sb.WriteString(strings.Repeat(null, 3))

	sb.WriteString(intField)
	sb.WriteString(isOpenVrField)
	sb.WriteString(null)
	if s.IsHidden {
		sb.WriteString(soh)
	} else {
		sb.WriteString(null)
	}
	sb.WriteString(strings.Repeat(null, 3))

	sb.WriteString(intField)
	sb.WriteString(lastPlayTimeField)
	sb.WriteString(null)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(s.LastPlayTimeEpoch))
	sb.Write(b)

	sb.WriteString(null)
	sb.WriteString(tagsField)
	for i, v := range s.Tags {
		sb.WriteString(null)
		sb.WriteString(soh)
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(null)
		sb.WriteString(v)
	}

	sb.WriteString(null)

	return sb.String()
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
