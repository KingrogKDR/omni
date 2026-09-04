package utils

import "fmt"

type StorageError struct {
	Err error
	CF  []byte
	Key []byte
	Op  string
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("storage: %s cf=%s key=%s: %v", e.Op, e.CF, e.Key, e.Err)
}
func (e *StorageError) Unwrap() error { return e.Err }
