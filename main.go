package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var cache = map[int]Book{}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	wg := &sync.WaitGroup{}
	cacheCh := make(chan Book)
	dbCh := make(chan Book)
	for i := 0; i < 10; i++ {
		id := rnd.Intn(10) + 1
		wg.Add(2)
		// send only channel
		go func(id int, wg *sync.WaitGroup, ch chan<- Book) {
			if b, ok := queryCache(id); ok {
				ch <- b
			}
			wg.Done()
		}(id, wg, cacheCh)
		// send only channel
		go func(id int, wg *sync.WaitGroup, ch chan<- Book) {
			if b, ok := queryDatabase(id); ok {
				ch <- b
			}
			wg.Done()
		}(id, wg, dbCh)

		// one goroutine to receive results from cache and database channels
		go func(cacheCh, dbCh <-chan Book) {
			select {
			case b := <-cacheCh:
				fmt.Println("from cache")
				fmt.Println(b)
				// because reading from cache is not always going to execute,
				// we need to force wait for db channel execution,
				// which always will complete
				// this will keep channels in sync
				<-dbCh
			case b := <-dbCh:
				fmt.Println("from database")
				fmt.Println(b)
			}
		}(cacheCh, dbCh)
		time.Sleep(time.Millisecond * 150)
	}
	// wait exactly for all wait groups to complete
	wg.Wait()
}

func queryCache(id int) (Book, bool) {
	b, ok := cache[id]
	return b, ok
}

func queryDatabase(id int) (Book, bool) {
	time.Sleep(time.Millisecond * 100)
	for _, b := range books {
		if b.ID == id {
			// now we have collisions between goroutines accessing the same shared memory
			cache[id] = b
			return b, true
		}
	}
	return Book{}, false
}
