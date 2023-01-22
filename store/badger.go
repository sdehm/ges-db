package store

import (
	"context"
	"encoding/binary"
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/ristretto/z"
)

type Store struct {
	db  *badger.DB
	seq *badger.Sequence
}

func NewStore(path string) (*Store, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	seq, err := db.GetSequence([]byte("seq"), 1000)
	if err != nil {
		return nil, err
	}
	return &Store{
		db:  db,
		seq: seq,
	}, nil
}

func (b *Store) Close() {
	b.seq.Release()
	b.db.Close()
}

func (b *Store) Set(value []byte) error {
	s, err := b.seq.Next()
	if err != nil {
		return err
	}
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, s)
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (b *Store) Get(key string) ([]byte, error) {
	var value []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		var v []byte
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			v = val
			return nil
		})
		value = v
		return nil
	})
	return value, err
}

func (b *Store) Clear() error {
	return b.db.DropAll()
}

func (b *Store) Iterate() {
	b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := binary.LittleEndian.Uint64(item.Key())
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%d, value=%s\n", k, v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (b *Store) Stream() error {
	stream := b.db.NewStream()

	stream.Send = func(buf *z.Buffer) error {
		list, err := badger.BufferToKVList(buf)
		if err != nil {
			return err
		}
		for _, kv := range list.Kv {
			fmt.Printf("key=%d, value=%s\n", binary.LittleEndian.Uint64(kv.Key), kv.Value)
		}
		return nil
	}

	return stream.Orchestrate(context.Background())
}
