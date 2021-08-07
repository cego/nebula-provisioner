package store

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/vault/shamir"
	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula"
	"io/ioutil"
	"os"
	"path"
)

type Store struct {
	l        *logrus.Logger
	config   *nebula.Config
	unsealed chan interface{}
	path     string
	db       *badger.DB
	keyParts [][]byte
}

const initFileName = "INITIALIZED"

func (s *Store) Initialize(numParts, threshold uint32) ([]string, error) {
	if s.IsInitialized() {
		return nil, fmt.Errorf("server is already initialized")
	}
	if s.IsOpen() {
		return nil, fmt.Errorf("server is already unsealed")
	}

	// Generating encryption key
	ek := make([]byte, 32)
	_, err := rand.Read(ek)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random encryption key")
	}

	secretParts, err := shamir.Split(ek, int(numParts), int(threshold))
	if err != nil {
		return nil, fmt.Errorf("failed to spilt encryption key using shamir: %v", err)
	}

	var keyParts []string
	for _, bytePart := range secretParts {
		keyParts = append(keyParts, hex.EncodeToString(bytePart))
	}

	err = s.open(ek)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize store: %v", err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			s.l.WithError(err).Println("Failed to close store after initializing")
		}
	}()

	return keyParts, nil
}

func (s *Store) open(encryptionKey []byte) error {
	if !s.IsInitialized() {

		err := ioutil.WriteFile(path.Join(s.path, initFileName), []byte(""), 0600)
		if err != nil {
			return fmt.Errorf("Failed to create file %s with error: %s\n", path.Join(s.path, initFileName), err)
		}
	}
	opts := badger.DefaultOptions(s.path)
	opts.EncryptionKey = encryptionKey
	opts.IndexCacheSize = 10 << 20 // 10MB
	opts.Logger = s.l

	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	s.db = db

	return nil
}

func (s *Store) IsInitialized() bool {
	if _, err := os.Stat(path.Join(s.path, initFileName)); err == nil {
		return true
	}
	return false
}

func (s *Store) IsOpen() bool {
	return s.db != nil && !s.db.IsClosed()
}

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Store) Unseal(keyPart string, removeExistingParts bool) error {
	if !s.IsInitialized() {
		return fmt.Errorf("server not initialized")
	}
	if s.IsOpen() {
		return fmt.Errorf("server is already unsealed")
	}

	decodedPart, err := hex.DecodeString(keyPart)
	if err != nil {
		return fmt.Errorf("failed to decode key part")
	}

	s.keyParts = appendIfMissing(s.keyParts, decodedPart)

	ek, err := shamir.Combine(s.keyParts)
	if err != nil {
		return fmt.Errorf("failed to combine encryption key using shamir: %s", err)
	}
	s.keyParts = nil

	if err := s.open(ek); err != nil {
		return fmt.Errorf("failed to open store after unsealed encryption key: %s", err)
	}

	s.unsealed <- true

	return nil
}

func NewStore(l *logrus.Logger, config *nebula.Config, unsealed chan interface{}) (*Store, error) {

	dbPath := config.GetString("db.path", "/tmp/nebula-provisioner/db")
	stat, err := os.Stat(dbPath)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dbPath)
	}

	return &Store{l, config, unsealed, dbPath, nil, make([][]byte, 0)}, nil
}

func appendIfMissing(slice [][]byte, b []byte) [][]byte {
	for _, ele := range slice {
		if bytes.Compare(ele, b) == 0 {
			return slice
		}
	}
	return append(slice, b)
}

func exists(txn *badger.Txn, prefix, key []byte) bool {
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	it := txn.NewKeyIterator(append(prefix, key...), opt)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		return true
	}
	return false
}
