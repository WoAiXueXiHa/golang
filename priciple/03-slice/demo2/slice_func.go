package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func PrintSlice(s *[]int) {
	// Go 是强类型，正常情况 *[]int 绝对不能转成 *reflect.SliceHeader
	// 而 unsafe.Pointer 类似 void*，可以接收任意类型指针
	// 对于 reflect.SliceHeader
	// type SliceHeader struct {
	// 	Data uintptr  // 对应底层数组地址，这个不是指针，就是存了地址数字的类型而已
	// 	Len  int      // 对应长度
	// 	Cap  int      // 对应容量
	// }

	ss := (*reflect.SliceHeader)(unsafe.Pointer(s))

	fmt.Printf("slice struct: %+v, slice is %v\n", ss, s)
}

func test(s []int) {
	PrintSlice((&s))
}

// 底层数组不变
func demo3_case1(s []int) {
	s[1] = 1000
	PrintSlice(&s)
}

// 底层数组变化
func demo3_case2(s []int) {
	s = append(s, 1000)
	s[1] = 1000
	PrintSlice(&s)
}
func demo3_infunc_modify() {
	s := make([]int, 5)
	demo3_case1(s)
	demo3_case2(s)
	PrintSlice(&s)
}

func demo2_slice_func() {
	s := make([]int, 5, 10)
	PrintSlice(&s)
	test(s)
}

func main() {
	demo2_slice_func()
	demo3_infunc_modify()
}
