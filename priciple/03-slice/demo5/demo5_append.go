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

func case1() {
	s1 := make([]int, 3, 3)
	s1 = append(s1, 1)

	PrintSlice(&s1)
}

func case2() {
	s1 := make([]int, 3, 4)
	s2 := append(s1, 1)

	PrintSlice(&s1)
	PrintSlice(&s2)
}

func case3() {
	s1 := make([]int, 3, 3)
	s2 := append(s1, 1)

	PrintSlice(&s1)
	PrintSlice(&s2)
}

func main() {
	case1()
	case2()
	case3()
}
