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

func (s *BadgerStore) Read(ctx context.Context, req ReadRequest) ([]byte, error) {
	var value []byte
	txErr := s.DB.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(req.Key)
		if err != nil {
			return fmt.Errorf("store get: %w", err)
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("read item: %w", err)
		}
		value = val
		return nil
	})
	if txErr != nil {
		return nil, fmt.Errorf("transaction read: %w", txErr)
	}
	return value, nil
}

func (s *BadgerStore) Write(ctx context.Context, req WriteRequest) error {
	txErr := s.DB.Update(func(txn *badger.Txn) error {
		switch req.Type {
		case PUT:
			err := txn.Set(req.Key, req.Value)
			if err != nil {
				return fmt.Errorf("store put: %w", err)
			}
		case DELETE:
			err := txn.Delete(req.Key)
			if err != nil {
				return fmt.Errorf("store delete: %w", err)
			}
		}
		return nil
	})
	if txErr != nil {
		return fmt.Errorf("transaction write: %w", txErr)
	}
	return nil
}

func (s *BadgerStore) List(ctx context.Context) (ListResponse, error) {
	var resp ListResponse
	txErr := s.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		resp.KeyValue = make(map[string]string, 1)
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				resp.KeyValue[string(k)] = string(v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if txErr != nil {
		return ListResponse{}, fmt.Errorf("transaction list: %w", txErr)
	}
	return resp, nil
}
