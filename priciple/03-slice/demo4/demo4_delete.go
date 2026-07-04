package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func PrintSlice(s *[]int) {
	ss := (*reflect.SliceHeader)(unsafe.Pointer(s))

	fmt.Printf("slice struct: %+v, slice is %v\n", ss, s)
}

func main() {
	s := []int{0, 1, 2, 3, 4}

	_ = s[4]
	PrintSlice(&s)
	// 删除第一个元素，从0开始计数
	// [0,1) + [2, len(s))
	s1 := append(s[:1], s[2:]...)
	{
		// 拷贝元素
		// 0, 1, 2, 3, 4
		// 0, 2, 3, 4, 4
	}

	PrintSlice(&s1)
	PrintSlice(&s)

	// 访问原切片
	_ = s[4]
	// 访问从原切片中删除了一个元素的切片
	_ = s1[4]

}
