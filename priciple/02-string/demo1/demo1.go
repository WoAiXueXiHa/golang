package main

import "fmt"

func main() {
	word := "你好, Go"
	for _, v := range word {
		fmt.Printf("%d:%c\n", v, v)
	}
}
