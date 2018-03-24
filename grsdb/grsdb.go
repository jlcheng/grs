package grsdb

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type DB struct {
	Repos []Repo `json:"repos"`
}

type Repo struct {
	Id         string     `json:"id"`
	FetchedSec int64      `json:"fetched_sec"`
	RStat      RStat_Json `json:"rstat,omitempty"`
}

func LoadBytes(reader io.Reader) (*DB, error) {
	bytes, err := ioutil.ReadAll(reader)
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

func LoadFile(filename string) (*DB, error) {
	reader, err := os.Open(filepath.FromSlash(filename))
	if err != nil {
		return nil, err
	}
	return LoadBytes(reader)
}

func SaveFile(writer DBWriter, filename string, db *DB) error {
	bytes, err := json.Marshal(db)
	if err != nil {
		return err
	}
	return writer(filename, bytes)
}

// DBWriter allows one to mock the API for writing to disk
type DBWriter func(filename string, bytes []byte) error

func FileDBWriter(filename string, bytes []byte) error {
	return ioutil.WriteFile(filepath.FromSlash(filename), bytes, 0644)
}

type BufferedDBWriter struct {
	Bytes []byte
}

func (w *BufferedDBWriter) Write(filename string, bytes []byte) error {
	w.Bytes = make([]byte, len(bytes))
	for i, b := range bytes {
		w.Bytes[i] = b
	}
	return nil
}

func (db DB) FindRepo(id string) *Repo {
	for idx, r := range db.Repos {
		if strings.Compare(id, r.Id) == 0 {
			return &db.Repos[idx]
		}
	}
	return nil
}
