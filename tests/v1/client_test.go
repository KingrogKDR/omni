package tests

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/KingrogKDR/omni/internal/client"
	"github.com/KingrogKDR/omni/internal/server"
	singleStorage "github.com/KingrogKDR/omni/internal/storage/single_storage"
	kvpb "github.com/KingrogKDR/omni/proto/gen/kv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func startTestClient(t *testing.T) *client.Client {
	t.Helper()

	store := singleStorage.NewSingleStorage(t.TempDir())
	if err := store.Start(); err != nil {
		t.Fatalf("failed to start storage: %v", err)
	}
	t.Cleanup(store.Stop)

	lis := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	kvpb.RegisterOmniServer(grpcServer, server.NewServer(store))
	go func() {
		_ = grpcServer.Serve(lis)
	}()
	t.Cleanup(grpcServer.Stop)

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		t.Fatalf("failed to dial bufconn: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return client.NewClientFromConn(conn)
}

func TestClientPutGet(t *testing.T) {
	c := startTestClient(t)
	ctx := context.Background()

	if err := c.Put(ctx, []byte("default"), []byte("k1"), []byte("v1")); err != nil {
		t.Fatalf("put failed: %v", err)
	}

	val, err := c.Get(ctx, []byte("default"), []byte("k1"))
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if string(val) != "v1" {
		t.Fatalf("got %q, want %q", val, "v1")
	}
}

func TestClientGetKeyNotFound(t *testing.T) {
	c := startTestClient(t)

	_, err := c.Get(context.Background(), []byte("default"), []byte("missing"))
	if !errors.Is(err, client.ErrKeyNotFound) {
		t.Fatalf("expected client.ErrKeyNotFound, got: %v", err)
	}
}

func TestClientDelete(t *testing.T) {
	c := startTestClient(t)
	ctx := context.Background()

	if err := c.Put(ctx, []byte("default"), []byte("k1"), []byte("v1")); err != nil {
		t.Fatalf("setup put failed: %v", err)
	}
	if err := c.Delete(ctx, []byte("default"), []byte("k1")); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err := c.Get(ctx, []byte("default"), []byte("k1"))
	if !errors.Is(err, client.ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound after delete, got: %v", err)
	}
}

func TestClientScanAll(t *testing.T) {
	c := startTestClient(t)
	ctx := context.Background()

	for _, k := range []string{"a", "b", "c"} {
		if err := c.Put(ctx, []byte("default"), []byte(k), []byte("v")); err != nil {
			t.Fatalf("setup put %s failed: %v", k, err)
		}
	}

	pairs, err := c.Scan(ctx, []byte("default"), &kvpb.ScanOptions{})
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(pairs))
	}
}
