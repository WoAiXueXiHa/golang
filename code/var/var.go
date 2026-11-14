package main

import (
	"fmt"
	"unsafe"
)

type User struct {
	ID     int
	Name   string
	Active bool
}

func main() {
	fmt.Println("=== nil silce VS 空 slice ===")
	var s1 []int  // nil slice 零值
	s2 := []int{} // 空 slice，有底层结构

	fmt.Printf("s1: len=%d, cap=%d, s1==nil: %v, ptr=%p\n",
		len(s1), cap(s1), s1 == nil, unsafe.SliceData(s1))

	fmt.Printf("s2: len=%d, cap=%d, s2==nil: %v, ptr=%p\n",
		len(s2), cap(s2), s2 == nil, unsafe.SliceData(s2))

	fmt.Println("=== new(map[string]int) 的陷阱 ===")

	mPtr := new(map[string]int) // 返回*map[string]int，没有底层数据结构
	fmt.Printf("mPtr 本身 (指针地址): %p, mPtr == nil: %v\n", mPtr, mPtr == nil)
	fmt.Printf("*mPtr 解引用后的值: %v, *mPtr == nil: %v\n", *mPtr, *mPtr == nil)

	// 如果此时直接执行 (*mPtr)["a"] = 1，会直接触发 panic: assignment to entry in nil map
	// 必须先用 make 初始化底层哈希表：
	*mPtr = make(map[string]int)
	(*mPtr)["a"] = 1
	fmt.Printf("make 初始化后写入成功: %v\n", *mPtr)

	fmt.Println("=== if 作用域内的 := 变量遮蔽 ===")
	x := 10
	fmt.Printf("外层初始: 值=%d, 内存地址=%p\n", x, &x)

	if true {
		x := 20
		fmt.Printf("if 内部:  值=%d, 内存地址=%p  <-- 地址变了！说明是新变量\n", x, &x)
	}
	fmt.Printf("if 结束后: 值=%d, 内存地址=%p  <-- 外层变量完全没受影响\n", x, &x)
}

// func main() {
// 	var u User
// 	fmt.Printf("%+v\n", u)

// 	var s []int

// 	fmt.Println(s == nil) // true

// 	s = append(s, 1)
// 	fmt.Println(s) // [1]

// 	var m map[string]int

// 	fmt.Println(m == nil) // true
// 	fmt.Println(m["a"])   // 0

// 	// m["a"] = 10 // panic!

// 	m = make(map[string]int)
// 	m["a"] = 10

// 	// var ch chan int
// 	// ch <- 10 // 永久阻塞，除非所在 goroutine 被其他机制终止
// }
