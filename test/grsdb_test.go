package test

import (
	"testing"
	"jcheng/grs/grsdb"
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
	db := &grsdb.DB{Repos:make([]grsdb.Repo,1)}
	db.Repos[0].Id = "/foo/bar"
	db.Repos[0].FetchedSec = 1
	grsdb.SaveFile(w.Write, "foo", db)
	if s := string(w.Bytes); s != "{\"repos\":[{\"id\":\"/foo/bar\",\"fetched_sec\":1}]}" {
		t.Errorf("unexpected saved data: %v\n", s)
	}
}
