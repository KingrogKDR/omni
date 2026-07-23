package main

import (
	"fmt"

	"github.com/KingrogKDR/omni/internal/kvstore"
)

func main() {
	// Start with defaults
	db, err := kvstore.Open()
	if err != nil {
		panic(err)
	}
	err = db.Put("name", []byte("Alice"))
	if err != nil {
		panic(err)
	}

	v, err := db.Get("name")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Name: %s\n", v)

	err = db.Delete("name")
	if err != nil {
		panic(err)
	}

	err = db.View()
	if err != nil {
		panic(err)
	}

	err = db.Close()
	if err != nil {
		panic(err)
	}

	// Or with custom config
	custom := kvstore.Config{DataDir: "./custom/omni", InMemOnly: true}
	db2, err := kvstore.Open(custom)
	if err != nil {
		panic(err)
	}

	err = db2.Put("name", []byte("John"))
	if err != nil {
		panic(err)
	}

	v, err = db2.Get("name")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Name: %s\n", v)
	_ = db2
}
