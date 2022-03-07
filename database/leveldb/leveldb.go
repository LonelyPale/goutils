package leveldb

import (
	"path"

	"github.com/syndtr/goleveldb/leveldb"
)

//todo

type GoLevelDB struct {
	db *leveldb.DB
}

func NewGoLevelDB(name string, dir string) (*GoLevelDB, error) {
	dbPath := path.Join(dir, name+".db")
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	database := &GoLevelDB{db: db}

	return database, nil
}
