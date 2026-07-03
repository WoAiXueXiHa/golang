package main

import (
	"fmt"
	"time"
)

// type Ticker struct {
//		C <-chan Time
// 		r runtineTimer
// }
// 每隔时间段d就向通道发送当时的时间，根据这个管道消息来触发事件
// ticker 只要定义完成，就从当前事件开始计时，每隔固定时间都会触发
// 只有关闭Ticker对象才不会继续发送时间消息

func demo1() {
	Watch := func() chan struct{} {
		ticker := time.NewTicker(time.Second)

		ch := make(chan struct{})

		go func(ticker *time.Ticker) {
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					fmt.Println("watch!!!")
				case <-ch:
					fmt.Println("Ticker Stop!!!")
					return
				}
			}
		}(ticker)
		return ch
	}

	ch := Watch()
	time.Sleep(5 * time.Second)
	ch <- struct{}{}
	close(ch)
}

// Ticker Stop!!! 打印不稳定
// func main() {
// 	demo1()
// }
