package main

import "fmt"

// 1. 定义接口
type Speaker interface {
	Speak() string
}

// 2. 隐式实现
type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return d.Name + ": wang"
}

// 3. 指针接收
type Cat struct {
	Name string
}

func (c *Cat) Speak() string {
	return c.Name + ": miao"
}

// 4. 接口组合
type Named interface {
	GetName() string
}

type NamedSpeaker interface {
	Speaker
	Named
}

func (d Dog) GetName() string {
	return d.Name
}
func (c Cat) GetName() string {
	return c.Name
}

// 5. interface 作为参数
func PrintSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

// 6. 类型断言
func CheckString(v any) {
	s, ok := v.(string)
	if ok {
		fmt.Println("string:", s)
	}
}

// 7. 类型 switch
func PrintType(v any) {
	switch x := v.(type) {
	case int:
		fmt.Println("int: ", x)
	case string:
		fmt.Println("string: ", x)
	case Speaker:
		fmt.Println("Speaker: ", x.Speak())
	default:
		fmt.Println("unknown")
	}
}

func main() {
	// 普通实现
	dog := Dog{Name: "旺财"}
	PrintSpeak(dog)

	// 指针接收
	cat := &Cat{Name: "咪咪"}
	PrintSpeak(cat)

	// 接口变量
	var s Speaker = dog
	fmt.Println(s.Speak())

	// 空接口，任意类型
	var v any
	v = 100
	fmt.Println(v)

	v = "hello"
	fmt.Println(v)

	// 类型断言
	CheckString(v)

	// 类型 switch
	PrintType(123)
	PrintType("hello")
	PrintType(dog)
	PrintType(cat)

	var ns NamedSpeaker = dog
	fmt.Println(ns.GetName())
	fmt.Println(ns.Speak())

	// 空接口
	var nilSpeaker Speaker
	fmt.Println(nilSpeaker == nil)

	// typed nil 放入 interface
	var nilCat *Cat = nil
	var speaker Speaker = nilCat
	fmt.Println(speaker == nil)

	// 空接口可以接收任何类型
	// 但是空类型不能接收空接口
}
