package main

import (
	"fmt"
)

type person struct {
	name string
	age  int
}

func newPerson(name string) *person {
	p := person{name: name}
	p.age = 20
	// C++ 绝对不允许，因为局部变量绑定在栈上
	// Go 这样可以，逃逸了，GC 回收这个变量
	return &p
}

func main() {
	fmt.Println(person{"Bob", 25})

	fmt.Println(person{name: "AAAA", age: 30})

	fmt.Println(person{name: "BBBB"})

	fmt.Println(&person{name: "CCCC", age: 50})

	fmt.Println(newPerson("DDDD"))

	s := person{name: "EEEE", age: 70}
	fmt.Println(s.name)

	sp := &s // 保存的是 s 的地址，并没有完整复制整个结构体
	// 语法糖，指针自动解引用
	fmt.Println(sp.age)
	// 结构体可变
	sp.age = 99 // s.age 也变了
	fmt.Println(sp.age)

	// 不通过指针可变吗？
	sq := s // 值拷贝，整个结构体被复制一份，此后两个结构体独立
	fmt.Println(sq.age)
	sq.age = 11
	fmt.Println(sq.age)

}
