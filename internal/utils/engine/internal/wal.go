package engine

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type WAL struct {
	entries []WALEntry
}

type Transaction uint8

const (
	PutTransaction Transaction = iota + 1
	DeleteTransaction
)

type WALEntry struct {
	transaction Transaction
	data        TransactionData
}

type TransactionData interface {
	isTransaction()
}

type PutData struct {
	key []byte
	val []byte
}

func (PutData) isTransaction() {}

type DeleteData struct {
	key []byte
}

func (DeleteData) isTransaction() {}

func NewWALEntry(data TransactionData) (WALEntry, error) {
	switch d := data.(type) {
	case PutData:
		return WALEntry{
			transaction: PutTransaction,
			data:        PutData{key: bytes.Clone(d.key), val: bytes.Clone(d.val)},
		}, nil
	case DeleteData:
		return WALEntry{
			transaction: DeleteTransaction,
			data:        DeleteData{key: bytes.Clone(d.key)},
		}, nil
	default:
		return WALEntry{}, errors.New("unsupported data")

	}
}

var byteOrder = binary.LittleEndian

func Encode(entry WALEntry) ([]byte, error) {
	var buf []byte
	switch d := entry.data.(type) {
	case PutData:
		buf = append(buf, byte(PutTransaction))
		buf = byteOrder.AppendUint64(buf, uint64(len(d.key)))
		buf = append(buf, d.key...)
		buf = byteOrder.AppendUint64(buf, uint64(len(d.val)))
		buf = append(buf, d.val...)
	case DeleteData:
		buf = append(buf, byte(DeleteTransaction))
		buf = byteOrder.AppendUint64(buf, uint64(len(d.key)))
		buf = append(buf, d.key...)
	default:
		return nil, errors.New("invalid data for encoding: Requires KEY VALUE for PUT and KEY for DELETE")
	}

	return buf, nil
}

func Decode(data []byte) (WALEntry, error) {
	if len(data) < 1 {
		return WALEntry{}, errors.New("truncated WAL entry: missing WAL data")
	}
	reader := bytes.NewReader(data)
	transactionByte, err := reader.ReadByte()
	if err != nil {
		return WALEntry{}, fmt.Errorf("truncated WAL entry: missing transaction type: %w", err)
	}
	transaction := Transaction(transactionByte)

	entry := WALEntry{}
	entry.transaction = transaction

	switch transaction {
	case PutTransaction:
		var keyLenBytes [8]byte
		_, err = io.ReadFull(reader, keyLenBytes[:])
		if err != nil {
			return WALEntry{}, fmt.Errorf("truncated WAL entry: missing key length: %w", err)
		}
		keyLen := byteOrder.Uint64(keyLenBytes[:])
		if keyLen > uint64(reader.Len()) {
			return WALEntry{}, errors.New("truncated WAL entry: key length is larger than the rest of the data")
		}

		key := make([]byte, keyLen)
		_, err = io.ReadFull(reader, key)
		if err != nil {
			return WALEntry{}, fmt.Errorf("truncated WAL entry: incomplete key: %w", err)
		}

		var valLenBytes [8]byte
		_, err = io.ReadFull(reader, valLenBytes[:])
		if err != nil {
			return WALEntry{}, fmt.Errorf("truncated WAL entry: missing value length: %w", err)
		}
		valLen := byteOrder.Uint64(valLenBytes[:])
		if valLen > uint64(reader.Len()) {
			return WALEntry{}, errors.New("truncated WAL entry: value length is larger than the rest of the data")
		}

		val := make([]byte, valLen)
		_, err = io.ReadFull(reader, val)
		if err != nil {
			return WALEntry{}, fmt.Errorf("truncated WAL entry: incomplete value: %w", err)
		}
		entry.data = PutData{
			key: key,
			val: val,
		}
	case DeleteTransaction:
		var keyLenBytes [8]byte
		_, err = io.ReadFull(reader, keyLenBytes[:])
		if err != nil {
			return WALEntry{}, fmt.Errorf("truncated WAL entry: missing key length: %w", err)
		}
		keyLen := byteOrder.Uint64(keyLenBytes[:])
		if keyLen > uint64(reader.Len()) {
			return WALEntry{}, errors.New("truncated WAL entry: key length is larger than the rest of the data")
		}

		key := make([]byte, keyLen)
		_, err = io.ReadFull(reader, key)
		if err != nil {
			return WALEntry{}, fmt.Errorf("truncated WAL entry: incomplete key: %w", err)
		}

		entry.data = DeleteData{
			key: key,
		}
	default:
		return WALEntry{}, errors.New("invalid transaction type")
	}

	return entry, nil
}
