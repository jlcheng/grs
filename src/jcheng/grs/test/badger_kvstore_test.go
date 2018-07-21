package test

import (
	"testing"
	"jcheng/grs/grsdb"
	"os"
	"io/ioutil"
)

func TestBadgerKVStore(t *testing.T) {
	oldwd, err := os.Getwd()
	d, err := ioutil.TempDir("", "grstest")
	if err != nil {
		t.Fatal("TestBadgerKVStore: ", err)
	}
	defer func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatal("TestBadgerKVStore.defer: ", err)
		}
		if err := os.RemoveAll(d); err != nil {
			t.Fatal("TestBadgerKVStore.defer: ", err)
		}
	}()

	kvs := grsdb.NewBadgerKVStore(d)
	err = kvs.SaveBytes("key", []byte("badger"))
	if err != nil {
		t.Fatal("TestBadgerKVStore", err)
	}
	bytes, err := kvs.LoadBytes("key")
	if err != nil {
		t.Fatal(err)
	}
	if got := string(bytes); got != "badger" {
		t.Fatalf("expected [badger], got [%v]\n", got)
	}
}
