package main

import (
	"fmt"
	"sync"
)

// channel 是 goroutine 之间传数据 + 做同步的管道
// 发送和接收都可能阻塞

// 发送：把数据扔到 chan 里，所以：ch<- 10 把 10 扔到 ch里
// 接收：从 chan 里拿出来数据，所以：x := <-ch 从 ch 里拿出数据，赋值给 x
// 不要这个值，等一下：<-ch，从 ch 接收一个值，但是我不要他
// 箭头指向谁，数据就流向谁

// 1. 什么时候会阻塞？
// 无缓冲 channel：发送和接收必须同时出现
// 无缓冲像当面传纸条：
// goroutine A          channel          goroutine B
//    ch <- 10   --->   [无缓存]   --->      <- ch
// 发送方必须等接收方伸手接
// 接收方也必须等发送方递过来

// 有缓冲 channel：缓冲区没满，发送不阻塞；缓冲区不空，接收不阻塞

func dmeo0() {
	ch1 := make(chan int)

	go func() {
		ch1 <- 10 // 没人接收，卡在这
	}()

	x := <-ch1 // 没人发送，卡在这
	fmt.Println(x)

	ch2 := make(chan int, 2)
	ch2 <- 1 // 不阻塞
	ch2 <- 2 // 不阻塞
	// ch2<- 3 阻塞，缓冲区满了，写不进
	fmt.Println(<-ch2) // 1
	fmt.Println(<-ch2) // 2
	// fmt.Println(<-ch2) 阻塞，缓冲区空了，读不到
}

// 2. channel 关闭之后会怎么样
// 关闭后，还能把关闭前剩余的数据读出来，读完之后，再读会得到 0+false
// 往已经关闭的 channel 里写入，会 panic
// 重复 close 也会 panic
func demo1() {
	ch := make(chan int, 2)

	ch <- 1
	ch <- 0
	close(ch)
	// ch <- 2 // panic: send on closed channel
	// close(ch) // panic: close of closed channel

	v, ok := <-ch
	fmt.Println(v, ok)

	v, ok = <-ch
	fmt.Println(v, ok)

	v, ok = <-ch
	fmt.Println(v, ok)
	// 用 ok 可以判断出这个0值到底是原来缓冲区里的，还是关闭后空的读到的
}

// 2. for range 读取 channel
// 一直读 channel，直到 channel 被关闭且数据读完
func demo2() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 1
	ch <- 1
	close(ch)

	for v := range ch {
		fmt.Println(v)
	}
	// 如果不关闭，for range 会一直等
	// fatal error: all goroutines are asleep - deadlock!
	// ch1 := make(chan int, 5)
	// go func() {
	// 	ch1 <- 2
	// 	ch1 <- 2
	// 	ch1 <- 2
	// }()
	// for v := range ch1 {
	// 	fmt.Println(v)
	// }
}

// 3. Go 不通过共享内存来通信，而是通过通信来共享内存
// 把数据通过 channel 交给某个 goroutine，谁拿到数据，谁处理数据
// 大家不一起抢同一块内存，而是通过 channel 把数据所有权交出去
func demo3() {
	ch := make(chan int)
	go func() {
		ch <- 10 // 把数据交出去
	}()

	x := <-ch // main 拿到了数据
	fmt.Println(x)
}

// 用 channel 实现 goroutine 之间的锁
func demo4() {
	// 空结构体不占用内存
	lock := make(chan struct{}, 1)
	var wg sync.WaitGroup

	cnt := 0

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			lock <- struct{}{} // 加锁：放入令牌，满了就阻塞
			cnt++              // 临界区
			<-lock             // 解锁：取出令牌
		}()
	}
	wg.Wait()
	fmt.Println(cnt)
}

func main() {
	// demo1()
	// demo2()
	demo4()
}
