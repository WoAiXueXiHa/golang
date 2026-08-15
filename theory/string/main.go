package main

import (
	"fmt"
	"strings"
)

func defintion() {
	s := "你好"
	fmt.Println(len(s)) // 6
	fmt.Println(s[0])   // 228，一个 byte
}

func exchange() {
	s := "hello"
	sByte := []byte(s)
	sByte[0] = 'H'
	fmt.Println(string(sByte))
}

func add() {
	item := []string{"Go", " ", "is", " ", "great"}
	var b strings.Builder

	for _, item := range item {
		b.WriteString(item)
	}

	fmt.Println(b.String())
}

func main() {
	//defintion()
	exchange()
}
