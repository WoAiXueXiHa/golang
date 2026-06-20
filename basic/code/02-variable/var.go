package main

import "fmt"

// 只声明不初始化，会给默认值
func testZero() {
	var a int
	var b bool
	var c float64
	var d string
	fmt.Println(a, b, c, "-", d, "-")
}

// 几种声明方式
func testVar() {
	// 最常规 var 变量名 类型 = 值
	var num int = 20

	// := 赋值
	str := "golang"

	// var 变量名 = 值    类似 C++ 的 auto
	var val = 2.10

	fmt.Printf("typeof num is %T, value is %d", num, num)
	fmt.Printf("typeof str is %T, value is %s", str, str)
	fmt.Printf("typeof val is %T, value is %f", val, val)

	// 多变量声明
	var v1, v2 = 100, "ll"
	n1, n2 := 1.0, 200

	fmt.Println(v1, v2, n1, n2)
}

// 交换两数的值
func exangeVal(val1, val2 int) {
	fmt.Printf("交换前：val1 = %d, val2 = %d\n", val1, val2)
	val1, val2 = val2, val1
	fmt.Printf("交换后：val1 = %d, val2 = %d\n", val1, val2)
}

func main() {
	testZero()
	fmt.Print("\n-------------------------\n")
	testVar()
	fmt.Print("\n-------------------------\n")
	exangeVal(100, 0)
}
