package main

import "fmt"

// 将数组传递到函数中，数组的地址不一样
func test(arr [3]int) {
	fmt.Printf("arr 内: %p\n", &arr)
}
func f1() {
	arr := [3]int{1, 2, 3}
	test(arr)
	fmt.Printf("arr 外: %p\n", &arr)
}

// 拷贝数组，修改旧数组，对新数组无影响
func f2() {
	arr1 := [3]int{1, 2, 3}
	arr2 := arr1

	arr1[0] = 100
	fmt.Println(arr1)
	fmt.Println(arr2)
}

func main() {
	f1()
	fmt.Printf("---------------------------\n")
	f2()
}
