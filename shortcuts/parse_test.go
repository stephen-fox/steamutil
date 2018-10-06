package shortcuts

import (
	"encoding/hex"
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

	var gotIds []int

	for _, s := range shortcuts {
		gotIds = append(gotIds, s.Id)
		switch s.Id {
		case 0:
			if s.AppName != "Install macOS High Sierra" {
				t.Error("Unexpected app name for", s.Id, "- '" + s.AppName + "'")
			}
			if s.ExePath != "/Users/sfox/Desktop/Install macOS High Sierra.app" {
				t.Error("Unexpected exe path for", s.Id, "- '" + s.ExePath + "'")
			}
			if s.StartDir != "/Users/sfox/Desktop/" {
				t.Error("Unexpected start dir for", s.Id, "- '" + s.StartDir + "'")
			}
			if s.LaunchOptions != "" {
				t.Error("Unexpected launch options for", s.Id, "- '" +  s.LaunchOptions + "'")
			}
			if s.IconPath != "" {
				t.Error("Unexpected icon path for", s.Id, "- '" + s.IconPath + "'")
			}
			if !s.AllowOverlay {
				t.Error("Unexpected allow overlay for", s.Id, "-", s.AllowOverlay)
			}
			if !s.AllowDesktopConfig {
				t.Error("Unexpected allow desktop config for", s.Id, "-", s.AllowDesktopConfig)
			}
			if s.IsHidden {
				t.Error("Unexpected is hidden for", s.Id, "-", s.IsHidden)
			}
			if s.IsOpenVr {
				t.Error("Unexpected is open cr for", s.Id, "-", s.IsOpenVr)
			}
			if s.LastPlayTimeEpoch != 0 {
				t.Error("Unexpected last play time epoch for", s.Id, "-", s.LastPlayTimeEpoch)
			}
			if len(s.Tags) != 2 {
				t.Error("Unexpected tags length for", s.Id, "-", len(s.Tags))
			}
			if s.Tags[0] != "junk" {
				t.Error("Unexpected tag[0] for", s.Id, "- '" + s.Tags[0] + "'")
			}
			if s.Tags[1] != "eee" {
				t.Error("Unexpected tag[1] for", s.Id, "- '" + s.Tags[1] + "'")
			}
		case 1:
			if s.AppName != "Automator" {
				t.Error("Unexpected app name for", s.Id, "- '" + s.AppName + "'")
			}
			if s.ExePath != "/Applications/Automator.app" {
				t.Error("Unexpected exe path for", s.Id, "- '" + s.ExePath + "'")
			}
			if s.StartDir != "/Applications/" {
				t.Error("Unexpected start dir for", s.Id, "- '" + s.StartDir + "'")
			}
			if s.LaunchOptions != "" {
				t.Error("Unexpected launch options for", s.Id, "- '" +  s.LaunchOptions + "'")
			}
			if s.IconPath != "" {
				t.Error("Unexpected icon path for", s.Id, "- '" + s.IconPath + "'")
			}
			if !s.AllowOverlay {
				t.Error("Unexpected allow overlay for", s.Id, "-", s.AllowOverlay)
			}
			if !s.AllowDesktopConfig {
				t.Error("Unexpected allow desktop config for", s.Id, "-", s.AllowDesktopConfig)
			}
			if s.IsHidden {
				t.Error("Unexpected is hidden for", s.Id, "-", s.IsHidden)
			}
			if s.IsOpenVr {
				t.Error("Unexpected is open cr for", s.Id, "-", s.IsOpenVr)
			}
			if s.LastPlayTimeEpoch != 1538319747 {
				t.Error("Unexpected last play time epoch for", s.Id, "-", s.LastPlayTimeEpoch)
			}
			if len(s.Tags) != 0 {
				t.Error("Unexpected tags length for", s.Id, "-", len(s.Tags))
			}
		case 2:
			if s.AppName != "Chess" {
				t.Error("Unexpected app name for", s.Id, "- '" + s.AppName + "'")
			}
			if s.ExePath != "/Applications/Chess.app" {
				t.Error("Unexpected exe path for", s.Id, "- '" + s.ExePath + "'")
			}
			if s.StartDir != "/Applications/" {
				t.Error("Unexpected start dir for", s.Id, "- '" + s.StartDir + "'")
			}
			if s.LaunchOptions != "abc" {
				t.Error("Unexpected launch options for", s.Id, "- '" +  s.LaunchOptions + "'")
			}
			if s.IconPath != "" {
				t.Error("Unexpected icon path for", s.Id, "- '" + s.IconPath + "'")
			}
			if !s.AllowOverlay {
				t.Error("Unexpected allow overlay for", s.Id, "-", s.AllowOverlay)
			}
			if !s.AllowDesktopConfig {
				t.Error("Unexpected allow desktop config for", s.Id, "-", s.AllowDesktopConfig)
			}
			if s.IsHidden {
				t.Error("Unexpected is hidden for", s.Id, "-", s.IsHidden)
			}
			if !s.IsOpenVr {
				t.Error("Unexpected is open cr for", s.Id, "-", s.IsOpenVr)
			}
			if s.LastPlayTimeEpoch != 1538335537 {
				t.Error("Unexpected last play time epoch for", s.Id, "-", s.LastPlayTimeEpoch)
			}
			if len(s.Tags) != 1 {
				t.Error("Unexpected tags length for", s.Id, "-", len(s.Tags))
			}
			if s.Tags[0] != "junk" {
				t.Error("Unexpected tag[0] for", s.Id, "- '" + s.Tags[0] + "'")
			}
		default:
			t.Error("Unexpected shortcut in slice. ID is", s.Id)
		}
	}

	if len(gotIds) != 3 {
		t.Error("Did not get the epxected number of shortcut entries. Got -", len(gotIds))
	}

	if gotIds[0] != 0 {
		t.Error("Shortcut ID at 0 is wrong. Got -", gotIds[0])
	}

	if gotIds[1] != 1 {
		t.Error("Shortcut ID at 1 is wrong. Got -", gotIds[1])
	}

	if gotIds[2] != 2 {
		t.Error("Shortcut ID at 2 is wrong. Got -", gotIds[2])
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
