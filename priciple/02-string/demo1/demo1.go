package main

import "fmt"

func main() {
	word := "Hello, Word"
	for _, v := range word {
		fmt.Printf("%d\n", v)
	}
}
