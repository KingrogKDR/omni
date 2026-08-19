package main

import (
	"log"
	"net"

	"github.com/KingrogKDR/omni/internal/server"
	kvpb "github.com/KingrogKDR/omni/proto/gen/kv"
	"google.golang.org/grpc"
)

const storeAddr = "localhost:28051"

func main() {
	omniServer := server.NewServer()

	grpcServer := grpc.NewServer()
	kvpb.RegisterOmniServer(grpcServer, omniServer)
	lis, err := net.Listen("tcp", storeAddr)
	if err != nil {
		log.Fatal(err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
