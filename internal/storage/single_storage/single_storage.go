package SingleStorage

import (
	"context"

	"github.com/KingrogKDR/omni/internal/storage"
)

// The caller is responsible for closing the engine when it is no longer
// needed. Closing releases resources associated with the engine.
type Engine interface {
	Start() error
	Read(ctx context.Context, op storage.ReadOp) ([]byte, error)
	WriteBatch(ctx context.Context, ops []storage.WriteOp) error
	Close() error
}

type SingleStorage struct {
	engine Engine
}

func NewSingleStorage(engine Engine) *SingleStorage {
	return &SingleStorage{
		engine: engine,
	}
}

func (s *SingleStorage) Start() {}

func (s *SingleStorage) Reader(ctx context.Context, operation storage.ReadOp) ([]byte, error) {
	return s.engine.Read(ctx, operation)
}

func (s *SingleStorage) Writer(ctx context.Context, batch []storage.WriteOp) error {
	return s.engine.WriteBatch(ctx, batch)
}

func (s *SingleStorage) Stop() {
	_ = s.engine.Close()
}
