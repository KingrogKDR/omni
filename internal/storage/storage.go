package storage

import (
	"context"
	"errors"
)

var ErrKeyNotFound = errors.New("key not found")

// StorageReader provides read-only access to a storage snapshot.
type StorageReader interface {
	GetCF(cf, key []byte) ([]byte, error)
	IterCF(ctx context.Context, cf []byte, opts ScanOptions) ([]Pair, error)
	Close()
}

type ScanOptions struct {
	Prefix  []byte
	Start   []byte
	End     []byte
	Limit   uint32
	Reverse bool
}

type Pair struct {
	Key []byte
	Val []byte
}

// WriteOp represents a supported write operation.
//
// Only types that implement the unexported isWriteOp method can be used
// as WriteOp values. This prevents unrelated types from being passed to
// the storage writer.
type WriteOp interface {
	isWriteOp()
}

type Put struct {
	CF  []byte
	Key []byte
	Val []byte
}

type Delete struct {
	CF  []byte
	Key []byte
}

func (Put) isWriteOp()    {}
func (Delete) isWriteOp() {}

// Storage provides access to the storage engine.
//
// It separates reading from writing: Reader returns a reader for
// read operations, while Writer applies a batch of write operations.
type Storage interface {
	Start() error
	Reader(ctx context.Context) (StorageReader, error)
	Writer(ctx context.Context, batch []WriteOp) error
	Stop()
}
