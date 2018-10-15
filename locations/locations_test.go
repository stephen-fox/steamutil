package locations

import (
	"log"
	"path"
	"testing"
)

func TestDataDirPath(t *testing.T) {
	p, i, err := DataDirPath()
	if err != nil {
		t.Error(err.Error())
	}

	if i == nil {
		t.Error("Info is nil")
	}

	log.Println(p)
}

func TestUserDataDirPath(t *testing.T) {
	if !IsInstalled() {
		t.Skip()
	}

	p, i, err := UserDataDirPath()
	if err != nil {
		t.Error(err.Error())
	}

	if i == nil {
		t.Error("Info is nil")
	}

	log.Println(p)
}

func TestUserIdsToDataDirPaths(t *testing.T) {
	if !IsInstalled() {
		t.Skip()
	}

	m, err := UserIdsToDataDirPaths()
	if err != nil {
		t.Error(err.Error())
	}

	log.Println(m)
}

func TestShortcutsFilePath(t *testing.T) {
	if !IsInstalled() {
		t.Skip()
	}

	p, i, err := ShortcutsFilePath("34161670")
	if err != nil {
		log.Println(err.Error())
		return
	}

	if i == nil {
		t.Error("Info is nil")
	}

	if path.Base(p) != shortcutsFileName {
		t.Error("File name should be", shortcutsFileName)
	}

	log.Println(p)
}
