// 两个协程交替打印字母和数字

// 两个channel做同步，数字协程收到信号打印数字，然后通知字母协程，字母协程打印完再通知数字，形成循环

package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	num := make(chan struct{})
	letter := make(chan struct{})

	go func() {
		defer wg.Done()
		for i := 1; i <= 26; i++ {
			<-num
			fmt.Println(i)
			letter <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for c := 'A'; c <= 'Z'; c++ {
			<-letter

			fmt.Printf("%c\n", c)
			if c != 'Z' {
				num <- struct{}{}
			}
		}
	}()

	num <- struct{}{}

	wg.Wait()

}
