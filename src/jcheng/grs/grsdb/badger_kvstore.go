package grsdb

import (
	"github.com/dgraph-io/badger"
)

type BadgerKVStore struct {
	dir string
}
func NewBadgerKVStore(dir string) *BadgerKVStore {
	return &BadgerKVStore{
		dir: dir,
	}
}


func (s *BadgerKVStore) SaveBytes(key string, val []byte) error {
	opts := badger.DefaultOptions
	opts.Dir = s.dir
	opts.ValueDir = s.dir
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(txn *badger.Txn) error {
		txn.Set([]byte(key), val)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *BadgerKVStore) LoadBytes(key string) ([]byte, error) {
	opts := badger.DefaultOptions
	opts.Dir = s.dir
	opts.ValueDir = s.dir
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var iPtr *badger.Item
	err = db.View(func(txn *badger.Txn) error {
		item, rerr := txn.Get([]byte(key))
		if rerr != nil {
			return rerr
		}
		iPtr = item
		return nil
	})
	if err != nil {
		return nil, err
	}
	return iPtr.Value()
}