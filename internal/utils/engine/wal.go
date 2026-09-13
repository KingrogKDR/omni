package engine

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/KingrogKDR/omni/internal/storage"
)

type Transaction uint8

const (
	PutTransaction Transaction = iota + 1
	DeleteTransaction
)

type WALEntry struct {
	transaction Transaction
	data        storage.WriteOp
}

func NewWALEntry(data storage.WriteOp) (WALEntry, error) {
	switch d := data.(type) {
	case storage.Put:
		return WALEntry{
			transaction: PutTransaction,
			data:        storage.Put{CF: bytes.Clone(d.CF), Key: bytes.Clone(d.Key), Val: bytes.Clone(d.Val)},
		}, nil
	case storage.Delete:
		return WALEntry{
			transaction: DeleteTransaction,
			data:        storage.Delete{Key: bytes.Clone(d.Key)},
		}, nil
	default:
		return WALEntry{}, errors.New("unsupported data")

	}
}

type TransactionCodec struct{}

func (TransactionCodec) Encode(data any) ([]byte, error) {
	entry, ok := data.(WALEntry)
	if !ok {
		return nil, errors.New("expected WALEntry")
	}

	return encodeTransaction(entry)
}

func (TransactionCodec) Decode(data []byte) (any, error) {
	return decodeTransaction(data)
}

var transactionByteOrder = binary.LittleEndian

func encodeTransaction(entry WALEntry) ([]byte, error) {
	var buf []byte
	switch d := entry.data.(type) {
	case storage.Put:
		buf = append(buf, byte(PutTransaction))
		buf = transactionByteOrder.AppendUint64(buf, uint64(len(d.CF)))
		buf = append(buf, d.CF...)
		buf = transactionByteOrder.AppendUint64(buf, uint64(len(d.Key)))
		buf = append(buf, d.Key...)
		buf = transactionByteOrder.AppendUint64(buf, uint64(len(d.Val)))
		buf = append(buf, d.Val...)
	case storage.Delete:
		buf = append(buf, byte(DeleteTransaction))
		buf = transactionByteOrder.AppendUint64(buf, uint64(len(d.Key)))
		buf = append(buf, d.Key...)
	default:
		return nil, errors.New("invalid data for encoding: Requires KEY VALUE for PUT and KEY for DELETE")
	}

	return buf, nil
}

func decodeTransaction(data []byte) (WALEntry, error) {
	if len(data) < 1 {
		return WALEntry{}, errors.New("truncated WAL entry: missing WAL data")
	}
	reader := bytes.NewReader(data)
	transactionByte, err := reader.ReadByte()
	if err != nil {
		return WALEntry{}, fmt.Errorf("truncated WAL entry: missing transaction type: %w", err)
	}
	transaction := Transaction(transactionByte)

	entry := WALEntry{
		transaction: transaction,
	}

	switch transaction {
	case PutTransaction:
		cf, err := readField(reader, "cf")
		if err != nil {
			return WALEntry{}, err
		}

		key, err := readField(reader, "key")
		if err != nil {
			return WALEntry{}, err
		}

		val, err := readField(reader, "value")
		if err != nil {
			return WALEntry{}, err
		}
		entry.data = storage.Put{
			CF:  cf,
			Key: key,
			Val: val,
		}
	case DeleteTransaction:
		key, err := readField(reader, "key")
		if err != nil {
			return WALEntry{}, err
		}

		entry.data = storage.Delete{
			Key: key,
		}
	default:
		return WALEntry{}, errors.New("invalid transaction type")
	}

	if reader.Len() != 0 {
		return WALEntry{}, errors.New("invalid WAL entry: trailing data")
	}

	return entry, nil
}

func readField(reader *bytes.Reader, field string) ([]byte, error) {
	var lenBytes [8]byte
	if _, err := io.ReadFull(reader, lenBytes[:]); err != nil {
		return nil, fmt.Errorf("truncated WAL entry: missing %s length: %w", field, err)
	}
	length := transactionByteOrder.Uint64(lenBytes[:])
	if length > uint64(reader.Len()) {
		return nil, fmt.Errorf("truncated WAL entry: %s length is larger than the rest of the data", field)
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, fmt.Errorf("truncated WAL entry: incomplete %s: %w", field, err)
	}
	return data, nil
}
