package storage

import "context"

type WriteTypes int
type ReadTypes int

const (
	PUT    WriteTypes = 0
	DELETE WriteTypes = 1
)

const (
	GET ReadTypes = 0
)

type Storage interface {
	Write(ctx context.Context, key []byte, value []byte) error
	Read(ctx context.Context, key []byte) ([]byte, error)
}
