package kvstore

import (
	"bufio"
	"fmt"
	"strings"
)

// reads and parses input from user
func ReadInput(reader *bufio.Reader) ([]string, uint32, error) {
	fmt.Print("> ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return nil, 0, fmt.Errorf("couldn't read string: %v\n", err)
	}
	input := text[:len(text)-1]
	inputSlice := strings.Split(input, " ")
	inputArgc := uint32(len(inputSlice))

	return inputSlice, inputArgc, nil
}
