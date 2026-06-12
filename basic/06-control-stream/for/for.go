package main

import "fmt"

// 类似 C++ 的 for 循环
// for init; condition; post {
// ...
// }

// 类似 C++ 的 while
// for condition {
// ...
// }

// 死循环
// for {
// ...
// }

// for range
func testFor() {
	for i := 0; i < 5; i++ {
		fmt.Printf("cur i is %d\n", i)
	}
	j := 0
	for {
		if j == 5 {
			break
		}
		fmt.Printf("cur j is %d\n", j)
		j++
	}

	var strArr = []string{"aa", "bb", "cc", "dd"}
	var sliceArr = make([]string, 0)
	sliceArr = strArr[1:3]
	for i, s := range sliceArr {
		fmt.Printf("slice index is %d, str is %s\n", i, s)
	}

	var dic = map[string]int{
		"aa":      1,
		"bnbbbbb": 2,
	}
	for k, v := range dic {
		fmt.Printf("key is %s, value is %d\n", k, v)
	}
}

func main() {
	testFor()
}
