package locations

import (
	"io/ioutil"
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

	p, i, err := v.ShortcutsFilePath(getSomeUserId(v))
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

func TestDefaultDataVerifier_GridDirPath(t *testing.T) {
	if !IsInstalled() {
		t.Skip()
	}

	v, err := NewDataVerifier()
	if err != nil {
		t.Fatal(err.Error())
	}

	_, _, err = v.GridDirPath(getSomeUserId(v))
	if err != nil {
		t.Fatal(err.Error())
	}
}

func getSomeUserId(dv DataVerifier) string {
	if !IsInstalled() {
		panic("Steam must be installed to get a user ID")
	}

	userDirPath, _, err := dv.UserDataDirPath()
	if err != nil {
		panic(err.Error())
	}

	infos, err := ioutil.ReadDir(userDirPath)
	if err != nil {
		panic(err.Error())
	}

	for _, info := range infos {
		if info.IsDir() {
			return info.Name()
		}
	}

	panic("No user IDs were found")

	return ""
}
