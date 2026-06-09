package main

import (
	"fmt"
	"math"
)

// func 函数名(参数列表) 返回类型{}

func swap(x, y string) (string, string) {
	return y, x
}

// 参数传递：go 的函数参数传递都是值传递，没有引用传递
// 传值，拷贝原来的值，传指针，拷贝指针的内容
// 拷贝的指针的内容指向原内容，故可以顺着指针进行修改
func add(a int) {
	a++
	fmt.Printf("add() 中的 a 值为 %d\n", a)
}

// 函数变量作为回调函数
// 给函数起别名，声明一个函数类型
type fc func(int) int

// 回调函数：把一个函数交给另外一个函数，由另外一个函数在合适的时候调用它
func callBack(x int) int {
	fmt.Printf("我是回调，x: %d\n", x)
	return x
}
func CallBack(x int, f fc) {
	// 传进来的是callback函数，函数执行需要传入一个ingt类型参数，所以传入x
	f(x)
}

// 闭包：匿名函数，是一个内联的语句或者表达式
func getNumber() func() int {
	i := 0
	return func() int {
		i += 1
		return i
	}
}
func main() {
	a, b := swap("hello", "golang")
	fmt.Println(a, b)
	fmt.Println("-------------------")
	num := 100
	fmt.Printf("调用add前：num=%d\n", num)
	add(num)
	fmt.Printf("调用add后： num=%d\n", num)

	fmt.Println("-------------------")

	// 函数变量，go 中一切皆变量
	// 声明函数变量
	getSquareRoot := func(x float64) float64 {
		return math.Sqrt(x)
	}
	fmt.Println(getSquareRoot(16))

	// 执行回调函数
	fmt.Println("-------------------")
	CallBack(1, callBack)
	fmt.Println("-------------------")

	nextNumber := getNumber()

	// 调用 nextNumber 函数，变量自增1并返回
	fmt.Println(nextNumber())
	fmt.Println(nextNumber())
	fmt.Println(nextNumber())

	// 创建新的函数 nextNumber1，查看结果
	nextNumber1 := getNumber()
	fmt.Println(nextNumber1())
	fmt.Println(nextNumber1())
}
