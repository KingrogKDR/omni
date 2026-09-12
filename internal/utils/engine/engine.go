package main

import (
	"log"

	wal "github.com/KingrogKDR/kWALity"
)

func main() {
	opts := wal.DefaultWALOptions()
	_, err := wal.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

}
