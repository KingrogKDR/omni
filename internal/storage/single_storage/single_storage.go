package singleStorage

import (
	"context"
	"log"

	"github.com/KingrogKDR/omni/internal/storage"
	"github.com/KingrogKDR/omni/internal/utils"
)

// The caller is responsible for closing the engine when it is no longer needed
type Engine interface {
	WriteBatch(ctx context.Context, ops []storage.WriteOp) error
	NewReader(ctx context.Context) (storage.StorageReader, error)
	Close() error
}

type SingleStorage struct {
	engine Engine
	path   string
}

func NewSingleStorage(path string) *SingleStorage {
	return &SingleStorage{
		path: path,
	}
}

func (s *SingleStorage) Start() error {
	engine, err := utils.NewBadgerEngine(s.path)
	if err != nil {
		return err
	}
	s.engine = engine
	return nil
}

func (s *SingleStorage) Reader(ctx context.Context) (storage.StorageReader, error) {
	return s.engine.NewReader(ctx)
}

func (s *SingleStorage) Writer(ctx context.Context, batch []storage.WriteOp) error {
	return s.engine.WriteBatch(ctx, batch)
}

func (s *SingleStorage) Stop() {
	if err := s.engine.Close(); err != nil {
		log.Println("engine close error:", err)
	}
}
