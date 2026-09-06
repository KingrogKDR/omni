package storage

import (
	"context"
)

// ReadOp provides read-only access to a storage snapshot.
type ReadOp interface {
	isReadOp()
}

type Get struct {
	Cf  []byte
	Key []byte
}

type PrefixScan struct {
	Cf     []byte
	Prefix []byte
}

type RangeScan struct {
	Cf    []byte
	Start []byte
	End   []byte
}

func (Get) isReadOp()        {}
func (PrefixScan) isReadOp() {}
func (RangeScan) isReadOp()  {}

type IteratorOptions struct {
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
	Start()
	Reader(ctx context.Context, op ReadOp, opts IteratorOptions) ([]Pair, error)
	Writer(ctx context.Context, batch []WriteOp) error
	Stop()
}
