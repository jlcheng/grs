package grsdb

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

func BadgerDbAsString(dir string) error {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			value, _ := item.Value()
			fmt.Printf("%v=%v\n", string(item.Key()), string(value))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
