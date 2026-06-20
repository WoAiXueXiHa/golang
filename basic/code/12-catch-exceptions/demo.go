package main

import "fmt"

// 2. panic 传递
func testPanic1() {
	fmt.Println("testPanic1上半部分")
	testPanic2()
	fmt.Println("testPanic1下半部分")
}

func testPanic2() {
	defer func() {
		recover()
	}()

	fmt.Println("testPanic2上半部分")
	testPanic3()
	fmt.Println("testPanic2下半部分")
}

func testPanic3() {
	fmt.Println("testPanic3上半部分")
	panic("在testPanic3出现了panic")
	fmt.Println("testPanic3下半部分")
}

func main() {

	// // 1. recover 捕获异常
	// defer func() {
	// 	if error := recover(); error != nil {
	// 		fmt.Println("出现了panic，使用recover获取信息：", error)
	// 	}
	// }()

	// fmt.Println("111111111111")
	// panic("出现panic")
	// fmt.Println("222222222222")

	// 2. panic 传递
	fmt.Println("程序开始")
	testPanic1()
	fmt.Println("程序结束")

}
