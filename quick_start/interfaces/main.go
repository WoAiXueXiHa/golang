package main

import (
	"fmt"
)

// 接口：方法签名的集合，一份能力合同

// 任何拥有 Speak() string 方法的类型，都可以作为 Speaker
// 只关心方法，不关心结构体字段
type Speaker interface {
	Speak() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return d.Name + ": 汪汪"
}

// 方法名不能靠“值接收者/指针接收者”来重载
// func (d *Dog) Speak() string {
// 	return d.Name + ": 汪汪"
// }

// 怎么选择值/指针？ 不修改对象->值，修改对象->指针，要传地址

type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return c.Name + ": 喵喵"
}

// 接口作为函数参数，s 参数必须具备 Speak() string 的方法
func MakeSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

// 空接口 interface{} any 两者等价
// 可以理解为 C++ 的 void*，Go 的空接口可以接收任何类型
func PrintValue(v any) {
	fmt.Println(v)
	// 但是在函数内部，v 的静态类型是 any， 不能直接调用原类型的方法
	// fmt.Println(v.Speak()) // 错误，any 没有声明 Speak 方法
}

// 接口变量中可能保存不同类型，断言用于把接口值取回具体类型
// value := interfaceValue.(具体类型)

func PrintDog(v any) {
	dog, ok := v.(Dog)
	// 不要随意写 dog := v.(Dog)
	// 如果实际类型不是 Dog，程序会发生 panic
	if !ok {
		fmt.Println("不是 Dog")
		return
	}
	fmt.Println(dog.Name)
}

// 也可以断言接口
func TrySpeak(v any) {
	speaker, ok := v.(Speaker)
	if !ok {
		fmt.Println("不会说话")
		return
	}
	fmt.Println(speaker.Speak())
}

// 接口嵌套
type Reader interface {
	Read() string
}

type Writer interface {
	Write(string)
}

type ReadWriter interface {
	Reader
	Writer
}

// 实现 ReadWriter 的类型必须同时拥有读写两个方法

type File struct {
	Content string
}

func (f *File) Read() string {
	return f.Content
}
func (f *File) Write(content string) {
	f.Content = content
}

func SaveAndPrint(rw ReadWriter) {
	rw.Write("hello go")
	fmt.Println(rw.Read())
}

func main() {
	fmt.Println("--- 接口基本使用 ---")

	var s Speaker

	s = Dog{Name: "旺财"}
	fmt.Println(s.Speak())

	s = Cat{Name: "咪咪"}
	fmt.Println(s.Speak())

	// 接口作为函数参数，调用时可以传入任何实现了 Speaker 的类型
	dog := Dog{Name: "大黄"}
	cat := Cat{Name: "小咪"}
	MakeSpeak(dog)
	MakeSpeak(cat)
	// 此时函数不再依赖具体的 Dog/Cat，只依赖 Speaker

	fmt.Println("--- 空接口 ---")
	PrintValue(100)
	PrintValue("hello")
	PrintValue(true)
	PrintValue(Dog{Name: "旺财"})

	fmt.Println("--- 断言 ---")
	PrintDog(Dog{Name: "旺财"}) // 旺财
	PrintDog("hello")         // 不是 Dog

	fmt.Println("--- 接口嵌套 ---")
	file := &File{}
	SaveAndPrint(file)
}

// 总结
// 1. 接口是一组方法的清单要求
// 2. 方法匹配自动实现接口
// 3. 接口作为函数参数，让函数接收不同实现
// 4. any 等价于 interface{}，可以接收任何类型，但是不能给具体类型赋值
// 5. 类型断言用于从接口中取出具体的类型
// 6. 接口嵌套用于组合多个能力
