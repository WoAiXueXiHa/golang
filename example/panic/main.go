package main

import (
	"fmt"
)

// 抛异常：panic()
// 异常捕获：defer + recover()
// 异常传递：panic 沿着调用栈向上传递，并执行沿途的 defer

// 1. 使用 recover 捕获 panic
// recover() 只有在 defer 延迟函数里调用，才能捕获当前 goroutine 正在传播的 panic
func demo1() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Catched panic: ", err)
		}
	}()

	fmt.Println("Begin...")
	fmt.Println("The exe is error...")
	fmt.Println("This line won't be executed")
}

// 2. panic 传递方式
// 假设调用链路：main->a->b->c
func a() {
	defer fmt.Println("a defer")
	b()
}

// 3. 在中间层 recover，panic 就会停止传递
func b() {
	//defer fmt.Println("b defer")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("b catch panic:", err)
		}
	}()
	c()
}

func c() {
	defer fmt.Println("c defer")
	panic("crash")
}

// recover 只能捕获同一个 goroutine 里的 panic
// 所以每个 goroutine 自己内部要 recover
func demo2() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("aaa catch: ", err)
		}
	}()

	go func() {
		panic("bbb crash")
	}()

	select {}
}
func main() {
	//demo1()
	// a()
	// fmt.Println("main end")
	demo2()
}
