package main

import "fmt"

func main() {
	localStr := "case3"
	if localStr == "case3" {
		fmt.Println("into true logic")
	} else {
		fmt.Println("into false logic")
	}

	var dic = map[string]int{
		"apple":  1,
		"banana": 2,
	}
	// num 是对应键值对的值，ok是bool类型，判断这个键值是否存在
	// if ...; ok{...} 先执行前面的赋值和判断，然后判读 ok 是否为true
	// if 初始化; 条件判断 {} 声明周期在这个 if 结束之后也就结束了
	if num, ok := dic["orange"]; ok {
		fmt.Printf("orange num %d\n", num)
	}
	if num, ok := dic["apple"]; ok {
		fmt.Printf("appele num %d\n", num)
	}

	switch localStr {
	case "case1":
		fmt.Println("case1")
	case "case2":
		fmt.Println("case2")
	case "case3":
		fmt.Println("case3")
	default:
		fmt.Println("default")
	}
}
