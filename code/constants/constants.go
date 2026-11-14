package main

import "fmt"

const (
	a = 10
	b = 3.14
	c = 'a'
)

func main() {
	x := a
	y := b
	z := c

	fmt.Printf("%T\n", x)
	fmt.Printf("%T\n", y)
	fmt.Printf("%T\n", z)

}
