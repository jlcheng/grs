package test

import (
	"io/ioutil"
	"jcheng/grs/grsdb"
	"os"
	"testing"
	"fmt"
	"github.com/theckman/go-flock"
	"path/filepath"
)

func TestMemKVStore(t *testing.T) {
	var s = grsdb.NewMemKVStore()
	initStr := "init contents of foo"
	s.SaveBytes("foo", []byte(initStr))

	got, err := s.LoadBytes("foo")
	if err != nil {
		t.Fatal("TestMemKVStore")
	}
	if string(got) != initStr {
		t.Fatal("TestMemKVStore")
	}
}

func TestDiskKVStore(t *testing.T) {
	oldwd, err := os.Getwd()
	d, err := ioutil.TempDir("", "grstest")
	if err != nil {
		t.Fatalf("TestDiskKVStore: %v\n", err)
	}
	defer func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatal("TestDiskKVStore.defer: ", err)
		}
		if err := os.RemoveAll(d); err != nil {
			t.Fatal("TestDiskKVStore.defer: ", err)
		}
	}()
	s, err := grsdb.InitDiskKVStore(d)
	if err != nil {
		t.Fatal("TestDiskKVStore: ", err)
	}

	initStr := "init contents of foo"
	err = s.SaveBytes("foo", []byte(initStr))
	if err != nil {
		t.Fatal("TestDiskKVStore: ", err)
	}

	got, err := s.LoadBytes("foo")
	if err != nil {
		t.Fatal("TestDiskKVStore: ", err)
	}
	if string(got) != initStr {
		t.Fatal("TestDiskKVStore")
	}
}

func TestDiskKVStore_Lock(t *testing.T) {
	oldwd, err := os.Getwd()
	d, err := ioutil.TempDir("", "grstest")
	if err != nil {
		t.Fatalf("TestDiskKVStore: %v\n", err)
	}
	defer func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatal("TestDiskKVStore.defer: ", err)
		}
		if err := os.RemoveAll(d); err != nil {
			t.Fatal("TestDiskKVStore.defer: ", err)
		}
	}()
	s, err := grsdb.InitDiskKVStore(d)
	if err != nil {
		t.Fatal("TestDiskKVStore: ", err)
	}
	full := filepath.Join(d, "/foo")
	lockpath := fmt.Sprintf("%v.lock", full)
	f := flock.NewFlock(lockpath)
	fmt.Println(f)
	fmt.Println(f.Locked())
	_, err = f.TryLock()
	fmt.Println(f.Locked())
	if err != nil {
		t.Fatal(err)
	}
	err = s.SaveBytes("foo", []byte("bar"))
	if err != nil {
		t.Fatal(err)
	}
}
