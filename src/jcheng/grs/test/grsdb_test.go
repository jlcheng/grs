package test

import (
	"jcheng/grs/grsdb"
	"testing"
)

func TestFindRepo_Ok(t *testing.T) {
	db := &grsdb.DB{Repos: make([]grsdb.RepoDTO, 3)}
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
	db := &grsdb.DB{Repos: make([]grsdb.RepoDTO, 3)}
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

func TestFindOrCreateRepo_Ok(t *testing.T) {
	db := &grsdb.DB{Repos: make([]grsdb.RepoDTO, 3)}
	db.Repos[0].Id = "1"
	db.Repos[0].FetchedSec = 1
	db.Repos[1].Id = "11"
	db.Repos[1].FetchedSec = 2
	db.Repos[2].Id = "111"
	db.Repos[2].FetchedSec = 3

	r := db.FindOrCreateRepo("4")
	if len(db.Repos) != 4 {
		t.Fatal("TestFindOrCreateRepo_Ok")
	}
	if &db.Repos[3] != r {
		t.Fatal("TestFindOrCreateRepo_Ok")
	}
}
