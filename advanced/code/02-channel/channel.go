package main

import (
	"fmt"
	"time"
)

// channel 是一个可以收发数据的管道
// 通信，那么就需要至少两个携程

// ===============================================================
// 1. 搞清楚什么时候会阻塞
// 当一个 channel 只有发送方没有接收方，或者只有接收方没有发送方的时候，就会阻塞
// ===============================================================
func demo1_block() {
	// ch := make(chan int)	// 无缓冲区的 channel

	// // 只发送，没有接收方
	// ch <- 45		// main 协程必须等到一个协程来接收数据

	// // 只接收，没有发送方
	// val := <-ch	// main 协程必须等到一个协程来发送数据
	// fmt.Println(val)
	
	// 正确处理方式
	ch := make(chan int)

	// 开启一个新的协程，接收方
	go func() {
		val := <-ch
		fmt.Println("接收到数据了：",val)
	}()

	// main 协程是发送方
	// 顺序不能错：必须先开启子协程（go func），保证救场的接收者有机会被调度，这里的代码是串行
	// 至于运行后是发送方先等、还是接收方先等，Go 运行时会自动协调，谁先到谁就原地等待。
	ch <- 42
	fmt.Println("程序结束")
}

// ===============================================================
// 2. 管道关闭之后还能继续读，但是读到的是零值
// ===============================================================
func demo2_after_close_recv() {
	// 因为这里有缓冲区，发的数据会放到缓冲区，保证不超过缓冲区大小就能这样写
	ch := make(chan string, 5)
	ch <- "hello"
	ch <- ","
	close(ch)
	
	go func() {
		for i := 0; i < 5; i++ {
			v := <-ch
			fmt.Printf("接收到的数据:%s\n", v)
		}
	}()
	time.Sleep(1 * time.Second)
}

// ===============================================================
// 3. 判定读取
// 如果我要写入零值，那么上面的读取方式就存在问题了
// 无法区分读到的是关闭后的，还是关闭前的
// ===============================================================
func demo3_judge_recv() {
	ch := make(chan int, 5)
	ch <- 1
	ch <- 200
	close(ch) 
	go func() {
		for i := 0; i < 5; i++ {
			v, ok := <-ch
			if ok {
				fmt.Printf("v=%d\n", v)
			} else {
				fmt.Printf("channel数据已经读完,v=%d\n", v)
			}
		}
	}()
	time.Sleep(1 * time.Second)
}

// ===============================================================
// 4. for range 读取
// 有数据就读取，直到另一端关闭了管道，for range 会感知
// ===============================================================
func demo4_forRange_recv() {
	ch := make(chan string, 5)
	ch <- "golang"
	ch <- "cpp"
	close(ch)	// 这里不关闭也是只读取这两份数据，协程不会退出，会阻塞等待 main 协程结束
	go func() {
		for v := range ch {
			fmt.Printf("v = %s\n", v)
		}
	}()
	time.Sleep(time.Second * 1)
}

// ===============================================================
// 5. 使用单向 channel
// 定义单向读：
// ch := make(chan int)
// type RecvChannel = <-chan int
// var recv RecvChannel = ch
//
// 定义单向写：
// ch := make(chan int)
// type SendChannel = chan<- int
// var send SendChannel = ch
// ===============================================================
type RecvChannel = <-chan int
type SendChannel = chan<- int

func demo5_simplex_channel() {
	var ch = make(chan int)

	go func() {
		var send SendChannel = ch
		fmt.Println("send: 100")
		send<- 100
	}()

	go func() {
		var recv RecvChannel = ch
		num := <-recv
		fmt.Printf("receive: %d\n", num)
	}()
	time.Sleep(1*time.Second)
}

// ===============================================================
// 6. golang中不以共享内存来通信，而以通信来共享内存
// C++中，多线程并发计算需要共享一个全局变量 sum，为了防止数据被同时修改，必须加锁
// 而 go 中，每个协程修改自己的局部sum，修改完之后通过 channel 来传递数据
// ===============================================================

func demo6_goroutine_comm() {
	Sum := func(s []int, c chan int) {
		sum := 0
		for _, v := range s {
			sum += v
		}
		c<- sum
	}

	s := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int)

	go func() {
		Sum(s[:len(s) / 2], c)	// 17
		time.Sleep(1 * time.Second)
	}()

	go Sum(s[len(s) / 2:], c)	// -5

	x, y := <-c, <-c

	fmt.Println(x, y, x+y)
} 

// ===============================================================
// 7. 实现 goroutine 之间的锁
// 利用缓冲队列满了之后，继续往 channel 里写数据，就会阻塞的特性
func demo7_mutex() {
	add := func(ch chan bool, num *int) {
		ch<- true
		*num = *num + 1
		<-ch
	}

	ch := make(chan bool, 1)

	var num int
	for i := 0; i < 100; i++ {
		go add(ch, &num)
	}

	time.Sleep(2 * time.Second)
	fmt.Printf("num的值: %d\n", num)
}

func main() {
	// demo1_block()
	// demo2_after_close_recv()
	// demo3_judge_recv()
	// demo4_forRange_recv()
	// demo5_simplex_channel()
	// demo6_goroutine_comm()
	demo7_mutex()
}
