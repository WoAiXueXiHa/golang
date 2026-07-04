package main

import "fmt"

func main() {
	s := "Hello"
	strByte := []byte(s)
	strByte[0] = 'h'
	fmt.Println(string(strByte))
}
