// ===========================================================
// 0. 前置知识回顾
// 程序：磁盘上存储的静态文件
// 进程：程序加载到内存中，是程序的一个实体，资源的容器，进行资源分配
// 线程：内核的执行流，CPU 调度的实体，被执行的上下文
// 协程：用户的执行流，代码/运行时调度，拥有极小的栈空间（2KB）
// 
// 并发：一个 CPU 上，不同线程轮换执行，宏观上感觉是同时发生
// 并行：多个 CPU 上，不同线程同时执行
//
// ===========================================================

package main
 
import (
	"fmt"
	"sync"
	"time"
	"runtime"
) 

// ===========================================================
// 1. goroutine 基本使用
// ===========================================================
func demo1_basic() {
	fmt.Println("=== 1. goroutine 基本使用 ===")
	go fmt.Println("Hello from goroutine")

	// 若不等待，主协程退出时，这个协程可能还没执行
	time.Sleep(10 * time.Millisecond)
	fmt.Println("主协程结束")
}

// ===========================================================
// 2. 主协程：Go 给 main() 创建的 协程
// 主协程退出 -> 所有子协程立即被杀死，不等待
// ===========================================================
func demo2_main_goroutine() {
	fmt.Println("\n=== 2. 主协程退出，子协程被杀死 ===")
	
	go func() {
		time.Sleep(500 * time.Millisecond)

		fmt.Println("这句话大概率不会被打印....")
	}()

	fmt.Println("主协程结束 -> 子协程被杀死...")
}

// ===========================================================
// 3. WaitGroup：等待一组协程完成
// 像一个计数器：
// Add(n)	计数器 + n，登记 n 个任务
// Done()	计数器 - 1，完成一个任务
// Wait()	阻塞等待到计数器归零
// ===========================================================
func worker(name string, wg *sync.WaitGroup) {
	defer wg.Done() 
	for i := 0; i < 3; i++ {
		fmt.Printf("[%s] 步骤 %d\n", name, i + 1)
		time.Sleep(50 * time.Millisecond)
	}
}

func demo3_WaitGroup() {
	fmt.Println("\n=== 3. WaitGroup 多协程同步 ===")
	var wg sync.WaitGroup

	wg.Add(2)	// 明确创建两个协程
	// 必须要取地址，否则两个协程操作的是副本
	go worker("worker A", &wg)
	go worker("worker B", &wg)

	wg.Wait()
	fmt.Println("All works were done!")

}

// ===========================================================
// 4. 闭包变量捕获陷阱
// ===========================================================
func demo4_closure() {
	fmt.Println("\n === 4. 闭包变量捕获陷阱 ===")
	
	// 错误写法❌
	fmt.Println("\n--- ❌闭包直接捕获循环变量❌")
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1) 
		go func() {
			defer wg.Done()
			// ❌ i 是循环变量，所有协程共享同一个变量地址
			fmt.Printf("协程 [%d]\n", i)
		}()
	}
	wg.Wait()

	// 正确写法✔
	fmt.Println("\n--- 1. ✔循环内创建局部变量✔ ---")
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		i := i 		// 遮蔽外层的循环变量，用其它变量名也可以
		go func() {
			defer wg.Done()
			fmt.Printf("协程 [%d]\n", i)
		}()
	}
	wg.Wait()

	fmt.Println("\n--- 2. ✔通过参数传入，强制拷贝✔ ---")
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		i := i 		// 遮蔽外层的循环变量，用其它变量名也可以
		go func(id int) {
			defer wg.Done()
			fmt.Printf("协程 [%d]\n", id)
		}(i)
	}
	wg.Wait()

}

// ===========================================================
// 5. 观察 goroutine 交错执行
// goroutine 调度是非确定性的，每次运行的打印顺序可能不同
// 使用 Sleep 会让出执行权，更容易观察
// ===========================================================
func demo5_interleaving() {
	fmt.Println("\n === 5. goroutine 交错执行 ===")
	var wg sync.WaitGroup

	printer := func(label string, n int) {
		defer wg.Done() 
		for i:= 1; i <= n; i++ {
			fmt.Printf("%s%d ", label, i)
			time.Sleep(time.Millisecond)	// 让出时间片，给其它 goroutine
		}
	}

	fmt.Println("每次执行顺序可能不同👇")
	wg.Add(3) 
	go printer("a", 5)
	go printer("b", 5)
	go printer("c", 5)
	wg.Wait()
	fmt.Println()

	// 对比：去掉 Sleep 试试
	fmt.Println("\n去掉 Sleep：一个 goroutine 可能连续跑完整个循环👇")
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 8; i++ {
			fmt.Printf("X无Sleep->%d,", i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i <= 8; i++ {
			fmt.Printf("Y无Sleep%d,", i)
		}
	}()
	wg.Wait()
	fmt.Println()

}

// ===========================================================
// 6. 创建上万个 goroutine
// 每个 goroutine 初始栈约 2KB，OS 线程约 1MB
func demo6_light_weight() {
	fmt.Println("\n=== 6. goroutine 轻量级特性 ===")
	fmt.Printf("启动前 goroutne 数: %d\n", runtime.NumGoroutine())

	const N = 10000
	var wg sync.WaitGroup
	wg.Add(N)

	start := time.Now() 
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
		}()
	}

	fmt.Printf("创建 %d 个 goroutine 后：%d 个\n", N, runtime.NumGoroutine())
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("全部完成耗时：%v（如果串行需要 %v）\n", elapsed, N*10*time.Millisecond)
}

func main() {
	// demo1_basic()
	// demo2_main_goroutine()
	// demo3_WaitGroup()
	// demo4_closure()
	// demo5_interleaving()
	demo6_light_weight()
}