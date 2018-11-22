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
	tenEntriesVdfName    = "10-entries.vdf"
)

func TestReadVdfV1File(t *testing.T) {
	rp, err := shortcutsVdfV1TestPath()
	if err != nil {
		t.Fatal(err.Error())
	}

	f, err := os.Open(rp + threeEntriesVdfName)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	shortcuts, err := ReadVdfV1File(f)
	if err != nil {
		t.Fatal(err.Error())
	}

	var gotIds []int

	for _, s := range shortcuts {
		gotIds = append(gotIds, s.Id)
		switch s.Id {
		case 0:
			if s.AppName != "Install macOS High Sierra" {
				t.Fatal("Unexpected app name for", s.Id, "- '" + s.AppName + "'")
			}
			if s.ExePath != "/Users/sfox/Desktop/Install macOS High Sierra.app" {
				t.Fatal("Unexpected exe path for", s.Id, "- '" + s.ExePath + "'")
			}
			if s.StartDir != "/Users/sfox/Desktop/" {
				t.Fatal("Unexpected start dir for", s.Id, "- '" + s.StartDir + "'")
			}
			if s.LaunchOptions != "" {
				t.Fatal("Unexpected launch options for", s.Id, "- '" +  s.LaunchOptions + "'")
			}
			if s.IconPath != "" {
				t.Fatal("Unexpected icon path for", s.Id, "- '" + s.IconPath + "'")
			}
			if !s.AllowOverlay {
				t.Fatal("Unexpected allow overlay for", s.Id, "-", s.AllowOverlay)
			}
			if !s.AllowDesktopConfig {
				t.Fatal("Unexpected allow desktop config for", s.Id, "-", s.AllowDesktopConfig)
			}
			if s.IsHidden {
				t.Fatal("Unexpected is hidden for", s.Id, "-", s.IsHidden)
			}
			if s.IsOpenVr {
				t.Fatal("Unexpected is open cr for", s.Id, "-", s.IsOpenVr)
			}
			if s.LastPlayTimeEpoch != 0 {
				t.Fatal("Unexpected last play time epoch for", s.Id, "-", s.LastPlayTimeEpoch)
			}
			if len(s.Tags) != 2 {
				t.Fatal("Unexpected tags length for", s.Id, "-", len(s.Tags))
			} else {
				if s.Tags[0] != "junk" {
					t.Fatal("Unexpected tag[0] for", s.Id, "- '" + s.Tags[0] + "'")
				}
				if s.Tags[1] != "eee" {
					t.Fatal("Unexpected tag[1] for", s.Id, "- '" + s.Tags[1] + "'")
				}
			}
		case 1:
			if s.AppName != "Automator" {
				t.Fatal("Unexpected app name for", s.Id, "- '" + s.AppName + "'")
			}
			if s.ExePath != "/Applications/Automator.app" {
				t.Fatal("Unexpected exe path for", s.Id, "- '" + s.ExePath + "'")
			}
			if s.StartDir != "/Applications/" {
				t.Fatal("Unexpected start dir for", s.Id, "- '" + s.StartDir + "'")
			}
			if s.LaunchOptions != "" {
				t.Fatal("Unexpected launch options for", s.Id, "- '" +  s.LaunchOptions + "'")
			}
			if s.IconPath != "" {
				t.Fatal("Unexpected icon path for", s.Id, "- '" + s.IconPath + "'")
			}
			if !s.AllowOverlay {
				t.Fatal("Unexpected allow overlay for", s.Id, "-", s.AllowOverlay)
			}
			if !s.AllowDesktopConfig {
				t.Fatal("Unexpected allow desktop config for", s.Id, "-", s.AllowDesktopConfig)
			}
			if s.IsHidden {
				t.Fatal("Unexpected is hidden for", s.Id, "-", s.IsHidden)
			}
			if s.IsOpenVr {
				t.Fatal("Unexpected is open cr for", s.Id, "-", s.IsOpenVr)
			}
			if s.LastPlayTimeEpoch != 1538319747 {
				t.Fatal("Unexpected last play time epoch for", s.Id, "-", s.LastPlayTimeEpoch)
			}
			if len(s.Tags) != 0 {
				t.Fatal("Unexpected tags length for", s.Id, "-", len(s.Tags))
			}
		case 2:
			if s.AppName != "Chess" {
				t.Fatal("Unexpected app name for", s.Id, "- '" + s.AppName + "'")
			}
			if s.ExePath != "/Applications/Chess.app" {
				t.Fatal("Unexpected exe path for", s.Id, "- '" + s.ExePath + "'")
			}
			if s.StartDir != "/Applications/" {
				t.Fatal("Unexpected start dir for", s.Id, "- '" + s.StartDir + "'")
			}
			if s.LaunchOptions != "abc" {
				t.Fatal("Unexpected launch options for", s.Id, "- '" +  s.LaunchOptions + "'")
			}
			if s.IconPath != "" {
				t.Fatal("Unexpected icon path for", s.Id, "- '" + s.IconPath + "'")
			}
			if !s.AllowOverlay {
				t.Fatal("Unexpected allow overlay for", s.Id, "-", s.AllowOverlay)
			}
			if !s.AllowDesktopConfig {
				t.Fatal("Unexpected allow desktop config for", s.Id, "-", s.AllowDesktopConfig)
			}
			if s.IsHidden {
				t.Fatal("Unexpected is hidden for", s.Id, "-", s.IsHidden)
			}
			if !s.IsOpenVr {
				t.Fatal("Unexpected is open vr for", s.Id, "-", s.IsOpenVr)
			}
			if s.LastPlayTimeEpoch != 1538335537 {
				t.Fatal("Unexpected last play time epoch for", s.Id, "-", s.LastPlayTimeEpoch)
			}
			if len(s.Tags) != 1 {
				t.Fatal("Unexpected tags length for", s.Id, "-", len(s.Tags))
			} else {
				if s.Tags[0] != "junk" {
					t.Fatal("Unexpected tag[0] for", s.Id, "- '" + s.Tags[0] + "'")
				}
			}
		default:
			t.Fatal("Unexpected shortcut in slice. ID is", s.Id)
		}
	}

	if len(gotIds) != 3 {
		t.Fatal("Did not get the epxected number of shortcut entries. Got -", len(gotIds))
	}

	if gotIds[0] != 0 {
		t.Fatal("Shortcut ID at 0 is wrong. Got -", gotIds[0])
	}

	if gotIds[1] != 1 {
		t.Fatal("Shortcut ID at 1 is wrong. Got -", gotIds[1])
	}

	if gotIds[2] != 2 {
		t.Fatal("Shortcut ID at 2 is wrong. Got -", gotIds[2])
	}
}

// TODO: Actually test values.
func TestReadVdfV1File10Entries(t *testing.T) {
	rp, err := shortcutsVdfV1TestPath()
	if err != nil {
		t.Fatal(err.Error())
	}

	f, err := os.Open(rp + tenEntriesVdfName)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	scs, err := ReadVdfV1File(f)
	if err != nil {
		t.Fatal(err.Error())
	}

	for i, s := range scs {
		if i != s.Id {
			t.Fatal("Got unexpected shortcut ID at index", i, "-", s.Id)
		}
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
