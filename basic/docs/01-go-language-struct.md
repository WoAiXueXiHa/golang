# Go 项目结构

## 这一章要记住什么

这一章主要看三个点：

- `package main` 是程序入口包。
- `main()` 是可执行程序的入口函数。
- 自己写的包可以通过 `import` 引入，导出的函数名必须首字母大写。

---

## 1. main 包和 main 函数

代码里入口文件是：

```go
package main

func main() {
	fmt.Println("hello, golang")
}
```

`package main` 表示这是一个可以编译成可执行程序的包。

`func main()` 是程序真正开始执行的地方。

```text
go run main.go
      |
      v
找到 package main
      |
      v
执行 func main()
```

## 总结一下

写 Go 可执行程序时，入口文件必须属于 `main` 包，并且要有 `main()` 函数。

---

## 2. import 引入标准库和自定义包

代码里引入了两个包：

```go
import (
	"fmt"
	myMath "learn/golang/basic/01-go-language-struct/myMath"
)
```

`fmt` 是标准库，负责格式化输入输出。

`myMath` 是自己写的包路径别名。真实包名在文件里是：

```go
package mymath
```

这里用了别名：

```go
myMath "learn/golang/basic/01-go-language-struct/myMath"
```

所以调用时写：

```go
myMath.Add(1, 3)
myMath.Mul(3, 3)
```

```text
main.go
  |
  | import myMath
  v
myMath 包
  |
  | 导出函数
  v
Add / Mul
```

## 总结一下

`import` 后面写的是包路径，代码里使用的是包名或别名。

如果路径名、目录名、包名不完全一样，可以通过别名把调用方式固定下来。

---

## 3. 函数名首字母大写表示导出

自定义包里有两个函数：

```go
func Add(x, y int) int {
	return x + y
}

func Mul(x, y int) int {
	return x * y
}
```

`Add` 和 `Mul` 首字母都是大写，所以可以被其他包访问。

如果写成 `add` 或 `mul`，只能在当前包内部使用。

```text
首字母大写：Add  -> 包外可以访问
首字母小写：add  -> 只能包内访问
```

## 总结一下

Go 不用 `public`、`private` 这种关键字。

首字母大写就是导出，首字母小写就是包内私有。

---

## 易错点

1. `package main` 和 `func main()` 是可执行程序入口，不是所有包都需要。
2. `import` 写的是包路径，调用时用的是包名或别名。
3. 自定义包里的函数想给外部用，函数名必须首字母大写。
4. 同一个目录下的 Go 文件应该属于同一个包。

---

## 快问快答

### Q1：`package main` 有什么作用？

答：

它表示当前包是可执行程序入口包，可以编译运行。

### Q2：Go 里怎么控制函数能不能被其他包访问？

答：

看首字母。大写可以被包外访问，小写只能在当前包内访问。

### Q3：`import myMath "xxx/myMath"` 是什么意思？

答：

这是给导入的包起别名。后面调用这个包里的导出内容时，用 `myMath.Add()` 这种形式。

---

## 一句话总结

Go 项目的入口靠 `package main` 和 `main()`，包之间复用靠 `import`，对外暴露能力靠首字母大写。

