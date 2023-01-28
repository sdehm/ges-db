package store

import (
	"fmt"
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

func TestIterate(t *testing.T) {
	s, err := NewStoreMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	k, err := s.Set([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	k2, err := s.Set([]byte("world"))
	if err != nil {
		t.Fatal(err)
	}
	keys := make([][]byte, 0)
	values := make([]string, 0)
	s.Iterate(func(k []byte, v []byte) {
		fmt.Printf("key=%d, value=%s \n", k, v)
		values = append(values, string(v))
		keys = append(keys, k)
	})
	if values[0] != "hello" {
		t.Fatalf("value not equal: %s", values[0])
	}
	if values[1] != "world" {
		t.Fatalf("value not equal: %s", values[1])
	}
	if string(keys[0]) != string(append([]byte("ges"), k...)) {
		t.Fatalf("key not equal: %d", keys[0])
	}
	if string(keys[1]) != string(append([]byte("ges"), k2...)) {
		t.Fatalf("key not equal: %d", keys[1])
	}
}