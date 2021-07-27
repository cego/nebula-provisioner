package server

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/slackhq/nebula"
	"os"
)

type store struct {
	db *badger.DB
}

func (s store) Close() error {
	return s.db.Close()
}

func NewStore(config *nebula.Config) (*store, error) {

	dbPath := config.GetString("db.path", "/tmp/nebula-provisioner/db")
	stat, err := os.Stat(dbPath)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dbPath)
	}

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		return nil, err
	}


	return &store{db}, nil
}
