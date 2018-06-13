package grsio

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewCopyDir(t *testing.T) {
	oldwd, tmpdir := MkTmpDir(t, "grstest", "TestNewCopyDir")
	defer CleanTmpDir(t, oldwd, tmpdir, "TestNewCopyDir")

	err := os.MkdirAll(filepath.Join(tmpdir, "lvl1", "lvl2", "lvl3"), 0755)
	if err != nil {
		t.Fatal("TestNewCopyDir")
	}

	Touch(t, filepath.Join(tmpdir, "lvl1", "lvl1_1"), "TestNewCopyDir")
	Touch(t, filepath.Join(tmpdir, "lvl1", "lvl2", "lvl2_1"), "TestNewCopyDir")
	Touch(t, filepath.Join(tmpdir, "lvl1", "lvl2", "lvl2_2"), "TestNewCopyDir")

	tmpdir2 := tmpdir + "_2"
	CopyDir(tmpdir, tmpdir2)
	defer os.RemoveAll(tmpdir2)
	AssertExists(t, tmpdir2, true, "TestNewCopyDir")
	AssertExists(t, filepath.Join(tmpdir2, "lvl1"), true, "TestNewCopyDir")
	AssertExists(t, filepath.Join(tmpdir2, "lvl1", "lvl1_1"), false, "TestNewCopyDir")
	AssertExists(t, filepath.Join(tmpdir2, "lvl1", "lvl2"), true, "TestNewCopyDir")
	AssertExists(t, filepath.Join(tmpdir2, "lvl1", "lvl2", "lvl2_1"), false, "TestNewCopyDir")
	AssertExists(t, filepath.Join(tmpdir2, "lvl1", "lvl2", "lvl2_2"), false, "TestNewCopyDir")
	AssertExists(t, filepath.Join(tmpdir2, "lvl1", "lvl2", "lvl3"), true, "TestNewCopyDir")
}

func AssertExists(t *testing.T, path string, isDir bool, errId string) {
	stat, err := os.Stat(path)
	if err != nil {
		t.Fatal(errId)
	}
	if stat.IsDir() != isDir {
		t.Fatal(errId)
	}
}

func Touch(t *testing.T, path string, errId string) {
	f, err := os.Create(path)
	if f != nil {
		f.Close()
	}
	if err != nil {
		t.Fatal(errId)
	}
}

func MkTmpDir(t *testing.T, prefix string, errid string) (oldwd string, d string) {
	var err error
	oldwd, err = os.Getwd()
	if err != nil {
		t.Fatal(errid, err)
	}
	d, err = ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatal(errid, err)
	}
	return oldwd, d
}

func CleanTmpDir(t *testing.T, oldwd string, tmpdir string, errid string) {

	if err := os.Chdir(oldwd); err != nil {
		t.Fatal(errid, err)
	}
	if err := os.RemoveAll(tmpdir); err != nil {
		t.Fatal(errid, err)
	}
}
