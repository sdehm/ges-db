package store

import (
	"context"
	"encoding/binary"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/ristretto/z"
)

type Store struct {
	db  *badger.DB
	seq *badger.Sequence
}

func uint64ToBytes(u uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, u)
	return b
}

func bytesToUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func makeKey(p string, t string, i uint64) []byte {
	return append([]byte(p), append([]byte(t), uint64ToBytes(i)...)...)
}

func keyToId(key []byte) uint64 {
	return bytesToUint64(key[len(key)-8:])
}

func NewStore(path string) (*Store, error) {
	return newStore(badger.DefaultOptions(path))
}

func NewStoreMemory() (*Store, error) {
	return newStore(badger.DefaultOptions("").WithInMemory(true))
}

func newStore(opts badger.Options) (*Store, error) {
	db, err := badger.Open(opts)
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

func (b *Store) Set(value []byte, t string) ([]byte, error) {
	s, err := b.seq.Next()
	if err != nil {
		return nil, err
	}
	k := makeKey("event", t, s)
	e := badger.NewEntry(k, value)
	err =  b.db.Update(func(txn *badger.Txn) error {
		// return txn.Set(k, value)
		return txn.SetEntry(e)
	})
	return k, err
}

func (b *Store) Get(key []byte) ([]byte, error) {
	var value []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
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

func (b *Store) Iterate(f func (key []byte, value []byte)) {
	b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("event")
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				f(item.Key(), v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (b *Store) Stream(f func (key []byte, value []byte)) error {
	stream := b.db.NewStream()
	stream.Prefix = []byte("event")

	stream.Send = func(buf *z.Buffer) error {
		list, err := badger.BufferToKVList(buf)
		if err != nil {
			return err
		}
		for _, kv := range list.Kv {
			f(kv.Key, kv.Value)
		}
		return nil
	}

	return stream.Orchestrate(context.Background())
}
