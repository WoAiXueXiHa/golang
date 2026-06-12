package main

import "fmt"

// type interfaceName interface {
// 		methondName1(参数列表) 返回类型列表
// 		methondName2(参数列表) 返回类型列表
// ...
// }

// 1. 接口基本使用
// 和结构体一样的定义方式，换个interface关键字
// 定义接口类型 Phone，接口里有两个方法
type Phone interface {
	Call()
	SendMessage()
}

type Apple struct {
	PhoneName string
}

func (this *Apple) Call() {
	fmt.Printf("%s有打电话功能\n", this.PhoneName)
}

func (this *Apple) SendMessage() {
	fmt.Printf("%s有发信息功能\n", this.PhoneName)
}

type Oppo struct {
	PhoneName string
}

func (this *Oppo) Call() {
	fmt.Printf("%s有打电话功能\n", this.PhoneName)
}

func (this *Oppo) SendMessage() {
	fmt.Printf("%s有发信息功能\n", this.PhoneName)
}

// 2. 一个类型实现多个接口
type MyWriter interface {
	MyWriter(s string)
}

type MyRead interface {
	MyReader()
}

type MyReadWriter struct {
}

func (this *MyReadWriter) MyWriter(s string) {
	fmt.Printf("call MyWriteReader MyWriter %s\n", s)
}

func (this *MyReadWriter) MyReader() {
	fmt.Println("call MyWriteReader MyReader")
}

// 5. 接口作为函数参数
type Reader interface {
	Read() int
}

type MyReader1 struct {
	a, b int
}

func (this *MyReader1) Read() int {
	return this.a + this.b
}

func DoJob(r Reader) {
	fmt.Printf("myReader is %d\n", r.Read())
}

// 如果函数的形参是空接口，实参可以是任意类型
func DoJob1(val interface{}) {
	fmt.Printf("val is %v\n", val)
}

// 6. 接口嵌套，一个接口中包含了其他接口，要实现外部接口，需要先实现内部嵌套接口的对应方法
type A interface {
	run1()
}

type B interface {
	run2()
}

// 定义嵌套接口C
type V interface {
	A
	B
	run3()
}

type Runner struct{}

func (this *Runner) run1() {
	fmt.Println("run1...")
}

func (this *Runner) run2() {
	fmt.Println("run2...")
}

func (this *Runner) run3() {
	fmt.Println("run3...")
}

func main() {
	fmt.Println("------------- 接口基本使用 ------------")
	phoneA := Apple{"apple"}
	phoneB := Oppo{"oppo"}
	phoneA.Call()
	phoneA.SendMessage()
	phoneB.Call()
	phoneB.SendMessage()

	// 多态的体现
	var phoneC Phone
	phoneC = new(Apple)                 // new 返回的是 Apple 这个结构体指针
	phoneC.(*Apple).PhoneName = "APPLE" // 接口的断言，啥类型对应啥类型的断言
	phoneC.Call()
	phoneC.SendMessage()

	fmt.Println("------------- 一个类型定义多个接口，使用多个接口 ------------")
	myRead := new(MyReadWriter)
	myRead.MyReader()

	myWriter := MyReadWriter{}
	myWriter.MyWriter("hello")

	fmt.Println("------------- 空接口 ------------")
	// 3. 空接口，空接口可以存储任意类型的数值
	// 可以用空接口作为参数，表示接收任意类型的参数
	var any interface{}
	any = 10
	fmt.Println(any)

	any = "Vect"
	fmt.Println(any)

	any = map[string]int{
		"aa": 1,
		"bb": 2,
	}
	fmt.Println(any)

	// 4. 断言
	// a := 1
	// var i interface{} = a
	// cannot use i (variable of type interface{}) as int value in variable declaration: need type assertion
	// var b int = i
	// 可以将任意类型变量赋值给空接口 interface{} 类型，但是反过来不行
	// 为了让这个操作能够完成，使用断言
	fmt.Println("------------- 断言 ------------")
	var x interface{}
	x = 9
	val, ok := x.(int) // x是啥类型，括号里就断言啥类型
	fmt.Printf("val is %d\n, ok is %t\n", val, ok)

	// 断言过程中，如果没有bool判断，断言成功程序正常运行
	// interface{}存储的值和要断言的类型不一致，报panic
	// panic: interface conversion: interface {} is string, not int
	// var y interface{}
	// y = "golang"
	// num := y.(int)
	// fmt.Println(num)

	// 5. 接口作为函数参数
	fmt.Println("------------- 接口作为函数参数 ------------")
	myReader := &MyReader1{2, 10}
	DoJob(myReader)
	v := 20
	DoJob1(v)

	// 6. 接口嵌套
	fmt.Println("------------- 接口嵌套 ------------")
	runer := new(Runner)
	runer.run1()
	runer.run2()
	runer.run3()
}
