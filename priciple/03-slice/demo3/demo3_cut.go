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

func case1(s []int) {
	s = s[1:]
	PrintSlice(&s)
}

func case2(s []int) {
	s = s[1:3]
	PrintSlice(&s)
}

func case3(s []int) {
	s = s[len(s)-1:]
	PrintSlice(&s)
}

func case4(s []int) {
	s1 := s[2:]
	PrintSlice(&s1)
}

func main() {
	s := make([]int, 5)

	case1(s)
	case2(s)
	case3(s)
	case4(s)

	PrintSlice(&s)
}
