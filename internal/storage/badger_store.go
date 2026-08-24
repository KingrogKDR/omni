package storage

import (
	"context"
	"fmt"

	"github.com/dgraph-io/badger"
)

type BadgerStore struct {
	DB *badger.DB
}

func NewBadgerStore(opts badger.Options) (*BadgerStore, error) {
	badger_store, err := badger.Open(opts)

	return &BadgerStore{
		DB: badger_store,
	}, err
}

func (s *BadgerStore) Read(ctx context.Context, key []byte) ([]byte, error) {
	var value []byte
	txErr := s.DB.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return fmt.Errorf("store read: %w", err)
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("read item: %w", err)
		}
		value = val
		return nil
	})
	if txErr != nil {
		return nil, fmt.Errorf("transaction: %w", txErr)
	}
	return value, nil
}

func (s *BadgerStore) Write(ctx context.Context, key []byte, value []byte) error {
	txErr := s.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		if err != nil {
			return fmt.Errorf("store read: %w", err)
		}
		return nil
	})
	if txErr != nil {
		return fmt.Errorf("transaction: %w", txErr)
	}
	return nil
}
