package main

import (
	"fmt"
	"math"
)

// 函数基础用法
// func function_name (参数列表) 返回类型 {
//	...
// }

// 1. 值传递和指针传递，go 中只有值传递，指针传递也是拷贝了指针的副本，有指针的副本也能找到具体的值
func swapByValue(x,y int) {
	x, y = y, x
}

func swapByPtr(x, y *int) {
	*x, *y = *y, *x
}

// 3. 函数可以作为另外一个函数的实参传递
// 声明一个函数类型
// func(int) int 这个匿名函数现在叫 fc，起了个别名
type fc func(int) int

// 接收一个int类型和一个func(int)int类型参数
// 把外面传进来的x交给f执行
func callBack(x int, f fc) {
	res := f(x)
	fmt.Println("回调返回值：", res)
}

// 完全符合 func(int)int 的函数类型，可以直接作为 fc 类型作为变量传给 callBack
func cb(x int) int {
	fmt.Printf("我是回调，x: %d\n", x)
	return x
}

// 4. 闭包：匿名函数 + 捕获自己作用域外的变量
func getNumber() func() int {
	i := 0
	return func() int {
		i += 10
		return i
	}
}

func main() {
	// 1. 值传递和指针传递
	a, b := 10, 20
	fmt.Println("调用前 a =", a, "b =", b)
	fmt.Println("------------- 值传递 ---------------")
	swapByValue(a, b)
	fmt.Println("值传递调用结束，main中 a =", a, "b =", b)
	fmt.Println("------------- 指针传递 ---------------")
	swapByPtr(&a, &b)
	fmt.Println("指针传递调用结束，main中 a =", a, "b =", b)

	// 2. 函数作为变量
	// func(x float64) float64 { return mat.Sqrt(x) }
	// 把整个函数作为一个变量赋值给 getSquareRoot
	// getSquareRoot 存的是一个函数类型的值，类型是 func(x float64) float64
	// 这个函数没有名字，就是匿名函数
	fmt.Println("------------- 函数作为变量 ---------------")
	getSquareRoot := func(x float64) float64 {
		return math.Sqrt(x)
	}
	fmt.Println(getSquareRoot(16))

	// 3. 函数可以作为另外一个函数的实参传递
	fmt.Println("------------- 函数作为另一个函数的实参 ---------------")
	callBack(99, cb)

	// 4. 闭包
	fmt.Println("------------- 闭包 ---------------")
	nxetNumber := getNumber()
	fmt.Println(nxetNumber())
	fmt.Println(nxetNumber())
	fmt.Println(nxetNumber())	
}