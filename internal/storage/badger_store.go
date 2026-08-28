package storage

import (
	"context"
	"fmt"

	"github.com/dgraph-io/badger"
)

type BadgerStore struct {
	DB *badger.DB
}

func NewBadgerStore() (*BadgerStore, error) {
	opts := badger.DefaultOptions("./data")
	badger_store, err := badger.Open(opts)

	return &BadgerStore{
		DB: badger_store,
	}, err
}

func (s *BadgerStore) Read(ctx context.Context, req ReadRequest) (ReadResponse, error) {
	var resp ReadResponse
	txErr := s.DB.Update(func(txn *badger.Txn) error {
		var err error
		switch req.Type {
		case GET:
			item, err := txn.Get(req.Key)
			if err != nil {
				return fmt.Errorf("store get: %w", err)
			}
			val, err := item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("read item: %w", err)
			}
			resp.Value = val
			return nil
		case LIST:
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			resp.Map = make(map[string]string, 1)
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				key := item.Key()
				err := item.Value(func(val []byte) error {
					resp.Map[string(key)] = string(val)
					return nil
				})
				if err != nil {
					return err
				}
			}
		case SCAN:
			opt := badger.DefaultIteratorOptions
			opt.Prefix = req.Prefix
			it := txn.NewIterator(opt)
			defer it.Close()

			resp.Map = make(map[string]string, 1)
			limit := uint32(0)
			for it.Rewind(); it.Valid(); it.Next() {
				if limit >= req.Limit {
					break
				}
				item := it.Item()
				key := item.Key()
				err := item.Value(func(val []byte) error {
					resp.Map[string(key)] = string(val)
					return nil
				})
				if err != nil {
					return err
				}
				limit += 1
			}
		default:
			err = fmt.Errorf("Not a valid read method\n")
		}
		return err
	})
	if txErr != nil {
		return ReadResponse{}, fmt.Errorf("transaction read: %w", txErr)
	}
	return resp, nil
}

func (s *BadgerStore) Write(ctx context.Context, req WriteRequest) error {
	txErr := s.DB.Update(func(txn *badger.Txn) error {
		var err error
		switch req.Type {
		case PUT:
			err = txn.Set(req.Key, req.Value)
			if err != nil {
				return fmt.Errorf("store put: %w", err)
			}
		case DELETE:
			err = txn.Delete(req.Key)
			if err != nil {
				return fmt.Errorf("store delete: %w", err)
			}
		default:
			err = fmt.Errorf("Not a valid write method\n")
		}
		return err
	})
	if txErr != nil {
		return fmt.Errorf("transaction write: %w", txErr)
	}
	return nil
}
