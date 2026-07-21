package main

import (
	"context"
	"fmt"
	"time"
)

// context 是在函数调用链、goroutine 之间传递取消信号、超时时间、请求数据的工具
// 控制超时、主动取消、通知多个goroutine停止、传递请求级信息

// type Context interface {
// 	Deadline() (deadline time.Time, ok bool)	// 截止时间
// 	Done() <-chan struct{}						// 取消信号
// 	Err() error 								// 取消原因
// 	Value(key any) any 							// 取值
// }

// 1. context 如何创建？
// ctx := context.Background() 根context，在main、测试、请求入口创建
// ctx := context.TODO() 暂时不知道用什么context时占位
// ctx, cancel := context.WithCancel(context.Background())
// defer cancel()   23 24行，手动取消
// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// defer cancel()   2s后自动取消
// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
// defer cancel()   指定截止时间取消

// 2. 使用场景：业务超时
// 假设一个系统里用户提交一个请求，需要调用AI，但最多等两秒
func demo1() {
	callAI := func(ctx context.Context, question string) (string, error) {
		select {
		case <-time.After(3 * time.Second):
			return "AI answer: " + question, nil
		// <-ctx.Done() 等 context 被取消、超时、截止，一旦触发，Done()这个channel会关闭，所有监听它的goroutine都能收到信号
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	answer, err := callAI(ctx, "Simply introducing the program of Go")
	if err != nil {
		fmt.Println("Failed to request: ", err)
		return
	}

	fmt.Println(answer)
}

// 主动取消 goroutine
// 比如启动一个后台任务，不想让他一直跑
func demo2() {
	worker := func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Worker stopped: ", ctx.Err())
				return

			default:
				fmt.Println("Worker is running...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	go worker(ctx)
	time.Sleep(2 * time.Second)
	cancel()

	time.Sleep(time.Second)
}

// WithValue 使用
type ctxKey string

const userIDKey ctxKey = "user_id"

func handle(ctx context.Context) {
	userID := ctx.Value(userIDKey).(int64)
	fmt.Println("The user: ", userID)
}
func main() {
	demo1()
	demo2()

	ctx := context.WithValue(context.Background(), userIDKey, int64(1000))
	handle(ctx)
}
