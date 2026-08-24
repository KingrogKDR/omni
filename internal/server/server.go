package server

import (
	"context"

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
	value, err := s.store.Read(ctx, req.Key)
	if err != nil {
		return &kvpb.GetResponse{}, err
	}
	return &kvpb.GetResponse{
		Value: value,
	}, nil
}

func (s *Server) Put(ctx context.Context, req *kvpb.PutRequest) (*kvpb.PutResponse, error) {
	err := s.store.Write(ctx, req.Key, req.Value)
	if err != nil {
		return &kvpb.PutResponse{}, err
	}
	return &kvpb.PutResponse{}, nil
}

func (s *Server) Delete(ctx context.Context, req *kvpb.DeleteRequest) (*kvpb.DeleteResponse, error) {
	return &kvpb.DeleteResponse{}, nil
}
