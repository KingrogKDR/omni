package main

import (
	"fmt"
	"log"
	"net"

	"github.com/KingrogKDR/omni/internal/server"
	"github.com/KingrogKDR/omni/internal/storage"
	kvpb "github.com/KingrogKDR/omni/proto/gen/kv"
	"github.com/dgraph-io/badger"
	"google.golang.org/grpc"
)

const ServerAddr string = "localhost:28051"

func main() {
	opts := badger.DefaultOptions("./data")
	badger_store, err := storage.NewBadgerStore(opts)
	if err != nil {
		log.Fatalf("Error opening store: %v", err)
	}
	defer badger_store.DB.Close()

	omniServer := server.NewServer(badger_store)

	grpcServer := grpc.NewServer()
	kvpb.RegisterOmniServer(grpcServer, omniServer)
	lis, err := net.Listen("tcp", ServerAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("grpc server listening on", ServerAddr)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}

	// omniClient := client.NewClient(ServerAddr)
	// client.Put()

}
