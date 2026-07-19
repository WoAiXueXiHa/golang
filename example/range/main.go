package main

import (
	"fmt"
	"sync"
)

func demo0() {
	nums := []int{2, 3, 4}
	sum := 0
	// 索引i 值nums[i]
	for _, num := range nums {
		sum += num
	}
	fmt.Println("sum:", sum)

	for i, num := range nums {
		if num == 3 {
			fmt.Println("index:", i)
		}
	}

	kvs := map[string]string{
		"a": "apple",
		"b": "banana",
	}
	for k, v := range kvs {
		fmt.Printf("%s -> %s\n", k, v)
	}

	for k := range kvs {
		fmt.Println("key:", k)
	}

	// 字符串中迭代的是 unicode 码点
	for i, v := range "go" {
		fmt.Println(i, v)
	}
}

// for range 的坑
// 1. &v 拿到的是迭代变量的地址，而不是切片元素的地址
// 1.22 之前，v 是同一个临时变量反复用
func demo1() {
	nums := []int{1, 2, 3}
	for _, v := range nums {
		fmt.Printf("%p -> %d\n", &v, v)
		// 1.21 及以前，地址一样
		// 1.22 及以后，地址可能不一样，每轮都有新的 v
	}

	// 想拿原切片元素地址，这样：
	for i := range nums {
		fmt.Printf("%p -> %d\n", &nums[i], nums[i])
	}
}

// 2. 闭包捕获迭代变量
func demo2() {
	vals := []string{"a", "b", "c"}
	var wg sync.WaitGroup

	for _, v := range vals {
		wg.Add(1)
		v := v // 1.22 以前必须这样：创建新变量，让闭包捕获当前迭代的副本
		go func() {
			defer wg.Done()
			fmt.Println(v)
		}()
	}

	wg.Wait() // 等所有 goroutine 执行完
}

// 3.修改迭代变量，不会改原切片，因为 v 是元素副本
// 想修改，必须用索引遍历
func demo3() {
	nums := []int{1, 2, 3}
	for _, v := range nums {
		v *= 10
	}
	fmt.Println(nums)
}

// 4. map 遍历顺序随机，每次顺序可能不同
func demo4() {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	for k, v := range m {
		fmt.Println(k, v)
	}
}

// 5. 字符串 range 返回的是 rune，不是 byte
func demo5() {
	s := "Go语言"
	// 一个汉字占3字节，按照 unicode 码点遍历
	for i, r := range s {
		fmt.Println(i, string(r))
	}
	// 按照单字节遍历
	for i := 0; i < len(s); i++ {
		fmt.Println(i, s[i])
	}
	// 按照单字符遍历
	for i, r := range []rune(s) {
		fmt.Println(i, string(r))
	}
}

// 6. range 中删除切片，容易跳过或漏删
func demo6() {
	// fmt.Println("错误示范：")
	// // range 在开始的时候就确定了遍历节奏，但在循环中改变了切片结构
	// nums := []int{1, 2, 3, 4, 5}

	// for i, v := range nums {
	// 	if v%2 == 0 {
	// 		nums = append(nums[:i], nums[i+1:]...)
	// 	}
	// }

	// fmt.Println(nums)

	// 正确方式：尾删
	nums := []int{1, 2, 3, 4, 5}

	for i := len(nums) - 1; i >= 0; i-- {
		if nums[i]%2 == 0 {
			nums = append(nums[:i], nums[i+1:]...)
		}
	}

	fmt.Println(nums)
}

func main() {
	//demo1()
	//demo2()
	//demo3()
	//demo4()
	//demo5()
	demo6()
}
