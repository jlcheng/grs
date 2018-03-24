package test

import (
	"bytes"
	"jcheng/grs/grsdb"
	"jcheng/grs/status"
	"reflect"
	"testing"
)

func TestLoadFile_Ok(t *testing.T) {
	db, err := grsdb.LoadFile("data/db.json")
	if err != nil {
		t.Error("cannot load data/db.json", err)
		return
	}
	if db.Repos[0].Id != "/foo/bar" {
		t.Errorf("Unexpected id: %v\n", db.Repos[0].Id)
		return
	}
}

func TestLoadFile_Fail(t *testing.T) {
	_, err := grsdb.LoadFile("")
	if err == nil {
		t.Error("Expected error, got none")
	}
}

func TestSaveFile_Ok(t *testing.T) {
	w := &grsdb.BufferedDBWriter{}
	db := &grsdb.DB{Repos: make([]grsdb.Repo, 2)}
	db.Repos[0].Id = "/foo/bar"
	db.Repos[0].FetchedSec = 1
	db.Repos[0].RStat = grsdb.RStat_Json{
		Dir:    status.DIR_VALID,
		Branch: status.BRANCH_DIVERGED,
		Index:  status.INDEX_MODIFIED,
	}
	grsdb.SaveFile(w.Write, "foo", db)
	reader := bytes.NewReader(w.Bytes)
	dbout, err := grsdb.LoadBytes(reader)
	if err != nil {
		t.Fatal("TestSaveFile_Ok:", err)
	}
	if !reflect.DeepEqual(*db, *dbout) {
		t.Fatal("Save/Load did not yield the same data")
	}
}

func TestFindRepo_Ok(t *testing.T) {
	db := &grsdb.DB{Repos: make([]grsdb.Repo, 3)}
	db.Repos[0].Id = "foo"
	db.Repos[0].FetchedSec = 1
	db.Repos[1].Id = "bar"
	db.Repos[1].FetchedSec = 2
	db.Repos[2].Id = "fizz"
	db.Repos[2].FetchedSec = 3

	if db.FindRepo("foo").FetchedSec != 1 ||
		db.FindRepo("bar").FetchedSec != 2 ||
		db.FindRepo("fizz").FetchedSec != 3 {
		t.Fatal("TestFindRepo_Ok")
	}
}

func TestFindRepo_Fail(t *testing.T) {
	db := &grsdb.DB{Repos: make([]grsdb.Repo, 3)}
	db.Repos[0].Id = "1"
	db.Repos[0].FetchedSec = 1
	db.Repos[1].Id = "11"
	db.Repos[1].FetchedSec = 2
	db.Repos[2].Id = "111"
	db.Repos[2].FetchedSec = 3

	if db.FindRepo("1111") != nil {
		t.Fatal("TestFindRepo_Fail")
	}
}
