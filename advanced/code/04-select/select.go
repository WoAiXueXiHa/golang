package main

import (
	"fmt"
	"time"
)

// go 的 select 和 Linux 的 select 不一样

// =================================================
// 1. select 的基本使用
// select {
//		case <-channel1:	// 从1中读取成功，执行
//			do ...
//		case channel2<- 2:  // 向2中写入2成功，执行
// 			do ...
//		default:
//			do ...
//}
// 空 select 永久阻塞
// =================================================

// =================================================
// 2. 没有default且case无法执行的select永久阻塞
// =================================================
// fatal error: all goroutines are asleep - deadlock!
func demo1() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	select {
	// 读数据，但是两个ch都没有数据，而且没有goroutine往里面写数据
	// 不可能读到数据，死锁
	case <-ch1:
		fmt.Println("receiced from ch1")
	case num := <-ch2:
		fmt.Printf("num is %d\n", num)
	}
}

// 2. 有单一case和default的select
func demo2() {
	ch := make(chan int, 1)
	select {
	case <-ch:
		fmt.Println("recevied from ch")
	default:
		fmt.Println("default....")
	}
}

// 3. 有多个case和default的select
func demo3() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	go func() {
		time.Sleep(1*time.Second)
		for i := 0; i < 3; i++ {
			select {
			case v := <-ch1:
				fmt.Println("recevied from ch1: ", v)
			case v := <-ch2:
				fmt.Println("recevied from ch2: ", v)
			default:
				fmt.Println("default....")
			}
		}
	}()
	ch1<- 1
	time.Sleep(1*time.Second)
	ch2<- 2
	time.Sleep(1*time.Second)
}

// 选择是随机执行的
// 两种 case 都可能打印
func demo4() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	ch1<- 66
	ch2<- 11

	select {
	case v := <-ch1:
		fmt.Println("recevied from ch1: ", v)
	case v:= <-ch2:
		fmt.Println("recevied from ch2: ", v)
	default:
		fmt.Println("default...")
	}

}

func main() {
	// demo1()
	// demo2()
	// demo3()
	demo4()
}
