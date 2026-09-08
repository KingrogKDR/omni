package engine

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Wal struct {
	entries []WalEntry
}

type Transaction uint8

const (
	PutTransaction Transaction = iota + 1
	DeleteTransaction
)

type WalEntry struct {
	lsn         uint64
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

func NewWalEntry(t Transaction, data TransactionData) (*WalEntry, error) {
	entry := &WalEntry{
		transaction: t,
		data:        data,
	}
	switch d := data.(type) {
	case PutData:
		entry.data = PutData{
			key: bytes.Clone(d.key),
			val: bytes.Clone(d.val),
		}

	case DeleteData:
		entry.data = DeleteData{
			key: bytes.Clone(d.key),
		}
	default:
		return nil, errors.New("unsupported data")

	}

	return entry, nil
}

var byteOrder = binary.LittleEndian

func Encode(entry WalEntry) ([]byte, error) {
	var buf []byte
	switch entry.transaction {
	case PutTransaction:
		buf = append(buf, 1)

	case DeleteTransaction:
		buf = append(buf, 2)

	default:
		return nil, errors.New("unsupported transaction")
	}
	switch d := entry.data.(type) {
	case PutData:
		buf = byteOrder.AppendUint64(buf, uint64(len(d.key)))
		buf = append(buf, d.key...)
		buf = byteOrder.AppendUint64(buf, uint64(len(d.val)))
		buf = append(buf, d.val...)
	case DeleteData:
		buf = byteOrder.AppendUint64(buf, uint64(len(d.key)))
		buf = append(buf, d.key...)
	default:
		return nil, errors.New("Invalid data for encoding: Requires KEY VALUE for PUT and KEY for DELETE")
	}

	return buf, nil
}

// TODO: Replace with a bytes Reader and also uatomatically trim trailing bytes in encoder and hold a check in decoder
func Decode(data []byte) (*WalEntry, error) {
	if len(data) < 1 {
		return nil, errors.New("data is too short")
	}
	entry := &WalEntry{}

	transactionType := data[0]

	switch transactionType {
	case 1:
		entry.transaction = PutTransaction

	case 2:
		entry.transaction = DeleteTransaction

	default:
		return nil, fmt.Errorf("unknown transaction type: %d", transactionType)
	}

	data = data[1:]
	if len(data) < 8 {
		return nil, errors.New("missing key length")
	}

	keyLen := byteOrder.Uint64(data[:8])
	data = data[8:]

	if keyLen > uint64(len(data)) {
		return nil, errors.New("invalid key length")
	}

	key := bytes.Clone(data[:keyLen])
	data = data[keyLen:]

	switch entry.transaction {
	case PutTransaction:
		if len(data) < 8 {
			return nil, errors.New("missing value length")
		}

		valLen := byteOrder.Uint64(data[:8])
		data = data[8:]

		if valLen > uint64(len(data)) {
			return nil, errors.New("invalid value length")
		}

		val := bytes.Clone(data[:valLen])

		entry.data = PutData{
			key: key,
			val: val,
		}
	case DeleteTransaction:
		entry.data = DeleteData{
			key: key,
		}
	}

	return entry, nil
}
