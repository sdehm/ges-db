package store

import (
	"testing"
)

func TestSetGet(t *testing.T) {
	s, err := NewStoreMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	k, err := s.Set([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	value, err := s.Get(k)
	if err != nil {
		t.Fatal(err)
	}
	if string(value) != "hello" {
		t.Fatal("value not equal")
	}
}