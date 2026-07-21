package main

import (
	"fmt"
	"sync"
)

// Mutex 互斥锁：保护共享变量，同一时刻只允许一个 goroutine 修改

func main() {
	var mu sync.Mutex
	var count int
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mu.Lock()
			defer mu.Unlock()
			count++
		}()
	}
	wg.Wait()
	fmt.Println(count)
}
