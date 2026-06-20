package main

import (
	"fmt"
	"sync"
	"time"
)

// 0、前置知识回顾
// 程序：磁盘上存储的静态文件
// 进程：程序加载到内存中，程序的一个实体，资源的容器，进行资源分配
// 线程：内核的执行流，CPU 调度实体，被执行的上下文
// 协程：用户的执行流，代码/运行时调度，极小的用户栈

// 并发：一个 CPU 上，不同线程轮换执行，宏观上感觉是同时发生，微观上并非
// 并行：多个 CPU 上，不同线程同时执行

// 1. 主协程，main()函数的是默认的主协程，后续创建的都是子协程
// func myGoroutine() {
// 	fmt.Println("myGoroutine")
// }

// func main() {
// 	// 只输出了主协程的内容
// 	// 主协程打印完之后就退出了，不会管子协程
// 	go myGoroutine()
// 	fmt.Println("end--------")
// 	// 执行完休眠两秒，这个两秒是猜测，具体不清楚
// 	time.Sleep(2 * time.Second)
// }

// // 2. 多协程调用和等待

// func worker(id int, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	fmt.Printf("worker[%d] 开始工作。。。\n", id)
// 	time.Sleep(time.Second * 1)
// 	fmt.Printf("worker[%d]结束工作。。。\n", id)
// }

// func main() {
// 	fmt.Println("----------- 主协程开始 --------------")
// 	var wg sync.WaitGroup
// 	// 另外一种写法，明确知道了要用几个协程
// 	// wg.Add(3)

// 	for i := 1; i <= 3; i++ {
// 		// 必须写 1，因为固定了 3 个协程， 如果写 2，就意味着 2 * 3 = 6 个协程
// 		wg.Add(1)

// 		go worker(i, &wg)
// 	}

// 	fmt.Println("------------ 主协程开始等待 ----------")

// 	wg.Wait()

// 	fmt.Println("-------- 主协程退出 -------------")
// }

// 3. 多协程异常捕获和 recover 捕获的三重防线
func badCapture() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("试图捕获子协程异常: ", err)
		}
	}()

	go func() {
		fmt.Println("子协程：我要 panic 了")
		panic("子协程 panic ")
	}()

	time.Sleep(time.Second * 1)
}

func testCaptureConut() {
	defer func() {
		// recover 捕获次数判定
		// 一次函数退栈过程中，内置的 recover() 只要成功执行并返回了 nil 的值
		// 当前发生的 panic 就会立刻被抹去
		if err1 := recover(); err1 != nil {
			fmt.Printf("第一次执行 recover 捕获成果: %v", err1)
		}

		if err2 := recover(); err2 != nil {
			fmt.Printf("第二次执行 recover， 成功捕获: %v", err2)
		} else {
			fmt.Println("第二次执行 recover， 返回 nil， 证明一个 panic 只能 reover 一次!")
		}
	}()

	panic("触发")
}

func SafeGo(wg *sync.WaitGroup, task func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 捕获范围闭环，让每个子协程的调用栈都安全
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("拦截成功，检测到子协程 panic， 已经成功隔离，原因为 %v\n", err)
			}
		}()

		// 真正执行业务逻辑
		task()
	}()
}

func main() {
	fmt.Println("------------ 多协程异常捕获和捕获次数测试 ---------------")

	testCaptureConut()

	time.Sleep(time.Second * 1)
	fmt.Println("-------------- 并发防止崩溃测试 ---------------")

	var wg sync.WaitGroup

	SafeGo(&wg, func() {
		fmt.Println("任务一：我是健康的业务，正常执行")
	})

	SafeGo(&wg, func() {
		fmt.Println("任务二：我很危险，我要踩雷了")
		a := 1
		b := 0
		fmt.Println(a / b)
	})

	wg.Wait()
	fmt.Println("虽然任务二崩溃了，但是因为 SafeGo 的隔离，主进程顺利完成！")

	// badCapture()		// 解除此行注释就能看到panic
}
