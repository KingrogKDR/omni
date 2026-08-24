package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/KingrogKDR/omni/internal/client"
)

const ServerAddr string = "localhost:28051"

func TestPutKey(t *testing.T) {
	ctx := context.Background()
	client, err := client.NewClient(ServerAddr)
	if err != nil {
		t.Fatal("Failed creating client:", err)
	}
	defer client.Close()

	_, err = client.Put(ctx, "Name", "Bob")
	if err != nil {
		t.Fatal("client error:", err)
	}

	fmt.Println("OK")
}

func TestGetKey(t *testing.T) {
	client, err := client.NewClient(ServerAddr)
	if err != nil {
		t.Fatal("Failed creating client:", err)
	}
	defer client.Close()

	resp, err := client.Get(context.Background(), "Name")
	if err != nil {
		t.Fatal("client error:", err)
	}

	fmt.Printf("Got value: %s\n", resp.Value)
}
