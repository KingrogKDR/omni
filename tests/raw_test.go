package tests

import (
	"context"
	"fmt"
	"strings"
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

	_, err = client.Put(ctx, "Fruit", "Apple")
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

func TestDeleteKey(t *testing.T) {
	client, err := client.NewClient(ServerAddr)
	if err != nil {
		t.Fatal("Failed creating client:", err)
	}
	defer client.Close()

	_, err = client.Delete(context.Background(), "Name")
	if err != nil {
		t.Fatal("client error:", err)
	}

	fmt.Println("OK")
}

func TestListAllData(t *testing.T) {
	client, err := client.NewClient(ServerAddr)
	if err != nil {
		t.Fatal("Failed creating client:", err)
	}
	defer client.Close()

	resp, err := client.List(context.Background())
	if err != nil {
		t.Fatal("client error:", err)
	}

	fmt.Println(strings.Repeat("-", 32))
	fmt.Printf("%#v\n", resp.KeyValue)
	fmt.Println("OK")
}

func TestScanPrefix(t *testing.T) {
	client, err := client.NewClient(ServerAddr)
	if err != nil {
		t.Fatal("Failed creating client:", err)
	}
	defer client.Close()

	resp, err := client.Scan(context.Background(), "Name", 3)
	if err != nil {
		t.Fatal("client error:", err)
	}

	fmt.Println(strings.Repeat("-", 32))
	fmt.Printf("%#v\n", resp.KeyValue)
	fmt.Println("OK")
}
