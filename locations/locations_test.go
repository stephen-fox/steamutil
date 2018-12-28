package locations

import (
	"log"
	"os"
	"path"
	"testing"
)

func TestDataDirPath(t *testing.T) {
	p, i, err := DataDirPath()
	if err != nil {
		t.Fatal(err.Error())
	} else if i == nil {
		t.Fatal("Info is nil")
	}

	log.Println(p)
}

func TestDefaultDataVerifier_DataDirPath(t *testing.T) {
	v, err := NewDataVerifier()
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = os.Stat(v.RootDirPath())
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDefaultDataVerifier_UserDataDirPath(t *testing.T) {
	if !IsInstalled() {
		t.Skip()
	}

	v, err := NewDataVerifier()
	if err != nil {
		t.Fatal(err.Error())
	}

	p, i, err := v.UserDataDirPath()
	if err != nil {
		t.Fatal(err.Error())
	} else if i == nil {
		t.Fatal("Info is nil")
	}

	log.Println(p)
}

func TestDefaultDataVerifier_UserIdsToDataDirPaths(t *testing.T) {
	if !IsInstalled() {
		t.Skip()
	}

	v, err := NewDataVerifier()
	if err != nil {
		t.Fatal(err.Error())
	}

	m, err := v.UserIdsToDataDirPaths()
	if err != nil {
		t.Fatal(err.Error())
	}

	log.Println(m)
}

func TestDefaultDataVerifier_ShortcutsFilePath(t *testing.T) {
	if !IsInstalled() {
		t.Skip()
	}

	v, err := NewDataVerifier()
	if err != nil {
		t.Fatal(err.Error())
	}

	p, i, err := v.ShortcutsFilePath("34161670")
	if err != nil {
		t.Fatal(err.Error())
		return
	} else if i == nil {
		t.Fatal("Info is nil")
	}

	if path.Base(p) != shortcutsFileName {
		t.Fatal("File name should be", shortcutsFileName)
	}

	log.Println(p)
}
