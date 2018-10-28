package shortcuts

import "testing"

func TestShortcut_EqualsNegative(t *testing.T) {
	s1 := Shortcut{
		Id:                1,
		AppName:           "woah",
		ExePath:           "/path/to/something",
		StartDir:          "/path/to",
		IconPath:          "/icon.png",
		ShortcutPath:      "",
		LaunchOptions:     "-one -two \"-three and some\"",
		LastPlayTimeEpoch: 1538448950,
		Tags: []string{
			"cool",
			"story",
		},
	}

	s2 := Shortcut{
		Id:                2,
		AppName:           "Another Cool app",
		ExePath:           "/what/to/my other cool app",
		StartDir:          "/what/to",
		IconPath:          "",
		ShortcutPath:      "/shortcut",
		LaunchOptions:     "-zero -two \"-three and some\"",
		LastPlayTimeEpoch: 1538448951,
		Tags: []string{
			"uncool",
			"story",
		},
	}

	if s1.Equals(s2) {
		t.Error("Shortcut 1 equals shortcut 2 even though they are different")
	}
}

func TestShortcut_EqualsPositive(t *testing.T) {
	s1 := Shortcut{
		Id:                1,
		AppName:           "woah",
		ExePath:           "/path/to/something",
		StartDir:          "/path/to",
		IconPath:          "/icon.png",
		ShortcutPath:      "",
		LaunchOptions:     "-one -two \"-three and some\"",
		LastPlayTimeEpoch: 1538448950,
		Tags: []string{
			"cool",
			"story",
		},
	}

	s2 := s1

	if !s1.Equals(s2) {
		t.Error("Shortcut 1 is not equal even though they are the same")
	}
}
