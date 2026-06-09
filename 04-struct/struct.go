package main

import "fmt"

func testStruct() {
	// 定义
	type Student struct {
		ID    int
		Name  string
		Age   int
		Score int
	}
	// 初始化
	// 键值对初始化， 属性: 值的方式，有的属性不写，默认值
	st := Student{
		ID:    100,
		Name:  "kun",
		Age:   29,
		Score: 100,
	}
	fmt.Printf("学生st: %v\n", st)

	// 值列表初始化，对应位置初始化，这种情况必须赋值完
	stu := &Student{
		101,
		"li",
		20,
		99,
	}
	fmt.Printf("学生stu: %v\n", stu)

	// 成员访问，用 . 访问，可以访问结构体变量或结构体指针
	fmt.Printf("学生st的分数: %d\n", st.Score)
	fmt.Printf("学生stu的名字: %s\n", st.Name)
}

func main() {
	testStruct()
}
