package server

import (
	"context"
	"fmt"

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
	val, err := s.store.Get(ctx, req.Key)
	if err != nil {
		return &kvpb.GetResponse{}, err
	}
	return &kvpb.GetResponse{
		Value: val,
	}, nil
}

func (s *Server) Put(ctx context.Context, req *kvpb.PutRequest) (*kvpb.PutResponse, error) {
	err := s.store.Put(ctx, req.Pair.Key, req.Pair.Value)
	if err != nil {
		return &kvpb.PutResponse{}, err
	}
	return &kvpb.PutResponse{}, nil
}

func (s *Server) Delete(ctx context.Context, req *kvpb.DeleteRequest) (*kvpb.DeleteResponse, error) {
	err := s.store.Delete(ctx, req.Key)
	if err != nil {
		return &kvpb.DeleteResponse{}, err
	}
	return &kvpb.DeleteResponse{}, nil
}

func (s *Server) List(req *kvpb.ListRequest, stream kvpb.Omni_ListServer) error {
	var limit uint32
	if req.Limit != nil {
		limit = *req.Limit
	}

	seq, iterErr := s.store.List(stream.Context(), req.Prefix, limit)

	for k, v := range seq {
		err := stream.Send(&kvpb.ListResponse{
			Pair: &kvpb.KeyValue{Key: k, Value: v},
		})

		if err != nil {
			return fmt.Errorf("list send: %w", err)
		}
	}

	if err := iterErr.Err(); err != nil {
		return fmt.Errorf("list iterrator: %w", err)
	}

	return nil
}
