package grsdb

import (
	"encoding/json"
	"strings"
)

type DB struct {
	Repos []Repo `json:"repos"`
}

type Repo struct {
	Id         string     `json:"id"`
	FetchedSec int64      `json:"fetched_sec"`
	RStat      RStat_Json `json:"rstat,omitempty"`
	MergedCnt  int        `json:"merged_cnt"`
}

type DBService interface {
	SaveDB(key string, root *DB) error
	LoadDB(key string) (*DB, error)
}

type DBServiceImpl struct {
	kvstore KVStore
}

func NewDBService(kvstore KVStore) (DBService) {
	return &DBServiceImpl{
		kvstore: kvstore,
	}
}

func (s *DBServiceImpl) SaveDB(key string, db *DB) error {
	bytes, err := json.Marshal(db)
	if err != nil {
		return err
	}
	return s.kvstore.SaveBytes(key, bytes)
}

func (s *DBServiceImpl) LoadDB(key string) (*DB, error) {
	bytes, err := s.kvstore.LoadBytes(key)
	if err != nil {
		return nil, err
	}
	db := &DB{}
	err = json.Unmarshal(bytes, db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) FindRepo(id string) *Repo {
	for idx, r := range db.Repos {
		if strings.Compare(id, r.Id) == 0 {
			return &db.Repos[idx]
		}
	}
	return nil
}


func (db *DB) FindOrCreateRepo(id string) *Repo {
	for idx, r := range db.Repos {
		if strings.Compare(id, r.Id) == 0 {
			return &db.Repos[idx]
		}
	}
	r := Repo{
		Id: id,
		FetchedSec: 0,
		RStat: RStat_Json{},
		MergedCnt: 0,
	}
	db.Repos = append(db.Repos, r)
	return &db.Repos[len(db.Repos)-1]
}