package shortcuts

import (
	"testing"
	"io/ioutil"
)

func TestShortcut_VdfString(t *testing.T) {
	s := Shortcut{
		Id:      0,
		AppName: "My Cool app",
		ExePath:  "/path/to/my cool app",
		StartDir: "/path/to",
		IconPath: "",
		ShortcutPath: "",
		LaunchOptions: "-one -two \"-three and some\"",
		LastPlayTime: "",
		Tags: "",
	}

	ioutil.WriteFile("/Users/sfox/Desktop/test", []byte(s.VdfString()), 0600)
}
