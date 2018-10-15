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

func TestDefaultDataVerifier_UserDataDirPath(t *testing.T) {
	if !IsInstalled() {
		t.Skip()
	}

	v, err := NewDataVerifier()
	if err != nil {
		t.Error(err.Error())
	}

	p, i, err := v.UserDataDirPath()
	if err != nil {
		t.Error(err.Error())
	}

	if i == nil {
		t.Error("Info is nil")
	}

	log.Println(p)
}

func TestDefaultDataVerifier_UserIdsToDataDirPaths(t *testing.T) {
	if !IsInstalled() {
		t.Skip()
	}

	v, err := NewDataVerifier()
	if err != nil {
		t.Error(err.Error())
	}

	m, err := v.UserIdsToDataDirPaths()
	if err != nil {
		t.Error(err.Error())
	}

	log.Println(m)
}

func TestDefaultDataVerifier_ShortcutsFilePath(t *testing.T) {
	if !IsInstalled() {
		t.Skip()
	}

	v, err := NewDataVerifier()
	if err != nil {
		t.Error(err.Error())
	}

	p, i, err := v.ShortcutsFilePath("34161670")
	if err != nil {
		t.Error(err.Error())
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
