package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	const s = "你好 Go"

	// 一个汉字占用三个字节
	fmt.Println("len: ", len(s))
	// 使用索引，会在每个索引处生成原始字节值
	for i := 0; i < len(s); i++ {
		fmt.Printf("%x ", s[i])
	}
	fmt.Println()
	// 字符个数
	fmt.Println("rune count: ", utf8.RuneCountInString(s))
	// 按字符遍历
	for i, v := range s {
		fmt.Printf("%#U starts at %d\n", v, i)
	}

	fmt.Println("\nDecodeRuneInString")
	for i, w := 0, 0; i < len(s); i += w {
		v, width := utf8.DecodeRuneInString(s[i:])
		fmt.Printf("%#U starts at %d\n", v, i)
		w = width

		examineRune(v)
	}
}

func examineRune(r rune) {
	if r == 'x' {
		fmt.Println("found x")
	} else if r == '你' {
		fmt.Println("found 你")
	}
}
