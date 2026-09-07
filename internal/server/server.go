package server

import (
	"context"
	"errors"

	"github.com/KingrogKDR/omni/internal/storage"
	kvpb "github.com/KingrogKDR/omni/proto/gen/kv"
)

type Server struct {
	kvpb.UnimplementedOmniServer

	store storage.Storage
}

func NewServer(store storage.Storage) *Server {
	return &Server{
		store: store,
	}
}

func (s *Server) Get(ctx context.Context, req *kvpb.GetRequest) (*kvpb.GetResponse, error) {
	reader, err := s.store.Reader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	val, err := reader.GetCF(req.Cf, req.Key)
	if errors.Is(err, storage.ErrKeyNotFound) {
		return &kvpb.GetResponse{NotFound: true}, nil
	}
	if err != nil {
		return nil, err
	}

	return &kvpb.GetResponse{
		Val: val,
	}, nil
}

func (s *Server) Put(ctx context.Context, req *kvpb.PutRequest) (*kvpb.PutResponse, error) {
	batch := make([]storage.WriteOp, 1)
	batch[0] = storage.Put{
		CF:  req.Cf,
		Key: req.Key,
		Val: req.Val,
	}
	err := s.store.Writer(ctx, batch)
	if err != nil {
		return nil, err
	}
	return &kvpb.PutResponse{Success: true}, nil
}

func (s *Server) Delete(ctx context.Context, req *kvpb.DeleteRequest) (*kvpb.DeleteResponse, error) {
	batch := make([]storage.WriteOp, 1)
	batch[0] = storage.Delete{
		CF:  req.Cf,
		Key: req.Key,
	}
	err := s.store.Writer(ctx, batch)
	if err != nil {
		return nil, err
	}
	return &kvpb.DeleteResponse{Success: true}, nil
}

func (s *Server) Scan(ctx context.Context, req *kvpb.ScanRequest) (*kvpb.ScanResponse, error) {
	so := req.GetScanOptions() // nil-safe even if req.ScanOptions is nil

	opts := storage.ScanOptions{
		Prefix:  so.GetPrefix(),
		Start:   so.GetStart(),
		End:     so.GetEnd(),
		Limit:   so.GetLimit(),
		Reverse: so.GetReverse(),
	}

	reader, err := s.store.Reader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	pairs, err := reader.IterCF(ctx, req.Cf, opts)

	if err != nil {
		return nil, err
	}

	kvPairs := make([]*kvpb.KeyValue, 0, len(pairs))
	for _, p := range pairs {
		kvPairs = append(kvPairs, &kvpb.KeyValue{
			Key:   p.Key,
			Value: p.Val,
		})
	}

	return &kvpb.ScanResponse{
		Pairs: kvPairs,
	}, nil

}
