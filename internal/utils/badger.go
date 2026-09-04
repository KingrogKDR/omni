package utils

import (
	"context"
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

func (b *BadgerEngine) WriteBatch(ctx context.Context, ops []storage.WriteOp) error {
	return b.db.Update(func(txn *badger.Txn) error {
		for _, operation := range ops {
			switch op := operation.(type) {
			case storage.PutOp:
				if err := txn.Set(cfKey(op.CF, op.Key), op.Val); err != nil {
					return fmt.Errorf("put %s: %w", op.Key, err)
				}
			case storage.DeleteOp:
				if err := txn.Delete(cfKey(op.CF, op.Key)); err != nil {
					return fmt.Errorf("delete %s: %w", op.Key, err)
				}
			default:
				return fmt.Errorf("unknown write op: %T", operation)
			}
		}
		return nil
	})
}

func (b *BadgerEngine) Close() error {
	return b.db.Close()
}
