package main

import "fmt"

func inSeq() func() int {
	// 变量 i 的生命周期延长
	i := 0
	return func() int {
		i++
		return i
	}
}

func main() {
	nextInt := inSeq()

	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt())

	nextInts := inSeq()
	fmt.Println(nextInts())
}
