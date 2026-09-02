package storage

import (
	"bytes"
	"context"
	"fmt"
	"iter"

	"github.com/dgraph-io/badger"
)

type BadgerStore struct {
	db *badger.DB
}

func NewBadgerStore() (*BadgerStore, error) {
	opts := badger.DefaultOptions("./data")
	badger_store, err := badger.Open(opts)

	return &BadgerStore{
		db: badger_store,
	}, err
}

func (s *BadgerStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	var value []byte
	txErr := s.db.View(func(txn *badger.Txn) error {
		var err error

		item, err := txn.Get(key)
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
		return nil, fmt.Errorf("transaction (GET): %w", txErr)
	}
	return value, nil
}

func (s *BadgerStore) Put(ctx context.Context, key, value []byte) error {
	txErr := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		if err != nil {
			return fmt.Errorf("store put: %w", err)
		}
		return nil
	})
	if txErr != nil {
		return fmt.Errorf("transaction (PUT): %w", txErr)
	}
	return nil
}

func (s *BadgerStore) Delete(ctx context.Context, key []byte) error {
	txErr := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		if err != nil {
			return fmt.Errorf("store delete: %w", err)
		}
		return nil
	})
	if txErr != nil {
		return fmt.Errorf("transaction (DELETE): %w", txErr)
	}
	return nil
}

func (s *BadgerStore) List(ctx context.Context, prefix []byte, limit uint32) (iter.Seq2[[]byte, []byte], *ReadError) {
	readErr := &ReadError{}
	seq := func(yield func(k, v []byte) bool) {
		txErr := s.db.View(func(txn *badger.Txn) error {
			opt := badger.DefaultIteratorOptions
			opt.Prefix = prefix

			it := txn.NewIterator(opt)
			defer it.Close()

			var count uint32
			for it.Rewind(); it.Valid() && (limit == 0 || count < limit); it.Next() {
				if ctx.Err() != nil {
					return fmt.Errorf("list: %w", ctx.Err())
				}

				item := it.Item()
				key := item.KeyCopy(nil)
				val, err := item.ValueCopy(nil)
				if err != nil {
					return fmt.Errorf("list value copy: %w", err)
				}

				// the loop exited normally
				// return in place because we must not call yield again after it's returned false, otherwise it causes runtime panic
				if !yield(key, val) {
					return nil
				}

				count++
			}

			return nil
		})
		if txErr != nil {
			readErr.err = fmt.Errorf("transaction (LIST): %w", txErr)
		}
	}
	return seq, readErr
}

func (s *BadgerStore) Scan(ctx context.Context, start, end []byte) (iter.Seq2[[]byte, []byte], *ReadError) {
	readErr := &ReadError{}
	seq := func(yield func(k, v []byte) bool) {
		txErr := s.db.View(func(txn *badger.Txn) error {
			opt := badger.DefaultIteratorOptions
			it := txn.NewIterator(opt)
			defer it.Close()

			for it.Seek(start); it.Valid(); it.Next() {
				if ctx.Err() != nil {
					return fmt.Errorf("scan: %w", ctx.Err())
				}

				item := it.Item()
				key := item.KeyCopy(nil)

				if end != nil && bytes.Compare(key, end) >= 0 {
					return nil
				}

				val, err := item.ValueCopy(nil)
				if err != nil {
					return fmt.Errorf("scan value copy: %w", err)
				}

				if !yield(key, val) {
					return nil
				}
			}
			return nil
		})
		if txErr != nil {
			readErr.err = fmt.Errorf("transaction (SCAN): %w", txErr)
		}
	}

	return seq, readErr
}

func (s *BadgerStore) Close() error {
	return s.db.Close()
}
