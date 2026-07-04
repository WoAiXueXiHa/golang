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
	s1 := []int{1, 2, 3}
	s2 := make([]int, len(s1))

	copy(s2, s1)

	PrintSlice(&s1)
	PrintSlice(&s2)

}
