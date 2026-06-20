# Go panic 和 recover

## 这一章要记住什么

这一章主要讲三个点：

- `panic` 会中断当前函数的正常执行流程。
- `recover` 只能在 `defer` 函数里拦截 panic。
- panic 会沿着调用栈向上传递，直到被 recover 捕获或程序崩溃。

---

## 1. panic 会打断当前执行流

代码里：

```go
func testPanic3() {
	fmt.Println("testPanic3上半部分")
	panic("在testPanic3出现了panic")
	fmt.Println("testPanic3下半部分")
}
```

执行到 `panic` 后，后面的普通代码不会继续执行。

```text
testPanic3上半部分
        |
        v
panic(...)
        |
        v
testPanic3下半部分不会执行
```

## 总结一下

`panic` 不是普通错误返回，它会直接打断当前函数后续代码。

---

## 2. panic 会向上传递

调用关系是：

```go
main -> testPanic1 -> testPanic2 -> testPanic3
```

`testPanic3` 发生 panic 后，会向上一层一层传递。

```text
main
 |
 v
testPanic1
 |
 v
testPanic2
 |
 v
testPanic3
 |
 v
panic
 |
 v
向上找 defer / recover
```

如果一直没有 recover，程序最终会崩溃。

## 总结一下

panic 会沿着调用栈向上传递，不会只停留在发生 panic 的那个函数里。

---

## 3. recover 必须配合 defer

代码里：

```go
func testPanic2() {
	defer func() {
		recover()
	}()

	fmt.Println("testPanic2上半部分")
	testPanic3()
	fmt.Println("testPanic2下半部分")
}
```

`recover()` 放在 `defer` 匿名函数里，所以能捕获从 `testPanic3` 传上来的 panic。

```text
testPanic2 注册 defer recover
        |
        v
调用 testPanic3
        |
        v
testPanic3 panic
        |
        v
回到 testPanic2，执行 defer
        |
        v
recover 捕获 panic
```

注意：`testPanic2下半部分` 仍然不会执行。

因为 `testPanic2` 是在调用 `testPanic3()` 时被 panic 打断的，recover 后它不会回到那一行后面继续跑。

## 总结一下

`recover` 能停止 panic 继续向上传播，但不能让当前函数回到 panic 发生点继续执行。

---

## 4. 为什么 testPanic1 下半部分能执行

代码里：

```go
func testPanic1() {
	fmt.Println("testPanic1上半部分")
	testPanic2()
	fmt.Println("testPanic1下半部分")
}
```

因为 panic 在 `testPanic2` 的 defer 里被 recover 了。

对 `testPanic1` 来说，`testPanic2()` 这个调用已经结束，所以它可以继续执行后面的代码。

```text
testPanic1
  |
  v
调用 testPanic2
  |
  v
testPanic2 内部捕获 panic
  |
  v
testPanic2 返回
  |
  v
testPanic1 继续执行下半部分
```

最终执行效果是：

```text
程序开始
testPanic1上半部分
testPanic2上半部分
testPanic3上半部分
testPanic1下半部分
程序结束
```

## 总结一下

在哪一层 recover，panic 就在哪一层停止继续向上传。

外层调用者会认为这个函数调用已经结束，然后继续执行自己的后续逻辑。

---

## 5. recover 捕获信息

代码注释里还有这种写法：

```go
defer func() {
	if error := recover(); error != nil {
		fmt.Println("出现了panic，使用recover获取信息：", error)
	}
}()
```

这个写法更完整。

因为 `recover()` 会返回 panic 抛出的值。

```text
panic("出现panic")
        |
        v
recover() 拿到 "出现panic"
```

## 总结一下

实际写代码时，不建议只写 `recover()` 什么都不处理。

至少应该把 panic 信息记录下来，方便排查问题。

---

## 易错点

1. `recover` 必须在 `defer` 函数里调用才有效。
2. `panic` 后面的普通代码不会继续执行。
3. recover 只能阻止 panic 继续向上传播，不能让当前函数回到原位置继续执行。
4. 不要把 panic 当作普通业务错误处理方式。
5. 捕获 panic 后最好记录日志，不要静默吞掉。

---

## 快问快答

### Q1：panic 和 error 有什么区别？

答：

`error` 是普通返回值，调用方显式处理。`panic` 会中断正常流程，沿调用栈向上传播，通常用于不可恢复的异常场景。

### Q2：recover 为什么必须放在 defer 里？

答：

因为 panic 发生后，普通代码已经不再继续执行，只有已经注册的 defer 会在栈展开过程中执行。

### Q3：recover 后，当前函数会从 panic 的地方继续执行吗？

答：

不会。recover 会停止 panic 继续传播，但当前函数不会回到 panic 发生点继续执行。

### Q4：为什么 `testPanic1下半部分` 能执行？

答：

因为 panic 在 `testPanic2` 里被 recover 了。对 `testPanic1` 来说，`testPanic2()` 调用已经返回，所以它能继续往下执行。

---

## 一句话总结

`panic` 负责打断流程，`defer` 负责兜底执行，`recover` 负责在合适的位置把 panic 拦下来。

