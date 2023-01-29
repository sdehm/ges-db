package store

import "testing"

// test struct to encode
type message struct {
	Id  int
	Text string
}

func TestEncodeDecode(t *testing.T) {
	msg := message{Id: 1, Text: "hello world"}
	data, err := Encode(msg)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("data is empty")
	}

	decoded, err := Decode[message](data)
	if err != nil {
		t.Fatal(err)
	}
	if decoded != msg {
		t.Fatal("decoded message not equal")
	}
}

