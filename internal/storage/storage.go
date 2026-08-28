package storage

import "context"

type WriteType int
type ReadType int

const (
	PUT WriteType = iota
	DELETE
)

const (
	GET ReadType = iota
	LIST
	SCAN
)

type WriteRequest struct {
	Type  WriteType
	Key   []byte
	Value []byte
}

type ReadRequest struct {
	Type   ReadType
	Key    []byte
	Prefix []byte
	Limit  uint32
}

type ReadResponse struct {
	Value []byte
	Map   map[string]string
}

type Storage interface {
	Write(ctx context.Context, req WriteRequest) error
	Read(ctx context.Context, req ReadRequest) (ReadResponse, error)
}
