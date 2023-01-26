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

func TestSetGet2Values(t *testing.T) {
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

	k2, err := s.Set([]byte("world"))
	if err != nil {
		t.Fatal(err)
	}

	value, err = s.Get(k2)
	if err != nil {
		t.Fatal(err)
	}
	if string(value) != "world" {
		t.Fatal("value not equal")
	}

	// check first value again
	value, err = s.Get(k)
	if err != nil {
		t.Fatal(err)
	}
	if string(value) != "hello" {
		t.Fatal("value not equal")
	}
}