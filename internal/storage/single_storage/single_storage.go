package SingleStorage

import (
	"context"

	"github.com/KingrogKDR/omni/internal/storage"
)

type Engine interface {
	NewReader(ctx context.Context) (storage.StorageReader, error)
	WriteBatch(ctx context.Context, ops []storage.WriteOp) error
	Close() error
}

type SingleStorage struct {
	engine Engine
}

func NewSingleStorage() *SingleStorage {
	return &SingleStorage{}
}

func (s *SingleStorage) Start() {}

func (s *SingleStorage) Reader(ctx context.Context) (storage.StorageReader, error) {
	return s.engine.NewReader(ctx)
}

func (s *SingleStorage) Write(ctx context.Context, batch []storage.WriteOp) error {
	return s.engine.WriteBatch(ctx, batch)
}

func (s *SingleStorage) Stop() {}
