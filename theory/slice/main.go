package main

import (
	"fmt"
	"slices"
)

func appendInFunc(s []int) {
	s = append(s, 99)
}

func appendReturn(s []int) []int {
	return append(s, 99)
}

func main() {
	a := []int{1, 2, 3, 4}
	b := a[:2]
	c := append(b, 100)

	fmt.Println("case1 a:", a)
	fmt.Println("case1 b:", b)
	fmt.Println("case1 c:", c)

	x := []int{1, 2, 3}
	appendInFunc(x)
	fmt.Println("case2 x:", x)

	x = appendReturn(x)
	fmt.Println("case3 x:", x)

	y := []int{1, 2, 3, 4}
	z := y[:2:2]
	z = append(z, 100)
	fmt.Println("case4 y:", y)
	fmt.Println("case4 z:", z)

	p := []*int{}
	v1, v2, v3 := 1, 2, 3
	p = append(p, &v1, &v2, &v3)
	p = slices.Delete(p, 1, 2)
	fmt.Println("case5 len/cap:", len(p), cap(p))

	m := [][]int{{1, 2}, {3, 4}}
	n := make([][]int, len(m))
	for i := range m {
		n[i] = slices.Clone(m[i])
	}
	n[0][0] = 100
	fmt.Println("case6 deep copy:", m[0][0], n[0][0])
}
