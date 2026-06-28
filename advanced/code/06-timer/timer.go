package main

import (
	"fmt"
	"time"
)

// Timer 是一种一次性时间定时器，在未来某个时刻，触发的事件只会执行一次
// type Timer struct {
//		c <-chan Time
// 		r runtineTimer
// }
// 这个Time类型的管道，用于事件通知
// 未达到设定时间的时候，管道内没有数据写入，一致阻塞
// 到达设定时间后，向管道内写入一个系统时间，触发事件

// 1. 创建Timer
func create() {
	timer := time.NewTimer(2 * time.Second) // 设置超时时间
	<-timer.C		// 从timer.C和管道读出数据，没有用变量接收
	fmt.Println("after 2s Time out")
}

// 2. 停止Timer
func cancel() {
	timer := time.NewTimer(2 * time.Second)
	res := timer.Stop()
	fmt.Println(res)	// true 没到超时时间停止 false 过了超时时间停止
}

// 3. 重置Timer
func reset() {
	timer := time.NewTimer(time.Second * 2)

	<-timer.C
	fmt.Println("time out1") 	// 2s 后打印

	res1 := timer.Stop()		// 已经过期了
	fmt.Printf("res1 is %t\n", res1)

	timer.Reset(time.Second)	// 重置超时时间

	res2 := timer.Stop()		// 这次没过期
	fmt.Printf("res2 is %t\n", res2)
}

// 4. time.AfterFunc
// func AfterFunc(d Duration, f func()) *Timer

func dur() {
	duration := time.Duration(1) * time.Second

	f := func() {
		fmt.Println("f has been called after 1s by time.AfterFunc")
	}

	timer := time.AfterFunc(duration, f)

	defer timer.Stop()

	time.Sleep(2 * time.Second)
}

// 5. time.After
// func After(d Duration) <-chan Time {
// return NewTimer(d).C
// }
// 返回timer里的管道，这个管道会在经过时段d之后写入数据
// 调用这个函数，相当于实现了定时器

func after() {
	ch := make(chan string)
	
	go func() {
		time.Sleep(time.Second * 3)
		ch<- "test"
	}()

	select {
	case val := <-ch:
		fmt.Printf("val is %s\n", val)
	case <-time.After(time.Second * 2):
		fmt.Println("timeout!!!")
	}
}

func main() {
	// create()
	// cancel()
	// reset()
	// dur()
	after()
}