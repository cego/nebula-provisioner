package store

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/vault/shamir"
	"github.com/sirupsen/logrus"
)

type Store struct {
	l        *logrus.Logger
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
	if len(encryptionKey) != 0 {
		opts.EncryptionKey = encryptionKey
	}
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

func NewStore(l *logrus.Logger, dataDir string, unsealed chan interface{}, encryptionEnabled bool) (*Store, error) {
	dbPath := filepath.Join(dataDir, "db")
	stat, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(dbPath, 0700)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if !stat.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", dbPath)
		}
	}

	s := &Store{l, unsealed, dbPath, nil, make([][]byte, 0)}

	if !encryptionEnabled {
		err = s.open([]byte{})
		if err != nil {
			return nil, err
		}
		unsealed <- true
	}

	return s, nil
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
	k := append(prefix, key...)
	it := txn.NewIterator(opt)
	defer it.Close()

	for it.Seek(k); it.ValidForPrefix(k); it.Next() {
		if bytes.Equal(it.Item().Key(), k) {
			return true
		}
	}
	return false
}

func containsIgnoreCase(s []string, e string) bool {
	e = strings.ToLower(e)
	for _, a := range s {
		if strings.ToLower(a) == e {
			return true
		}
	}
	return false
}

func containsByteSlice(array [][]byte, value []byte) bool {
	for _, v := range array {
		if bytes.Equal(v, value) {
			return true
		}
	}
	return false
}
