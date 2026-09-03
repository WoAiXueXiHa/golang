// 三个协程，每秒打印 cat dog fish，顺序不能变化，协程1打印cat，协程2打印dog，协程3打印fish
// 思路：用三个固定的 channel 来接收信号
// 先给cat发信号，cat打印完给dog发信号，dog打印完给fish发信号，fish打印完给cat，形成循环
package main

import (
	"fmt"
	"time"
)

func main() {
	cat := make(chan struct{})
	dog := make(chan struct{})
	fish := make(chan struct{})

	go func() {
		for {
			<-cat
			fmt.Println("cat")
			dog <- struct{}{}
		}
	}()

	go func() {
		for {
			<-dog
			fmt.Println("dog")
			fish <- struct{}{}
		}
	}()

	go func() {
		for {
			<-fish
			fmt.Println("fish")
			time.Sleep(time.Second)
			cat <- struct{}{}
		}
	}()

	cat <- struct{}{}

	select {}
}
