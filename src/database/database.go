package database

import (
	"errors"
)

const MongoDB = "mongodb"

var ErrNoDocuments = errors.New("document does not exist")
var ErrDuplicateKey = errors.New("document already exists")

type Iterator[T any] interface {
	Chan() chan T
	Err() error
}

func SliceBuf[T any](iter Iterator[T], bufSize int64) ([]T, error) {
	entries := make([]T, 0, bufSize)
	for entry := range iter.Chan() {
		entries = append(entries, entry)
	}
	return entries, iter.Err()
}
