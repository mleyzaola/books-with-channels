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
	// introduced this for/loop to show how we are not really caching database
	for j := 0; j < 10; j++ {
		for i := 0; i < 10; i++ {
			id := rnd.Intn(10) + 1
			wg.Add(2)
			go func(id int, wg *sync.WaitGroup) {
				if b, ok := queryCache(id); ok {
					fmt.Println("from cache")
					fmt.Println(b)
				}
				wg.Done()
			}(id, wg)
			go func(id int, wg *sync.WaitGroup) {
				if b, ok := queryDatabase(id); ok {
					fmt.Println("from database")
					fmt.Println(b)
				}
				wg.Done()
			}(id, wg)
			time.Sleep(time.Millisecond * 150)
		}
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
