package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/sdehm/ges-db/store"
)

func main() {
	badger, err := store.NewStore("/tmp/badger")
	if err != nil {
		log.Fatal(err)
	}
	defer badger.Close()
	_ = badger.Clear()
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			_, _ = badger.Set([]byte(fmt.Sprintf("value%d", i)), "test")
		}()
	}
	wg.Wait()

	// badger.Iterate(func(k []byte, v []byte) {
	// 	fmt.Printf("key=%d, value=%s\n", k, v)
	// })

	fmt.Println("streaming")
	_ = badger.Stream(func(k []byte, v []byte) {
		fmt.Printf("key=%d, value=%s\n", k, v)
	})

}
