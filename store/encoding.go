package store

import (
	"bytes"
	"encoding/gob"
)

func Encode[T any](v T) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode[T any](data []byte) (T, error) {
	buf := bytes.Buffer{}
	buf.Write(data)
	dec := gob.NewDecoder(&buf)
	var v T
	err := dec.Decode(&v)
	if err != nil {
		return v, err
	}
	return v, nil
}