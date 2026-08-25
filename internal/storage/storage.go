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
	Limit  int
}

type ListResponse struct {
	KeyValue map[string]string
}

type Storage interface {
	Write(ctx context.Context, req WriteRequest) error
	Read(ctx context.Context, req ReadRequest) ([]byte, error)
	List(ctx context.Context) (ListResponse, error)
}
