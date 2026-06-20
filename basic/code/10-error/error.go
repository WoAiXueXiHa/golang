package main

import (
	"errors"
	"fmt"
)

// errors.New() 源码
// func New(text string) error {
// 	return &errorString{text}
// }

// type errorString struct {
// 	s string
// }

// func (e *errorString) Error() string {
// 	return e.s
// }

// 1. error 基本使用
func getPositiveSelfAdd(num int) (int, error) {
	if num <= 0 {
		return -1, fmt.Errorf("num is not a positive number")
	}
	return num + 1, nil
}

// 2. 自定义 error 对象
type MyError struct {
	code int
	msg  string
}

func (m MyError) Error() string {
	return fmt.Sprintf("code=%d, msg=%v", m.code, m.msg)
}

func NewError(code int, msg string) error {
	return MyError{
		code: code,
		msg:  msg,
	}
}

func Code(err error) int {
	if e, ok := err.(MyError); ok {
		return e.code
	}
	return -1
}

func Msg(err error) string {
	if e, ok := err.(MyError); ok {
		return e.msg
	}
	return ""
}

func main() {
	fmt.Println("--------------- 1. error 基本使用 -------------")
	num1, err1 := getPositiveSelfAdd(1)
	fmt.Printf("num is %d, err is %v\n", num1, err1)

	num2, err2 := getPositiveSelfAdd(1)
	fmt.Printf("num is %d, err is %v\n", num2, err2)

	err3 := errors.New("hello")
	err4 := errors.New("hello")
	fmt.Println(err3 == err4)
	// 想要比较两个 error， 需要通过 Error() 拿到字符串信息
	fmt.Println(err3.Error() == err4.Error())

	fmt.Println("--------------- 2. 自定义 error 对象 -------------")
	err := NewError(100, "test MyError")
	fmt.Printf("code is %d, msg is %s\n", Code(err), Msg(err))
}
