package main

import (
	"fmt"
	"io"
	"os"
)

// 1. defer 执行顺序
func defer1() {
	fmt.Println("defer1...")
}

func defer2() {
	fmt.Println("defer2...")
}

func defer3() {
	fmt.Println("defer3...")
}

// 2. defer 使用：资源释放
// 这段代码是有问题的，第 30 行执行失败，程序直接 return 还没有关闭打开的文件 src 呢
func BadCopyFile(dstFile, srcFile string) (wr int64, err error) {
	src, err := os.Open(srcFile)
	if err != nil {
		return
	}
	dst, err := os.Create(dstFile)
	if err != nil {
		return
	}

	wr, err = io.Copy(dst, src)
	dst.Close()
	src.Close()
	return
}

func GoodCopyFile(dstFile, srcFile string) (wr int64, err error) {
	src, err := os.Open(srcFile)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstFile)
	if err != nil {
		return
	}
	defer dst.Close()

	wr, err = io.Copy(dst, src)
	return wr, err
}

// 4. defer 和 return
func deferRun1() {
	num := 1
	defer fmt.Printf("num is %d\n", num)

	num = 2
	return
}

func printArr(arr *[4]int) {
	for i := range arr {
		fmt.Println(arr[i])
	}
}

func deferRun2() {
	arr := [4]int{1, 2, 3, 4}
	defer printArr(&arr)

	arr[0] = 999
	return
}

func deferRun3() (res int) {
	num := 555

	defer func() {
		res++
	}()

	return num
}

func deferRun4() int {
	num := 777
	defer func() {
		num++
	}()

	return num
}

func main() {
	// defer defer1()
	// defer defer2()
	// defer defer3()

	// // 3. 配合 recover 一起处理 panic
	// fmt.Println("--------------- 配合 recover 一起处理 panic ----------------")
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println(r)
	// 	}
	// }()
	// a := 1
	// b := 0
	// fmt.Println("res: ", a/b)

	// 4. defer 和 return
	deferRun1()
	deferRun2()
	res := deferRun3()
	fmt.Println(res)

	ans := deferRun4()
	fmt.Println(ans)
}
