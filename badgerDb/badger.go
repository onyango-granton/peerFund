package badgerDb

import (
    "github.com/dgraph-io/badger/v3"
    "log"
)

var db *badger.DB

func OpenDB(path string) error {
    var err error
    db, err = badger.Open(badger.DefaultOptions(path))
    if err != nil {
        return err
    }
    return nil
}

func CloseDB() {
    if err := db.Close(); err != nil {
        log.Fatal(err)
    }
}

func StoreData(key, value string) error {
    return db.Update(func(txn *badger.Txn) error {
        return txn.Set([]byte(key), []byte(value))
    })
}

func GetData(key string) (string, error) {
    var value string
    err := db.View(func(txn *badger.Txn) error {
        item, err := txn.Get([]byte(key))
        if err != nil {
            return err
        }
        val, err := item.ValueCopy(nil)
        if err != nil {
            return err
        }
        value = string(val)
        return nil
    })
    return value, err
}
