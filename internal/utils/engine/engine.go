package main

import (
	"bytes"
	"fmt"
	"log"
)

func main() {
	data := []byte("Hello!")

	reader := bytes.NewReader(data)
	buf := make([]byte, 4)
	n, err := reader.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	buf2 := make([]byte, 2)
	_, err = reader.Read(buf2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("N:", n)
	fmt.Printf("Buf: %s\n", string(buf))
	fmt.Printf("Buf2: %s\n", string(buf2))
	fmt.Printf("Data: %s\n", string(data))

}
