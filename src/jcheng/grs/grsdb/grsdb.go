package grsdb

import (
	"encoding/json"
	"strings"
)

type DB struct {
	Repos []RepoDTO `json:"repos"`
}

type RepoDTO struct {
	Id         string     `json:"id"`
	FetchedSec int64      `json:"fetched_sec"`
	RStat      RStat_Json `json:"rstat,omitempty"`
	MergedCnt  int        `json:"merged_cnt"`
	MergedSec  int64      `json:"merged_sec"`
}

type DBService interface {
	SaveDB(key string, root *DB) error
	LoadDB(key string) (*DB, error)
}

type DBServiceImpl struct {
	kvstore KVStore
}

func NewDBService(kvstore KVStore) DBService {
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

func (db *DB) FindRepo(id string) *RepoDTO {
	for idx, r := range db.Repos {
		if strings.Compare(id, r.Id) == 0 {
			return &db.Repos[idx]
		}
	}
	return nil
}

func (db *DB) FindOrCreateRepo(id string) *RepoDTO {
	for idx, r := range db.Repos {
		if strings.Compare(id, r.Id) == 0 {
			return &db.Repos[idx]
		}
	}
	db.Repos = append(db.Repos, RepoDTO{})
	return &db.Repos[len(db.Repos)-1]
}
