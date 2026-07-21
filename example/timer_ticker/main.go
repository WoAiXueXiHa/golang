package main

import (
	"fmt"
	"time"
)

// Timer 是到点执行一次

// 1. 创建 Timer
func demo1() {
	timer := time.NewTimer(2 * time.Second)

	fmt.Println("Begin to wait")
	<-timer.C // 是一个 channel，到时间后，会往里面发一个时间值
	fmt.Println("Time over, running")
}

// 2. 停止 Timer
// 场景：订单 5s 未支付就取消，但用户提前支付了，就要停止计时器
func demo2() {
	timer := time.NewTimer(5 * time.Second)

	go func() {
		<-timer.C
		fmt.Println("The order has exceeded and will be automatically cancelled")
	}()

	time.Sleep(2 * time.Second)

	stopped := timer.Stop() // 如果 timer 还没触发，就取消它
	if stopped {
		fmt.Println("User has payed. Canceled the timeout timer")
	}
	time.Sleep(5 * time.Second)
}

// 3. 重置 timer，重新计时
// 场景：用户每次操作都把“自动退出时间”延后
func demo3() {
	timer := time.NewTimer(3 * time.Second)

	go func() {
		<-timer.C
		fmt.Println("User has not operated for too long and will automatically exit")
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("User has operated once.Timer was reset.")

	timer.Reset(3 * time.Second)

	time.Sleep(5 * time.Second)
}

// 4. time.After 使用 一次性定时器
func demo4() {
	queryAI := func() chan string {
		ch := make(chan string)

		go func() {
			time.Sleep(3 * time.Second)
			ch <- "AI returns the result."
		}()
		return ch
	}

	resultCh := queryAI()
	select {
	case result := <-resultCh:
		fmt.Println(result)
	case <-time.After(2 * time.Second):
		fmt.Println("AI was timeout and returned the catch-all result.")
	}
}

// 5. time.AfterFunc 使用，到时间后自动执行一个函数
func demo5() {
	timer := time.AfterFunc(2*time.Second, func() {
		fmt.Println("Order sending timeout reminder.")
	})

	time.Sleep(time.Second)

	stopped := timer.Stop()
	if stopped {
		fmt.Println("Order has payed. Canceled reminder")
	}

	time.Sleep(3 * time.Second)
}

// ticker 是周期性定时器
// 场景：每疫苗扫描一次待处理惹怒我
func demo6() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	done := time.After(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println("Scan outbox and deliver pending tasks")

		case <-done:
			fmt.Println("Stopped to scan.")
			return
		}
	}
}

func main() {
	demo1()
	demo2()
	demo3()
	demo4()
	demo5()
	demo6()

}
