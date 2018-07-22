package grsdb

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"strings"
)

type BdbViewerOptions struct {
	Dir string
	TextMode bool
	Add string
}

var DefaultViewerOptions = BdbViewerOptions{
	Dir:"",
	TextMode: false,
}

func BadgerDbAsString(vopt BdbViewerOptions) error {
	opts := badger.DefaultOptions
	opts.Dir = vopt.Dir
	opts.ValueDir = vopt.Dir
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(txn *badger.Txn) error {
		fmt.Println(vopt)
		if vopt.Add != "" {
			tokens := strings.SplitN(vopt.Add, ":", 2)
			txn.Set([]byte(tokens[0]), []byte(tokens[1]))
		}

		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			value, _ := item.Value()
			keyStr := string(item.Key())
			valStr := "<unknown format>"
			if vopt.TextMode {
				valStr = string(value)
			}
			fmt.Printf("%v=%v\n", keyStr, valStr)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
