## 🖼️ 核心视觉视窗：Go Error 与接口比对底层内存 ASCII 全景图

在拆解代码前，请先把这幅内存对比图牢牢刻在脑海里。为什么 `errors.New` 比较是 `false`，而你的自定义错误比较是 `true`，秘密全在这张图里：

```text
======================================================================================
                          【Error 接口比对的内存真相】
======================================================================================

1. 标准库 errors.New("hello")                    2. 你自定义的 MyError{100, "test"}
   (返回的是 &errorString{} 指针)                   (返回的是 结构体值类型)

 第一次调用 err3 := errors.New("hello")          第一次调用 errA := NewError(100, "test")
 err3 (iface 16字节)                             errA (iface 16字节)
 ┌───────────────────────────────┐               ┌───────────────────────────────┐
 │ tab *itab  │   data unsafe.Ptr│               │ tab *itab  │   data unsafe.Ptr│
 └──────┬───────────────┬────────┘               └──────┬───────────────┬────────┘
        │               │                               │               │
        ▼               ▼                               ▼               ▼
    【itab表】     【堆内存地址A】                     【itab表】     【堆内存地址X】
 ┌─────────────┐  ┌─────────────┐                ┌─────────────┐  ┌──────────────────┐
 │ type:       │  │ s: "hello"  │                │ type:       │  │ code: 100        │
 │ *errorString│  └─────────────┘                │  MyError    │  │ msg:  "test"     │
 └─────────────┘                                 └─────────────┘  └──────────────────┘
                                                                        ▲
 第二次调用 err4 := errors.New("hello")          第二次调用 errB := NewError(100, "test")
 err4 (iface 16字节)                             errB (iface 16字节)              │ (Go运行时判定
 ┌───────────────────────────────┐               ┌───────────────────────────────┐ │ 字段值完全相同)
 │ tab *itab  │   data unsafe.Ptr│               │ tab *itab  │   data unsafe.Ptr│ │
 └──────┬───────────────┬────────┘               └──────┬───────────────┬────────┘ │
        │               │                               │               │          │
        ▼               ▼                               ▼               ▼          │
    【itab表】     【堆内存地址B】                     【itab表】     【堆内存地址Y】    │
 ┌─────────────┐  ┌─────────────┐                ┌─────────────┐  ┌──────────────────┘
 │ type:       │  │ s: "hello"  │                │ type:       │  │ code: 100        │
 │ *errorString│  └─────────────┘                │  MyError    │  │ msg:  "test"     │
 └─────────────┘                                 └─────────────┘  └──────────────────┘

 💥 接口比对机制：                                💥 接口比对机制：
 tab地址相同 (✅)                                 tab地址相同 (✅)
 data地址不同 (❌ 0xAddrA != 0xAddrB)            data地址指向的【结构体内容完全相同】(✅)
 ────────────────────────────────                ───────────────────────────────────
 🎯 结果：err3 == err4 返回 FALSE                 🎯 结果：errA == errB 返回 TRUE ！

```

---

## 🛠️ 模块一：error 标准库的基本使用与“不等之谜”

### 1. 完整代码案例

```go
package main

import (
	"errors"
	"fmt"
)

// 1. error 基本使用
func getPositiveSelfAdd(num int) (int, error) {
	if num <= 0 {
		// fmt.Errorf 底层也是通过包装返回一个错误的指针对象
		return -1, fmt.Errorf("num is not a positive number")
	}
	return num + 1, nil
}

func main() {
	fmt.Println("--------------- 1. error 基本使用 -------------")
	// 注意：你原代码中 num2 的入参也是 1，无法触发错误分支，这里将其改为 0 触发错误
	num1, err1 := getPositiveSelfAdd(1)
	fmt.Printf("num is %d, err is %v\n", num1, err1)

	num2, err2 := getPositiveSelfAdd(0)
	fmt.Printf("num is %d, err is %v\n", num2, err2)

	err3 := errors.New("hello")
	err4 := errors.New("hello")
	
	// ⚡ 震撼的真相点：
	fmt.Println(err3 == err4) // 打印结果：false！
	
	// 想要比较两个 error 的内容，需要通过 Error() 拿到字符串信息进行值比对
	fmt.Println(err3.Error() == err4.Error()) // 打印结果：true
}

```

### 2. 代码逐行解释

