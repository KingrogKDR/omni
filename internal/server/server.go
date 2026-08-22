package server

import (
	"context"

	kvpb "github.com/KingrogKDR/omni/proto/gen/kv"
)

type Server struct {
	kvpb.UnimplementedOmniServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Get(ctx context.Context, req *kvpb.GetRequest) (*kvpb.GetResponse, error) {
	return &kvpb.GetResponse{}, nil
}

func (s *Server) Put(ctx context.Context, req *kvpb.PutRequest) (*kvpb.PutResponse, error) {
	return &kvpb.PutResponse{}, nil
}

func (s *Server) Delete(ctx context.Context, req *kvpb.DeleteRequest) (*kvpb.DeleteResponse, error) {
	return &kvpb.DeleteResponse{}, nil
}
