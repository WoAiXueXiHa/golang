package main

import (
	"fmt"
	"sync"
)

// 普通 map 并发读写会出问题
func demo() {
	m := map[string]int{}
	go func() {
		m["a"] = 1
	}()

	go func() {
		fmt.Println(m["a"])
	}()
}

func main() {
	// demo()

	var m sync.Map

	m.Store("go", 100)
	m.Store("redis", 90)

	v, ok := m.Load("go")
	if ok {
		fmt.Println(v.(int))
	}

	m.Delete("redis")

	m.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})
}

// Store(key, value)  // 写入
// Load(key)          // 读取
// Delete(key)        // 删除
// LoadOrStore(k, v)  // 有就读，没有就写
// Range(func)        // 遍历

// 使用场景：
// key 写一次，读很多次：缓存
// 多个 goroutine 操作不同的 key
