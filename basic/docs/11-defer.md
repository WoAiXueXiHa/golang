# Go defer

## 这一章要记住什么

这一章主要讲四个点：

- `defer` 会把函数调用延迟到当前函数退出前执行。
- 多个 `defer` 按后进先出顺序执行。
- `defer` 常用于资源释放。
- `defer` 和 `return` 配合时，要注意参数求值和返回值修改。

---

## 1. defer 执行顺序

代码里有三个函数：

```go
func defer1() { fmt.Println("defer1...") }
func defer2() { fmt.Println("defer2...") }
func defer3() { fmt.Println("defer3...") }
```

如果这样写：

```go
defer defer1()
defer defer2()
defer defer3()
```

执行顺序是：

```text
defer3
defer2
defer1
```

```text
注册顺序：
defer1 -> defer2 -> defer3

执行顺序：
defer3 -> defer2 -> defer1

像栈一样：
+--------+
| defer3 | 先执行
| defer2 |
| defer1 | 后执行
+--------+
```

## 总结一下

多个 `defer` 是后进先出，最后注册的最先执行。

---

## 2. defer 用于资源释放

代码里有两个复制文件的函数。

问题版本：

```go
src, err := os.Open(srcFile)
if err != nil {
	return
}
dst, err := os.Create(dstFile)
if err != nil {
	return
}
```

如果 `os.Create` 失败，`src` 已经打开，但没有关闭。

```text
打开 src 成功
    |
    v
创建 dst 失败
    |
    v
直接 return
    |
    v
src 没有 Close
```

改进版本：

```go
src, err := os.Open(srcFile)
if err != nil {
	return
}
defer src.Close()
```

只要资源打开成功，马上注册 `defer Close()`。

## 总结一下

资源谁打开，谁负责关闭。

打开成功后立刻 `defer Close()`，能减少中途返回导致资源泄露的风险。

---

## 3. defer 参数会提前求值

代码里：

```go
func deferRun1() {
	num := 1
	defer fmt.Printf("num is %d\n", num)

	num = 2
	return
}
```

输出的是 `1`，不是 `2`。

因为 `defer fmt.Printf(...)` 注册时，参数 `num` 已经被求值了。

```text
num = 1
defer fmt.Printf(..., num)
        |
        v
此时把 num 的值 1 记录下来

num = 2
return
执行 defer，打印记录好的 1
```

## 总结一下

`defer` 延迟的是函数调用的执行，但函数参数会在注册 `defer` 时就计算好。

---

## 4. defer 引用指针时看到的是最终数据

代码里：

```go
func deferRun2() {
	arr := [4]int{1, 2, 3, 4}
	defer printArr(&arr)

	arr[0] = 999
	return
}
```

这里 `defer` 传进去的是数组地址。

后面修改了数组内容，`defer` 执行时顺着地址看到的是修改后的数组。

```text
arr 地址传给 defer
      |
      v
arr[0] 改成 999
      |
      v
defer 执行时通过地址读取
      |
      v
看到 999
```

## 总结一下

如果 `defer` 参数是指针，注册时保存的是地址。

地址指向的数据后面变了，`defer` 执行时就会看到变化后的结果。

---

## 5. defer 和 return

代码里：

```go
func deferRun3() (res int) {
	num := 555
	defer func() {
		res++
	}()
	return num
}
```

这是有名返回值。

执行过程：

```text
return num
   |
   v
先把 num 赋给 res：res = 555
   |
   v
执行 defer：res++
   |
   v
真正返回 res：556
```

另一个函数：

```go
func deferRun4() int {
	num := 777
	defer func() {
		num++
	}()
	return num
}
```

这是无名返回值。

```text
return num
   |
   v
先把 num 的值 777 放进临时返回区
   |
   v
执行 defer：num++，局部变量变 778
   |
   v
返回临时返回区里的 777
```

## 总结一下

有名返回值可以在 `defer` 里被修改。

无名返回值会先把返回结果保存到临时位置，`defer` 改局部变量通常影响不到最终返回值。

---

## 易错点

1. `defer` 是当前函数退出前执行，不是当前代码块结束时执行。
2. 多个 `defer` 后进先出。
3. `defer` 的参数会在注册时求值。
4. 传指针给 `defer` 时，执行时能看到指针指向数据的后续变化。
5. 有名返回值可能被 `defer` 修改。

---

## 快问快答

### Q1：defer 一般用来做什么？

答：

最常见的是资源释放，比如关闭文件、释放锁、关闭连接。

### Q2：多个 defer 的执行顺序是什么？

答：

后进先出，最后注册的最先执行。

### Q3：`defer fmt.Println(num)` 里的 `num` 什么时候确定？

答：

注册 `defer` 时就确定，不是执行 `defer` 时才取值。

### Q4：defer 能修改返回值吗？

答：

如果是有名返回值，`defer` 里可以修改它，并影响最终返回结果。

---

## 一句话总结

`defer` 负责把收尾动作放到函数退出前执行，记住后进先出、参数提前求值、有名返回值可被修改。

