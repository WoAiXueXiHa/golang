package main

import "fmt"

func main() {
	var x float64 = 3.9
	var y int = int(x)

	var a float64 = -3.9
	var b int = int(a)

	fmt.Println(y, b) // 3 -3

	var aa int16 = 258
	bb := int8(aa)

	fmt.Println(bb) // 2

	var xx int16 = 255
	yy := int8(xx)

	fmt.Println(yy) // -1

	var xxx int64 = 9007199254740993
	yyy := float64(xxx)

	fmt.Printf("%.0f\n", yyy)
}