* `type error interface { Error() string }`：这是 Go 语言最根基的内置接口。任何实现了 `Error() string` 的类型都可以当作错误对象传递。
* `return -1, fmt.Errorf(...)`：当条件不满足时，函数返回特殊的错误信息，配合 `nil`（无错误状态）作为标准返回值组合。
* `err3 := errors.New("hello")`：调用标准库创建错误。你看你在注释里贴的源码：它底层返回的是一个**结构体指针** `&errorString{text}`。
* `err3 == err4`：在进行接口对象的直接比对。

### 3. 底层内存原理 🧠

配合 **【ASCII 图第 1 部分】**。为什么同样的字符串，`err3 == err4` 会是 `false`？

* Go 语言的接口变量比对铁律是：**只有当 `tab`（类型）相同，且 `data`（数据指针）也相同时，接口才相等。**
* 既然 `errors.New` 每次都用 `&` 去堆上申请一块**新内存**来存放 `errorString`，那么 `err3` 的 `data` 指针（比如 `0x00A1`）和 `err4` 的 `data` 指针（比如 `0x00B2`）在物理内存上**完全不相同**！
* **设计哲学**：Go 官方故意让它们不相等。防止你在不同的业务包里，由于不小心定义了同名字符串的错误（比如都叫 `"permission denied"`），导致在底层误判为是同一个错误类型。

---

## 🛠️ 模块二：自定义 error 对象与类型断言类型安全

### 1. 完整代码案例

```go
package main

import "fmt"

// 2. 自定义 error 对象
type MyError struct {
	code int
	msg  string
}

// 隐式实现了内置的 error 接口（注意：这里使用的是值接收者）
func (m MyError) Error() string {
	return fmt.Sprintf("code=%d, msg=%v", m.code, m.msg)
}

func NewError(code int, msg string) error {
	// 🎯 注意：这里直接返回了结构体值类型
	return MyError{
		code: code,
		msg:  msg,
	}
}

func Code(err error) int {
	// 使用上一章学过的接口安全类型断言
	if e, ok := err.(MyError); ok {
		return e.code
	}
	return -1
}

func Msg(err error) string {
	// 使用接口安全类型断言还原真实结构体
	if e, ok := err.(MyError); ok {
		return e.msg
	}
	return ""
}

func main() {
	fmt.Println("--------------- 2. 自定义 error 对象 -------------")
	err := NewError(100, "test MyError")
	fmt.Printf("code is %d, msg is %s\n", Code(err), Msg(err))
	
	// ⚡ 课外硬核延伸：
	errA := NewError(100, "test MyError")
	errB := NewError(100, "test MyError")
	fmt.Printf("自定义错误比对结果: %t\n", errA == errB) // 打印结果：true ！！
}

```

### 2. 代码逐行解释

* `type MyError struct { code int; msg string }`：大厂规范中常用的带自定义状态码（如 403、500）的错误结构体。
* `func (m MyError) Error() string`：实现接口。因为是用**值接收者**绑定的，所以 `MyError` 的值类型和指针类型都自动实现了接口。
* `if e, ok := err.(MyError); ok`：我们上一章学过的 **`comma, ok` 类型断言**。把抽象的 `error` 接口外壳剥开，重新在内存里还原成你特有的 `MyError` 结构体，进而捞出 `code` 和 `msg` 字段。

### 3. 底层内存原理 🧠

配合 **【ASCII 图第 2 部分】**。为什么在这里，同样的参数，`errA == errB` 却变成了 `true` 呢？

* **大厂面试必考点**：因为你在 `NewError` 里返回的是**值类型** `MyError{...}`，而不是指针 `&MyError{...}`！
* 当值类型被包裹进 `iface` 的 `data` 指针时，Go 语言在执行 `==` 接口比对时，如果发现动态类型是值类型，运行时会**进一步去解引用，比对它们内部的字段内容**！
* 由于 `errA` 和 `errB` 的 `code` 都是 `100`，`msg` 都是 `"test MyError"`，内容完全一致，所以接口判定它们相等！

---

## 🏗️ 架构师避坑：大厂经典的“Nil Error 线上灾难”

既然我们已经把 `error` 接口和双指针原理玩得炉火纯青了，我们来复盘一个在大厂微服务开发中**极具毁灭性的经典 Bug**。

请看下面这段代码：

```go
package main

import "fmt"

type CustomError struct{}
func (e *CustomError) Error() string { return "custom error" }

func CheckValid() error {
    var res *CustomError = nil 
    return res // 🧩 这里的 res 确实是个 nil 指针
}

func main() {
    err := CheckValid()
    if err != nil {
        fmt.Println("❌ 触发警报：系统判定发生了错误！") // 问题：这行会打印出来吗？
    } else {
        fmt.Println("✅ 完美：没有任何错误。")
    }
}

```
