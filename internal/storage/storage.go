package storage

import (
	"context"
	"iter"
)

type ReadError struct {
	err error
}

func (e *ReadError) Err() error { return e.err }

func (e *ReadError) SetErr(err error) { e.err = err }

type Writer interface {
	Put(context.Context, []byte, []byte) error
	Delete(context.Context, []byte) error
}

type Reader interface {
	Get(context.Context, []byte) ([]byte, error)
	List(context.Context, []byte, uint32) (iter.Seq2[[]byte, []byte], *ReadError)
	// Scan(context.Context, []byte, []byte)
}

type Storage interface {
	Reader
	Writer
}
