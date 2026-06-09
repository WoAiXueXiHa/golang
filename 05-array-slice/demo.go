package main

import "fmt"

func main() {
	// 数组初始化
	// var 变量名 = [数量写死了]类型{元素列表}
	var strArr = [10]string{"aa", "bb", "cc"}

	// 变量名 := [数量]类型{元素列表}

	arr := [3]int{1, 2, 3}

	// 切片初始化
	// []string 这个整体是切片的类型
	// make([]类型, len)
	var sliceArr = make([]string, 0)
	sliceArr = strArr[1:3]
	// map 初始化
	var dic = map[string]int{
		"apple":  1,
		"banana": 2,
	}

	fmt.Printf("strArr没有+v: %v\n", arr)
	fmt.Printf("strArr: %+v\n", strArr)
	fmt.Printf("arr: %+v\n", arr)
	fmt.Printf("sliceArr: %+v\n", sliceArr)
	fmt.Printf("dic: %+v\n", dic)

}
