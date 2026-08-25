package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/KingrogKDR/omni/internal/server"
	"github.com/KingrogKDR/omni/internal/storage"
	kvpb "github.com/KingrogKDR/omni/proto/gen/kv"
	"github.com/dgraph-io/badger"
	"google.golang.org/grpc"
)

const ServerDefaultAddr string = "localhost:28051"

func main() {
	opts := badger.DefaultOptions("./data")
	store, err := storage.NewBadgerStore(opts)
	if err != nil {
		log.Fatalf("Error opening store: %v", err)
	}
	defer store.DB.Close()

	omniServer := server.NewServer(store)

	grpcServer := grpc.NewServer()
	kvpb.RegisterOmniServer(grpcServer, omniServer)
	lis, err := net.Listen("tcp", ServerDefaultAddr)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	go func() {
		log.Println("grpc server listening on", ServerDefaultAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Println("grpc server stopped:", err)
			stop <- syscall.SIGTERM
		}
	}()

	waitForShutdown(grpcServer, stop)
}

func waitForShutdown(grpcServer *grpc.Server, stop chan os.Signal) {
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	log.Println("shutting down gRPC server...")
	grpcServer.GracefulStop()
}
