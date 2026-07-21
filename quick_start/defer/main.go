package main

import (
	"fmt"
)

// defer：把一个函数调用延迟到当前函数结束前执行

// 1. defer 什么时候执行，执行顺序
// 当前函数即将退出时执行，FILO
func demo1() {
	fmt.Println("1. enter demo1")

	defer fmt.Println("4. first defer")
	defer fmt.Println("3. second defer")

	fmt.Println("2. prepare to leave demo1")
}

// 2. defer 的参数会立刻求值
func demo2() {
	x := 10
	defer fmt.Println("defer1 saw x = ", x)

	// 想最后再读取变量，用匿名函数：
	defer func() {
		fmt.Println("defer2 saw x = ", x)
	}()

	x = 20
	fmt.Println("In main, x = ", x)
}

// 3. defer 使用场景：释放资源
type File struct {
	name string
}

// 注意：不能在函数内部定义方法，要外置
func (f *File) Close() {
	fmt.Println("close the file: ", f.name)
}

func demo3() {
	OpenFile := func(name string) *File {
		fmt.Println("open the file: ", name)
		return &File{name: name}
	}

	readFile := func() {
		file := OpenFile("data.txt")
		defer file.Close()

		fmt.Println("reading file")
		// 中途 return，也会执行 defer
		return
	}

	readFile()
}

// 4. defer + recover 捕获 panic
func demo4() {
	defer func() {
		// 注意：recover 必须放在 defer func(){} 里面才有意义
		err := recover()
		if err != nil {
			fmt.Println("catch panic: ", err)
		}
	}()

	fmt.Println("begin to run ")
	panic("the bad error")
	fmt.Println("This line won't be executed")
}

// 5.panic 会沿着调用栈向上传递
func a() {
	defer fmt.Println("a defer")
	b()
}

func b() {
	defer fmt.Println("b defer")
	c()
}

func c() {
	defer fmt.Println("c defer")
	// panic 从 c 开始往上炸：
	// c -> b -> a -> main
	// 每层函数退出前，都会先执行自己的defer
	// 如果某层 defer 里 recover 了，panic 就停止传播
	panic("crash")
}

// 6. defer 和 return -> return 先赋值，defer 后执行，最后真正返回
func demo6() {
	// 普通返回值
	fmt.Println("Nommal return value")
	f := func() int {
		x := 10

		defer func() { // defer 修改的是局部变量 x，最后真正返回 10
			x = 20
			fmt.Println("Modify defer, x = ", x)
		}()

		return x // return x 先把返回值定成10， 然后执行 defer
	}
	fmt.Println("f retrun: ", f())

	// 命名返回值
	fmt.Println("Named return value")
	g := func() (x int) {
		x = 10

		defer func() {
			x = 20
			fmt.Println("Modifer defer x = ", x)
		}()

		return // 这里的 x 本身就是返回值变量
	}
	fmt.Println("g return: ", g())
}
func main() {
	// demo1()
	// fmt.Println("5. back to main")
	// demo2()
	// demo3()
	// demo4()
	// fmt.Println("exe is running...")

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("main catch panic: ", err)
	// 	}
	// }()
	// a()
	// fmt.Println("This line won't be excuted")

	demo6()
}
