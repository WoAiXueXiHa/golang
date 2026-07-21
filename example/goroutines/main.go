package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// goroutine 是用 go 关键字启动一个并发执行的函数

// 1. 基本使用
// func Say() {
// 	fmt.Println("hello from goroutine")
// }
// func main() {
// 	go Say()

// 	fmt.Println("hello from main")
// 	time.Sleep(time.Second)
// 	// main 本身也是一个 goroutine，main结束，整个程序结束
// 	// 没执行的 goroutine 直接被杀掉
// }

// // 2. WaitGroup 正确等待 goroutine 结束
// func main() {
// 	var wg sync.WaitGroup

// 	for i := 1; i <= 3; i++ {
// 		wg.Add(1) // 登记一个任务，开 goroutine 前加

// 		go func(id int) {
// 			defer wg.Done() // 任务完成，计数-1，一个 goroutine 结束时减
// 			fmt.Println("worker:", id)
// 		}(i)
// 	}

// 	wg.Wait()	// main 等待所有 goroutine 结束
// 	fmt.Println("all done~")
// }

// 4. 创建一万个 goroutine
func demo() {
	var wg sync.WaitGroup
	var cnt int64

	n := 10000

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			atomic.AddInt64(&cnt, 1) // 保证数据安全
		}()
	}
	wg.Wait()
	fmt.Println("cnt=", cnt)
}

func main() {
	var wg sync.WaitGroup
	// 全是 5
	// 原因：闭包捕获的是这个 i 变量本身
	// 创建了 5 个goroutine 没来的及执行，i 已经变成5了
	var i int
	// for i = 0; i < 5; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		fmt.Println(i)
	// 	}()
	// }

	// 正确方式：把循环变量当成参数传进去
	for i = 0; i < 5; i++ {
		wg.Add(1)

		go func(x int) {
			defer wg.Done()
			fmt.Println(x)
		}(i)
	}
	demo()
	wg.Wait()

}
