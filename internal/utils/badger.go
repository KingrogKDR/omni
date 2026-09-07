package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/KingrogKDR/omni/internal/storage"
	"github.com/dgraph-io/badger"
)

type BadgerEngine struct {
	db *badger.DB
}

func cfKey(cf, key []byte) []byte {
	out := make([]byte, 0, len(cf)+1+len(key))
	out = append(out, cf...)
	out = append(out, '_')
	out = append(out, key...)
	return out
}

func NewBadgerEngine(path string) (*BadgerEngine, error) {
	bo := badger.DefaultOptions(path)
	db, err := badger.Open(bo.WithLogger(nil))
	return &BadgerEngine{
		db: db,
	}, err
}

func (b *BadgerEngine) WriteBatch(ctx context.Context, ops []storage.WriteOp) error {
	return b.db.Update(func(txn *badger.Txn) error {
		for _, operation := range ops {
			switch op := operation.(type) {
			case storage.Put:
				if err := txn.Set(cfKey(op.CF, op.Key), op.Val); err != nil {
					return fmt.Errorf("PUT %s: %w", op.Key, err)
				}
			case storage.Delete:
				if err := txn.Delete(cfKey(op.CF, op.Key)); err != nil {
					return fmt.Errorf("DELETE %s: %w", op.Key, err)
				}
			default:
				return fmt.Errorf("unknown write op: %T", operation)
			}
		}
		return nil
	})
}

func (b *BadgerEngine) NewReader(ctx context.Context) (storage.StorageReader, error) {
	txn := b.db.NewTransaction(false) // read-only
	return &BadgerReader{
		txn: txn,
	}, nil
}

type BadgerReader struct {
	txn *badger.Txn
}

func (r *BadgerReader) GetCF(cf, key []byte) ([]byte, error) {
	item, err := r.txn.Get(cfKey(cf, key))
	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil, storage.ErrKeyNotFound
	}
	if err != nil {
		return nil, err
	}

	return item.ValueCopy(nil)
}

func (r *BadgerReader) IterCF(ctx context.Context, cf []byte, opts storage.ScanOptions) ([]storage.Pair, error) {
	if len(opts.Prefix) > 0 && len(opts.Start) > 0 && !bytes.HasPrefix(opts.Start, opts.Prefix) {
		return nil, fmt.Errorf("SCAN: start must begin with prefix when both are set")
	}

	bo := badger.DefaultIteratorOptions
	bo.Prefix = cfKey(cf, nil)
	if len(opts.Prefix) > 0 {
		bo.Prefix = cfKey(cf, opts.Prefix)
	}
	bo.Reverse = opts.Reverse

	it := r.txn.NewIterator(bo)
	defer it.Close()

	if len(opts.Start) > 0 {
		it.Seek(cfKey(cf, opts.Start))
	} else {
		it.Rewind()
	}

	endKey := cfKey(cf, opts.End)
	limit := uint32(0)

	pairs := make([]storage.Pair, 0)
	for ; it.Valid() && (opts.Limit == 0 || limit < opts.Limit); it.Next() {
		item := it.Item()
		if len(opts.End) > 0 && ((opts.Reverse && bytes.Compare(item.Key(), endKey) <= 0) || bytes.Compare(item.Key(), endKey) >= 0) {
			break
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		pairs = append(pairs, storage.Pair{
			Key: item.KeyCopy(nil)[len(cf)+1:],
			Val: val,
		})
		limit++
	}

	return pairs, nil
}

func (r *BadgerReader) Close() {
	r.txn.Discard()
}

func (b *BadgerEngine) Close() error {
	return b.db.Close()
}
