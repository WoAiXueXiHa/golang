package main

import (
	"fmt"
	"context"
	"time"
)

// =================================================
// 1. context.WithCancel 使用
// 取消控制函数
// =================================================
func demo1_WithCancel() {
	Watch := func(ctx context.Context, name string) {
		for {
			select {
			case <-ctx.Done():
				// 主 goroutine 调用 cancel 后，会发送一个信号到 ctx.Done()这个channel
				fmt.Printf("%s exit!\n", name)	
				return 
			default:
				fmt.Printf("%s watching...\n", name)
				time.Sleep(time.Second)
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	go Watch(ctx, "goroutine1")
	go Watch(ctx, "goroutine2")

	time.Sleep(6 * time.Second)

	fmt.Println("end watching!!!")
	cancel()
	time.Sleep(time.Second)
}

// =================================================
// 2. context.WithDeadline 和 context.WithTimeout使用
// 取消控制函数
// =================================================
func demo2_WithDeadline() {
	Watch := func(ctx context.Context, name string) {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("%s exit!\n", name)
				return 
			default:
				fmt.Printf("%s watching...\n", name)
				time.Sleep(time.Second)
			}
		}
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
	defer cancel()
	go Watch(ctx, "goroutine1")
	go Watch(ctx, "goroutine2")

	time.Sleep(6 * time.Second)
	fmt.Println("end watching!!!")
}

func demo3_WithTimeout() {
	Watch := func(ctx context.Context, name string) {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("%s exit!\n", name)
				return
			default:
				fmt.Printf("%s watching...\n", name)
				time.Sleep(time.Second)
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()
	go Watch(ctx, "goroutine1")
	go Watch(ctx, "goroutine2")

	time.Sleep(6 * time.Second)
	fmt.Println("end watching!!!")
}

// =================================================
// 3. context.WithValue 使用
// =================================================
func demo4() {
	func1 := func(ctx context.Context) {
		fmt.Printf("name is %s\n", ctx.Value("name").(string))
	}

	ctx := context.WithValue(context.Background(), "name", "vect")
	go func1(ctx)
	time.Sleep(time.Second)
}


func main() {
	// demo1_WithCancel()
	// demo2_WithDeadline()
	// demo3_WithTimeout()
	demo4()
}

