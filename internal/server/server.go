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
	readReq := storage.ReadRequest{
		Type: storage.GET,
		Key:  req.Key,
	}
	value, err := s.store.Read(ctx, readReq)
	if err != nil {
		return &kvpb.GetResponse{}, err
	}
	return &kvpb.GetResponse{
		Value: value,
	}, nil
}

func (s *Server) Put(ctx context.Context, req *kvpb.PutRequest) (*kvpb.PutResponse, error) {
	writeReq := storage.WriteRequest{
		Type:  storage.PUT,
		Key:   req.Key,
		Value: req.Value,
	}

	err := s.store.Write(ctx, writeReq)
	if err != nil {
		return &kvpb.PutResponse{}, err
	}
	return &kvpb.PutResponse{}, nil
}

func (s *Server) Delete(ctx context.Context, req *kvpb.DeleteRequest) (*kvpb.DeleteResponse, error) {
	writeReq := storage.WriteRequest{
		Type: storage.DELETE,
		Key:  req.Key,
	}

	err := s.store.Write(ctx, writeReq)
	if err != nil {
		return &kvpb.DeleteResponse{}, err
	}
	return &kvpb.DeleteResponse{}, nil
}

func (s *Server) List(ctx context.Context, req *kvpb.ListRequest) (*kvpb.ListResponse, error) {
	resp, err := s.store.List(ctx)
	if err != nil {
		return &kvpb.ListResponse{}, err
	}
	return &kvpb.ListResponse{
		KeyValue: resp.KeyValue,
	}, nil
}
