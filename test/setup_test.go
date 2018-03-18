package test

import (
	"io/ioutil"
	"jcheng/grs/config"
	"os"
	"path/filepath"
	"testing"
)

func TestSetupUserPrefDir(t *testing.T) {
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal("TestSetupUserPrefDir: ", err)
	}
	d, err := ioutil.TempDir("", "grstest")
	if err != nil {
		t.Fatal("TestSetupUserPrefDir: ", err)
	}
	defer func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatal("TestSetupUserPrefDir.defer: ", err)
		}
		if err := os.RemoveAll(d); err != nil {
			t.Fatal("TestSetupUserPrefDir.defer: ", err)
		}
	}()
	d2 := filepath.Join(d, ".grs.d")
	config.SetupUserPrefDir(d2)
	f, err := os.Stat(d2)
	if err != nil {
		t.Fatal("TestSetupUserPrefDir: ", err)
	}
	if !f.IsDir() {
		t.Fatal("TestSetupUserPrefDir: ", err)
	}
}
