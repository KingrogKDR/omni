package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var KV map[string]string = make(map[string]string)

func Put(key, value string) {
	KV[key] = value
}

func Get(key string) string {
	value, ok := KV[key]
	if !ok {
		return "Key not found"
	}
	return value
}

func View() {
	for key, value := range KV {
		fmt.Printf("%s: %s\n", key, value)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("couldn't read string: %v\n", err)
			continue
		}
		input := text[:len(text)-1]
		inputSlice := strings.Split(input, " ")
		inputArgc := len(inputSlice)
		switch inputSlice[0] {
		case "add":
			if inputArgc != 3 {
				fmt.Printf("invalid no. of arguments: got %d expected 3\nUsage: add KEY VALUE\n", inputArgc)
				break
			}
			Put(inputSlice[1], inputSlice[2])
			fmt.Println("Ok")
		case "get":
			if inputArgc != 2 {
				fmt.Printf("invalid no. of arguments: got %d expected 2\nUsage: get KEY\n", inputArgc)
				break
			}
			fmt.Printf("%s: %s\n", inputSlice[1], Get(inputSlice[1]))
		case "view":
			if inputArgc != 1 {
				fmt.Printf("invalid no. of arguments: got %d expected 1\nUsage: get KEY\n", inputArgc)
				break
			}
			View()
		default:
			fmt.Println("command not found\navailable commands: Usage\n add: add KEY VALUE\n get: get KEY\n view: view")
		}
	}
}
