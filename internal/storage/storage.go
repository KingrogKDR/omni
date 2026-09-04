package storage

import (
	"context"
	"iter"
)

// StorageReader provides read-only access to a storage snapshot.
//
// The caller is responsible for closing the reader when it is no longer
// needed. Closing releases resources associated with the reader.
type StorageReader interface {
	GetCF(ctx context.Context, cf []byte, key []byte) ([]byte, error)

	// IterCF returns a sequence of key-value pairs from the specified
	// column family.
	//
	// The sequence is limited to at most limit entries. A limit of zero
	// means no limit.
	//
	// The caller is responsible for consuming the sequence and handling
	// any error returned by the iterator.
	IterCF(ctx context.Context, cf []byte, limit uint32) (iter.Seq2[[]byte, []byte], *error)
	Close()
}

// WriteOp represents a supported write operation.
//
// Only types that implement the unexported isWriteOp method can be used
// as WriteOp values. This prevents unrelated types from being passed to
// the storage writer.
type WriteOp interface {
	isWriteOp()
}

type PutOp struct {
	Key []byte
	Val []byte
	Cf  []byte
}

type DeleteOp struct {
	Key []byte
	Cf  []byte
}

func (PutOp) isWriteOp()    {}
func (DeleteOp) isWriteOp() {}

// Storage provides access to the storage engine.
//
// It separates reading from writing: Reader returns a reader for
// read operations, while Write applies a batch of write operations.
type Storage interface {
	Reader(ctx context.Context) (StorageReader, error)
	Write(ctx context.Context, batch []WriteOp) error
}
