package main

import "fmt"

// 方法是一类特殊的函数，绑定在某种类型变量上的函数
type rect struct {
	width, height int
}

// 函数：func 函数名 返回值 函数体
// 方法：func 接收者 函数名 返回值 函数体
func (r *rect) area() int {
	return r.width * r.height
}

func (r rect) perim() int {
	return 2*r.width + 2*r.height
}

func main() {
	// 可以用值、指针调用方法
	r := rect{width: 10, height: 5}
	fmt.Println("area: ", r.area())
	fmt.Println("perim:", r.perim())

	rp := &r
	fmt.Println("area: ", rp.area())
	fmt.Println("perim:", rp.perim())
}
