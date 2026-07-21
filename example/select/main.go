package main

import (
	"fmt"
	"time"
)

// 1. 等任务结果或超时
// func callAI() <-chan string {
// 	resultCh := make(chan string)

// 	go func() {
// 		time.Sleep(2 * time.Second)
// 		resultCh <- "AI result"
// 	}()
// 	return resultCh
// }

// func main() {
// 	resultCh := callAI()
// 	// select 同时等待 resultCh 和 time.After
// 	// 谁先有动静，就执行谁
// 	select {
// 	case result := <-resultCh:
// 		fmt.Println("succeed:", result)
// 	case <-time.After(3 * time.Second):
// 		fmt.Println("timeout")
// 	}
// }

// 2. 等任务结果 or 用户取消
func demo1() {
	resultCh := make(chan string)
	cancelCh := make(chan struct{})

	go func() {
		time.Sleep(2 * time.Second)
		resultCh <- "AI task was completed"
	}()

	go func() {
		time.Sleep(time.Second)
		close(cancelCh)
	}()

	select {
	case result := <-resultCh:
		fmt.Println("Writing in DB:", result)
	case <-cancelCh:
		fmt.Println("User canceled task, aborted")
	}
}

// worker循环：有任务就处理，收到退出信号就停
type Task struct {
	ID      int
	Content string
}

func demo2() {
	worker := func(taskCh <-chan Task, stopCh <-chan struct{}) {
		for {
			select {
			case task := <-taskCh:
				fmt.Println("Carry out task:", task.ID, task.Content)

			case <-stopCh:
				fmt.Println("Worker exit")
				return
			}
		}
	}

	taskCh := make(chan Task)
	stopCh := make(chan struct{})

	go worker(taskCh, stopCh)

	taskCh <- Task{ID: 1, Content: "learn goroutine"}
	taskCh <- Task{ID: 2, Content: "learn select"}

	close(stopCh)
}

// 4. 非阻塞投递任务，队列满了就丢弃/返回失败
func demo3() {
	type Task struct {
		ID int
	}

	queue := make(chan Task, 2)

	queue <- Task{ID: 1}
	queue <- Task{ID: 2}

	select {
	case queue <- Task{ID: 3}:
		fmt.Println("Task was dispensed successfully")
	default:
		fmt.Println("Queue is full, failed to dispense")
	}
}

// 5. 非阻塞查询结果：没有就先返回处理中
func demo4() {
	resultCh := make(chan string)

	select {
	case result := <-resultCh:
		fmt.Println("Result of task: ", result)
	default:
		fmt.Println("Task is being processed")
	}
}
func main() {
	demo1()
	demo2()
	demo3()
	demo4()
}

// select 没有 default，会阻塞
