package main

import (
	"fmt"
	"time"
)

func worker(done <-chan struct{}) {
	for {
		select {
		case <-done: // 每一轮都要检查，有没有人通知我停下
			fmt.Println("woker: receive the canceling signal, exit")
			return // 真正让 goroutine 退出的地方
		default:
			fmt.Println("worker: working")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	done := make(chan struct{}) // 创建一个停止信号通道

	go worker(done) // 把停止信号交给worker

	time.Sleep(2 * time.Second)

	fmt.Println("main: inform worker to stop")
	close(done) // 广播停止信号

	time.Sleep(time.Second)
	fmt.Println("main: done~")
}
