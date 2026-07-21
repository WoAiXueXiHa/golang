package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 对简单变量做原子操作：计数器、开关状态
func main() {
	var count atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count.Add(1)
		}()
	}

	wg.Wait()
	fmt.Println(count.Load())
}

// Add(1)      // 原子加
// Load()      // 原子读
// Store(v)    // 原子写
// Swap(v)     // 替换并返回旧值
// CompareAndSwap(old, new) // 如果等于 old，就改成 new
