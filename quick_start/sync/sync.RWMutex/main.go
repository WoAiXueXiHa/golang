package main

import (
	"fmt"
	"sync"
)

// RWMutex 读写锁：读多写少
// 多个读可以同时进行
// 写的时候，不能读也不能写

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.data[key]
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

func main() {
	c := NewCache()
	c.Set("name", "LearnQ")
	fmt.Println(c.Get("name"))
}

// 死锁场景
// 1. 加锁后忘记解锁，后面的 goroutine 永远拿不到锁
// 2. 重复加同一把锁
// 3. 两把锁相互等待
// // goroutine 1
// a.Lock()
// b.Lock()

// // goroutine 2
// b.Lock()
// a.Lock()
