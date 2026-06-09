package main

import (
	"fmt"
	"unsafe"
)

func show() {
	const a, b = "val", 200
	fmt.Println(a, b)
}

// 常量用于枚举
const (
	Unkonw  = 0
	Success = 1
	Fail    = 2
)

// 常量可以用len() cap() unsafe.Sizeof() 函数计算表达式的值
func testCal() {
	const (
		a = "123"
		b = 10
		c = unsafe.Sizeof(a)
	)
	// typedef struct {
	// 		char* buffer
	//		size_t len
	// } string      理解成这个结构体
	fmt.Printf("a = %s, b = %d, c = %d\n", a, b, c)
}

// iota
// 第一行是0 每行递增1
// 如果有表达式，在表达式中参与计算，并保持递增性质
// 如果此行没有表达式，继承上面最近有表达式的行
func testIota() {
	const (
		a = iota
		b
		c
	)
	fmt.Println(a, b, c)

	const (
		val1, val2 = iota + 1, iota + 2 // 1 2
		val3, val4                      // 2 3

		val5, val6 = iota + 10, iota * 10 // 12 20
		val7, val8                        // 13 30
	)
	fmt.Println(val1, val2, val3, val4, val5, val6, val7, val8)
}

func main() {
	testCal()
	fmt.Println("-------------------------")
	testIota()
	fmt.Println("-------------------------")
	show()
}
