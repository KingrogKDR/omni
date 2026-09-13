package main

import (
	"log"

	wal "github.com/KingrogKDR/kWALity"
	"github.com/KingrogKDR/omni/internal/storage"
	"github.com/KingrogKDR/omni/internal/utils/engine"
)

func main() {
	opts := wal.DefaultWALOptions()
	_, err := wal.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	_, err = engine.NewWALEntry(storage.Put{
		CF:  []byte("cf-1"),
		Key: []byte("name"),
		Val: []byte("Alice"),
	})
	if err != nil {
		log.Fatal(err)
	}
}
