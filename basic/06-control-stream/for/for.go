package main

import (
	"fmt"
)

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

// // for range
// func testFor() {
// 	for i := 0; i < 5; i++ {
// 		fmt.Printf("cur i is %d\n", i)
// 	}
// 	j := 0
// 	for {
// 		if j == 5 {
// 			break
// 		}
// 		fmt.Printf("cur j is %d\n", j)
// 		j++
// 	}

// 	var strArr = []string{"aa", "bb", "cc", "dd"}
// 	var sliceArr = make([]string, 0)
// 	sliceArr = strArr[1:3]
// 	for i, s := range sliceArr {
// 		fmt.Printf("slice index is %d, str is %s\n", i, s)
// 	}

// 	var dic = map[string]int{
// 		"aa":      1,
// 		"bnbbbbb": 2,
// 	}
// 	for k, v := range dic {
// 		fmt.Printf("key is %s, value is %d\n", k, v)
// 	}
// }

// func main() {
// 	testFor()
// }

// for range 的坑

// // 在 1.22 版本之前，编译器为了节省内存，在整个 for 循环开始前
// // 只在内存中开辟了一块空间给 v，可以想象成一个垃圾桶
// // 第一轮：把 1 扔进 v，使用 &v 标记了这个垃圾桶的地址
// // 第二轮：把 2 扔进 v，覆盖了原来的 1， 用 &v 记录同一个垃圾桶的地址
// // 循环结束后，空间里最后剩下的是 2， res 切片里，存了两个一模一样的地址，解引用去拿，自然全是 2

// // 1.22 版本及以后，每轮循环都会创建一个新的临时地址 v，可以想象成一次性纸杯
// // 第一轮循环，生成一个一次性纸杯A，装入 1， &v 记录了纸杯A的地址
// // 第二轮循环，生成新的纸杯B，装入2， 记录了纸杯B的地址
// // res[0] 指向纸杯A res[1]指向纸杯B， 拿到的就是 1 2
// func main() {
// 	arr := [2]int{1, 2}
// 	res := []*int{}

// 	for _, v := range arr {
// 		res = append(res, &v)
// 	}

// 	fmt.Println(*res[0], *res[1])

// 	// 解决老版本的问题
// 	// 1. 使用局部变量拷贝 v
// 	for _, v := range arr {
// 		r := v
// 		res = append(res, &r)
// 	}

// 	// 2. 直接使用索引获取原来的元素
// 	for k := range arr {
// 		res = append(res, &arr[k])
// 	}
// }

// // 循环是否会停止？
// func main() {
// 	v := []int{1,2,3}
// 	for i := range v {
// 		v = append(v, i)
// 	}

// 	// 一定会，在循环开始之前，会对切片 v 的长度进行一次评估
// 	// 把这个长度用于控制循环的迭代次数，之后修改切片 v 的长度，并不会影响迭代次数

// 	// // 类似以下：
// 	// v := []int{1,2,3}
// 	// length := len(v)
// 	// for i := 0; i < length; i++ {
// 	// 	v = append(v, i)
// 	// }
// }

// // 使用迭代遍历闭包的问题
// // 老版本中，for 循环里的变量 i，只在内存中占用一块空间
// // 循环阶段，存入 1 2 3，但是此时只有一个 i 的地址
// // 另外一个循环开始遍历，去内存中找 i 的地址，i 里装的是 3，打印了三个3
// func main() {
// 	var funcs []func()

// 	for i := 0; i < 3; i++ {
// 		funcs = append(funcs, func() {
// 			fmt.Println(i) // 匿名函数 func 捕获了外层的 i
// 		})
// 	}

// 	for _, f := range funcs {
// 		f()
// 	}

// 	// 解决1：使用同名局部变量遮蔽
// 	for i := 0; i < 3; i++ {
// 		i := i // 在当前循环作用域中，新建一个变量 i
// 		funcs = append(funcs, func() {
// 			fmt.Println(i) // 此时捕获的是内部每轮循环的新的变量 i，地址都不一样
// 		})
// 	}

// 	// 2. 函数参数显式拷贝，使用值传递强制发生拷贝
// 	for i := 0; i < 3; i++ {
// 		funcs = append(funcs, func(val int) { // 定义接收参数 val
// 			fmt.Println(val)
// 		}(i)) // 每次创建时，把 i 作为实参传进去
// 	}
// }

// // for range 会创建每个元素的副本，不会直接操作原始切片中的元素
// // 因此，修改迭代变量不会影响原始切片

// func main() {
// 	slice := []int{1, 2, 3}

// 	for _, v := range slice {
// 		v *= 10
// 	}
// 	fmt.Println(slice) // [1 2 3]

// 	// 使用索引访问并修改原始切片中的元素
// 	for i := range slice {
// 		slice[i] *= 10
// 	}
// 	fmt.Println(slice)
// }

// // for range 遍历字典时，遍历顺序是随机的，每次运行程序时，顺序可能不同
// func main() {
// 	dic := map[string]int{
// 		"a": 1,
// 		"b": 2,
// 		"c": 3,
// 	}

// 	for k, v := range dic {
// 		fmt.Printf("key: %s, value: %d\n", k, v)
// 	}

// 	// 可以先对键排序，再遍历
// 	keys := make([]string, 0, len(dic))

// 	for i := range dic {
// 		keys = append(keys, i)
// 	}

// 	sort.Strings(keys)

// 	for _, key := range keys {
// 		fmt.Printf("key: %s, value: %d\n", key, dic[key])
// 	}
// }

// for range 遍历字符串时，每次迭代会返回 Unicode 代码点
// 不会返回字节，如果字符串包含多字节字符，就需要注意

func main() {
	str := "hello 世界"

	for i, r := range str {
		fmt.Printf("index: %d, rune: %c\n", i, r)
	}

	// 按照字节遍历，使用常规的 for 循环
	for i := 0; i < len(str); i++ {
		fmt.Printf("index: %d, byte: %x\n", i, str[i])
	}
}
