package grsdb

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/theckman/go-flock"
)

// KVStore interface is implemented by objects that can persist key/value pairs
type KVStore interface {
	SaveBytes(key string, val []byte) error
	LoadBytes(key string) ([]byte, error)
}

type MemKVStore struct {
	bins map[string][]byte
}

func (s *MemKVStore) SaveBytes(key string, val []byte) error {
	tmp := make([]byte, len(val))
	copy(tmp, val)
	s.bins[key] = tmp
	return nil
}

func (s *MemKVStore) LoadBytes(key string) ([]byte, error) {
	val, ok := s.bins[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("entry not found for [%v]", key))
	}
	return val, nil
}

func NewMemKVStore() KVStore {
	return &MemKVStore{
		bins: make(map[string][]byte),
	}
}

type DiskKVStore struct {
	path string
}

func InitDiskKVStore(path string) (KVStore, error) {
	path, err := filepath.Abs(path) // must be a directory
	if err != nil {
		return nil, err
	}
	return &DiskKVStore{
		path: filepath.FromSlash(path),
	}, nil
}

func (s *DiskKVStore) SaveBytes(key string, val []byte) error {
	full := filepath.Join(s.path, filepath.FromSlash(filepath.Clean(key)))

	// advisory lock
	lpath := fmt.Sprintf("%v.lock", full)
	f := flock.NewFlock(lpath)
	_, err := f.TryLock()
	if err != nil {
		fmt.Printf("cannot obtained lock on %v: %v\n", lpath, err)
		return err
	}
	// TODO: Consider handling failure to unlock. What can be done other than logging?
	defer func() {
		f.Unlock()
		fmt.Printf("released lock on %v\n", lpath)
	}()
	fmt.Printf("obtained lock on %v\n", lpath)
	return ioutil.WriteFile(full, val, 0644)
}

func (s *DiskKVStore) LoadBytes(key string) ([]byte, error) {
	full := filepath.Join(s.path, filepath.Clean(key))
	f, err := os.Open(full)
	if err != nil && !os.IsExist(err) {
		return nil, errors.New(fmt.Sprintf("key does not exist: %v", err))
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
