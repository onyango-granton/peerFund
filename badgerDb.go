package badgerDb

import (
    "fmt"

    "github.com/dgraph-io/badger/v3"
)

var db *badger.DB

// OpenDB initializes and opens the BadgerDB
func OpenDB(path string) error {
    var err error
    db, err = badger.Open(badger.DefaultOptions(path))
    if err != nil {
        return err
    }
    return nil
}

// CloseDB closes the BadgerDB
func CloseDB() {
    db.Close()
}

// StoreData stores a key-value pair in BadgerDB
func StoreData(key string, value string) error {
    return db.Update(func(txn *badger.Txn) error {
        err := txn.Set([]byte(key), []byte(value))
        if err != nil {
            return err
        }
        return nil
    })
}
