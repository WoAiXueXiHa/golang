package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// case1: 有缓冲 channel 的 FIFO 语义
	ch := make(chan int, 3)
	for i := 1; i <= 3; i++ {
		ch <- i
	}
	fmt.Println("case1 len/cap:", len(ch), cap(ch))
	fmt.Println("case1 order:", <-ch, <-ch, <-ch)

	// case2: 无缓冲 channel 是同步的：发送者阻塞直到接收者就绪
	done := make(chan struct{})
	go func() {
		ch2 := make(chan int)
		go func() {
			time.Sleep(100 * time.Millisecond)
			ch2 <- 42 // 阻塞，直到下面的接收者就绪
			fmt.Println("case2 sender finished")
		}()
		v := <-ch2 // 先就绪，唤醒上面的发送者
		fmt.Println("case2 recv:", v)
		close(done)
	}()
	<-done

	// case3: 关闭后先取完 buf 里的剩余元素，再返回零值
	ch3 := make(chan int, 2)
	ch3 <- 10
	ch3 <- 20
	close(ch3)
	fmt.Println("case3 drain:", <-ch3, <-ch3)
	v, ok := <-ch3
	fmt.Println("case3 after drain:", v, ok)

	// case4: 读写 nil channel 永久阻塞（用 select + default 观察）
	var nilCh chan int
	select {
	case <-nilCh:
		fmt.Println("case4 unreachable")
	default:
		fmt.Println("case4 recv from nil: default taken")
	}
	select {
	case nilCh <- 1:
		fmt.Println("case4 unreachable")
	default:
		fmt.Println("case4 send to nil: default taken")
	}

	// case5: 向已关闭的 channel 写入 → panic
	ch5 := make(chan int)
	close(ch5)
	func() {
		defer func() { fmt.Println("case5 panic:", recover()) }()
		ch5 <- 1
	}()

	// case6: 双重 close → panic
	ch6 := make(chan int)
	close(ch6)
	func() {
		defer func() { fmt.Println("case6 panic:", recover()) }()
		close(ch6)
	}()

	// case7: close nil channel → panic
	var ch7 chan int
	func() {
		defer func() { fmt.Println("case7 panic:", recover()) }()
		close(ch7)
	}()

	// case8: select 随机公平性：两个恒就绪的 case，命中接近 1:1
	chA := make(chan int, 1)
	chB := make(chan int, 1)
	chA <- 1
	chB <- 1
	var a, b int
	for i := 0; i < 10000; i++ {
		select {
		case <-chA:
			a++
			chA <- 1
		case <-chB:
			b++
			chB <- 1
		}
	}
	fmt.Println("case8 ratio a:b =", a, b)

	// case9: select 全部阻塞时走 timeout
	select {
	case <-time.After(50 * time.Millisecond):
		fmt.Println("case9 timeout")
	}

	// case10: 8 发送者 + 8 接收者并发收发（-race 验证线程安全）
	var wg sync.WaitGroup
	ch10 := make(chan int, 16)
	for i := 0; i < 8; i++ {
		wg.Add(2)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				ch10 <- n*1000 + j
			}
		}(i)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				<-ch10
			}
		}()
	}
	wg.Wait()
	fmt.Println("case10 concurrent ok, leftover:", len(ch10))
}
