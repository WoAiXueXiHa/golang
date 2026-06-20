# Go defer

## 这一章要记住什么

这一章主要讲五个点：

- `defer` 会把函数调用延迟到当前函数退出前执行，多个 `defer` 按后进先出顺序执行。
- `defer` 常用于资源释放——打开成功后立刻 `defer Close()`。
- `defer` 的参数在注册时就求值，传指针则能看到最终数据；闭包变量要看捕获的是什么。
- `defer` 和 `return` 配合时，有名返回值可以被 `defer` 修改，无名返回值不行。
- `defer` + `recover` 是 Go 里处理 panic 的标准组合。

---

## 1. defer 执行顺序 —— 后进先出

代码里有三个简单的 defer 函数：

```go
func defer1() { fmt.Println("defer1...") }
func defer2() { fmt.Println("defer2...") }
func defer3() { fmt.Println("defer3...") }
```

注册顺序和执行顺序相反：

```go
defer defer1()
defer defer2()
defer defer3()
```

```text
注册顺序：                 执行顺序：
defer1()  最先注册        defer3()  ← 最先执行
defer2()                  defer2()
defer3()  最后注册        defer1()  ← 最后执行

像一摞盘子（栈）：
+--------+
| defer3 | ← 最后放上去的，最先拿走
| defer2 |
| defer1 | ← 最先放上去的，最后拿走
+--------+
```

Go 内部用链表管理 defer 记录。每次执行 `defer`，把当前调用包装成一个 `_defer` 结构体，链到 goroutine 的 defer 链上。函数退出时从链头开始倒序执行。

**总结：** 多个 defer 是后进先出（LIFO），最后注册的最先执行。Go 内部用链表串联，函数返回时从头部倒序遍历。

---

## 2. defer 用于资源释放

谁打开，谁负责关闭——这是一条铁律。`defer` 让这条铁律的执行变得自然：

```go
// ❌ 问题版本：中途 return 导致资源泄露
func BadCopyFile(dstFile, srcFile string) (wr int64, err error) {
    src, err := os.Open(srcFile)
    if err != nil {
        return
    }
    dst, err := os.Create(dstFile)
    if err != nil {
        return          // ⚠️ src 已经打开了但没有 Close！
    }
    wr, err = io.Copy(dst, src)
    dst.Close()
    src.Close()
    return
}
```

```text
打开 src 成功 ✓
    |
    v
创建 dst 失败 ✗
    |
    v
直接 return → src 没人关了，文件句柄泄露
```

改进版本——打开成功后立刻 `defer Close()`：

```go
// ✅ 改进版本：打开后立刻 defer
func GoodCopyFile(dstFile, srcFile string) (wr int64, err error) {
    src, err := os.Open(srcFile)
    if err != nil {
        return
    }
    defer src.Close()   // ★ 打开成功，立刻注册关闭

    dst, err := os.Create(dstFile)
    if err != nil {
        return          // src 会在退出时由 defer 关闭
    }
    defer dst.Close()   // ★ 同样，打开成功就注册

    wr, err = io.Copy(dst, src)
    return wr, err      // 先关 dst，再关 src（defer 逆序执行）
}
```

```text
打开 src 成功 ✓ → defer src.Close()
打开 dst 成功 ✓ → defer dst.Close()
io.Copy 完成

函数退出时，defer 逆序执行：
  dst.Close()  ← 后注册的先执行 ✓
  src.Close()  ← 先注册的后执行 ✓
```

**总结：** 资源谁打开谁关闭。打开后立刻 `defer Close()`，这样无论后面有多少个 `return` 分支，资源都不会泄露。defer 的逆序特性天然保证关闭顺序和打开顺序相反（先关 dst 再关 src）。

---

## 3. defer 参数会提前求值

这是 defer 最容易踩的坑之一。看代码：

```go
func deferRun1() {
    num := 1
    defer fmt.Printf("num is %d\n", num)   // 此时 num = 1

    num = 2
    return
}
// 输出：num is 1   ← 不是 2！
```

虽然 `fmt.Printf` 的执行被延迟了，但它的**参数 `num` 在注册 defer 那一刻就已经求值完毕**。

```text
num = 1
  |
  v
defer fmt.Printf(..., num)
  |                   |
  |                   v 参数 num 求值 = 1，被复制保存
  v 注册 defer 记录
num = 2
return
  |
  v
执行 defer → 打印保存好的 1
```

`defer fmt.Printf("num is %d\n", num)` 等价于：

```go
// 伪代码：defer 立即捕获参数值
captured := num   // captured = 1，在注册时完成
// ... 后面 num 怎么变都和 captured 无关 ...
```

**总结：** `defer` 延迟的是函数调用的执行时机，不是参数求值时机。参数在 `defer` 语句执行那一刻就已经算好、复制、存好了。后面原变量怎么变，和 defer 拿到的参数没有关系。

---

## 4. defer 传指针 vs 传值

如果传给 defer 的不是值而是指针呢？

```go
func deferRun2() {
    arr := [4]int{1, 2, 3, 4}
    defer printArr(&arr)   // 传的是数组的地址 ← 注意取地址

    arr[0] = 999
    return
}
// 输出：999, 2, 3, 4  ← defer 看到了修改后的数据！
```

虽然参数（地址）在注册时就确定了，但地址指向的数据后来被改了。defer 执行时顺着地址读，自然看到最新数据。

```text
arr = [1, 2, 3, 4]  地址比如是 0xc0000180a0
  |
  v
defer printArr(&arr) → 保存地址 0xc0000180a0
  |
  v
arr[0] = 999 → 内存 [0xc0000180a0] 里的值变成 999
  |
  v
执行 defer → 读地址 0xc0000180a0 → 看到 999
```

```text
defer 保存的东西
+------------------+
| wrapper 函数指针 |     → 指向 printArr 的包装函数
| &arr             | → 指针 0xc0000180a0
+------------------+
         |
         v
  arr[0] = 999 后，这个地址里的内容就变了
  defer 执行时读地址 → 读到的就是修改后的数据
```

**对比：**

| 传值 `defer f(val)` | 传指针 `defer f(&val)` |
|---------------------|------------------------|
| 注册时复制 val 的值 | 注册时复制 val 的地址 |
| defer 看到注册时的快照 | defer 看到最新的数据 |
| 安全，不受后续修改影响 | 灵活但要注意时序 |

**总结：** 如果 defer 的参数是指针，注册时保存的是地址。地址后面指向的数据变了，defer 执行时就能看到变化后的结果。这既是特性也是陷阱——是特性还是陷阱，取决于你清不清楚自己在做什么。

---

## 5. defer 和 return

这是 defer 最精妙的部分，也是面试最爱问的。先看两个函数：

### 5.1 有名返回值 —— defer 可以改

```go
func deferRun3() (res int) {   // res 是有名返回值
    num := 555
    defer func() {
        res++                  // 闭包直接引用 res
    }()
    return num
}
// 返回值：556
```

### 5.2 无名返回值 —— defer 改了也白改

```go
func deferRun4() int {         // 无名返回值
    num := 777
    defer func() {
        num++                  // 闭包引用 num，不是返回值
    }()
    return num
}
// 返回值：777  ← num 变成 778 了，但返回的是 777！
```

### 为什么？`return` 不是原子操作

`return num` 这句话其实分三步执行：

```text
1. 赋值阶段：把 num 赋给返回值位置
   - 有名返回值：res = num     （res 就是返回值本身）
   - 无名返回值：~r0 = num     （~r0 是编译器生成的临时返回值变量）

2. defer 执行阶段：按后进先出执行所有 defer

3. 返回阶段：从返回值位置读出结果，真正 RET
```

```text
deferRun3 (有名返回值)              deferRun4 (无名返回值)

return num                           return num
    |                                    |
    v                                    v
res = 555  ← 赋值给返回值本身       ~r0 = 777  ← 赋值给临时变量
    |                                    |
    v                                    v
执行 defer: res++ → res = 556      执行 defer: num++ → num = 778
    |                               (改的是局部变量 num，~r0 纹丝不动)
    v                                    |
返回 res = 556 ✓                         v
                                    返回 ~r0 = 777  ← 白改了！
```

**核心差异一览：**

| | `deferRun3` (有名返回值) | `deferRun4` (无名返回值) |
|---|---|---|
| 闭包捕获 | `res`（返回值变量本身） | `num`（局部变量） |
| `return num` 做的事 | 直接赋值 `res = num` | 先复制 `~r0 = num`，再执行 defer |
| defer 里改的是什么 | 返回值 `res` | 局部变量 `num` |
| 最终返回值 | 被 defer 改了 ✓ | 不受 defer 影响 ✗ |

**总结：** 有名返回值可以在 defer 里修改并影响最终结果；无名返回值先把结果复制到临时位置再执行 defer，defer 里改局部变量影响不到已复制的返回值。这也是为什么 `defer` + `return err` 的惯用法（有名返回值 `err`）能正常工作的原因——defer 里可以给 `err` 附加上下文信息。

---

## 6. defer + panic + recover

源码里有一段被注释掉的 panic/recover 代码，这其实是 defer 的另一个经典用途：

```go
defer func() {
    if r := recover(); r != nil {
        fmt.Println(r)    // 捕获 panic，打印错误信息
    }
}()
a := 1
b := 0
fmt.Println("res: ", a/b)  // panic: integer divide by zero
```

```text
a / b 触发 panic
    |
    v
panic 沿调用栈向上传播
    |
    v
遇到 defer 定义的 recover()
    |
    v
recover 捕获到 panic 的值 → 程序不崩溃，继续执行 defer 后的代码
    |
    v
如果 defer 在 main 的最外层，程序正常退出
```

### 注意几个关键点：

1. **`recover` 只有在 defer 函数里才有效。** 直接写 `recover()` 或 defer 调用的非匿名函数里写 `recover()` 都返回 nil。
2. **`recover` 捕获的是 panic 的值**，可以是 string、error、任意类型。
3. **recover 后程序从触发 panic 的那个 goroutine 继续执行**，不会回到 panic 触发点。
4. **panic 只能被同一个 goroutine 里的 recover 捕获**——goroutine A 的 panic 不会跨 goroutine 传播到 goroutine B。

```go
// ✅ 正确：recover 在 defer 直接调用的匿名函数里
defer func() {
    if r := recover(); r != nil {
        fmt.Println("caught:", r)
    }
}()

// ❌ 无效：recover 不在 defer 里
r := recover()  // r 永远是 nil

// ❌ 无效：recover 被 defer 间接调用
defer myRecover()  // myRecover 里面调 recover → nil
// 必须由 defer 语句直接调用的那个函数字面量内部直接调用 recover
```

**总结：** `defer` + `recover` 是 Go 处理 panic 的唯一方式。recover 必须写在 defer 函数里才有效；它捕获 panic 的值，让程序有机会优雅降级而不是直接崩溃。

---

## 易错点

| # | 易错场景 | 为什么 |
|---|---------|--------|
| 1 | `defer` 在函数退出时执行，不在代码块退出时执行 | for 循环里的 defer 要等整个函数返回才执行，可能导致资源大量堆积 |
| 2 | 多个 `defer` 后进先出 | 先打开的资源后关闭——顺序反了可能出 bug |
| 3 | `defer f(val)` 的参数在注册时就求值 | 后面改 `val` 不影响已注册的 defer |
| 4 | `defer f(&val)` 传指针，能看到最新数据 | 地址不变，但地址指向的数据可能改变 |
| 5 | 有名返回值可能被 `defer` 修改 | 闭包捕获的是返回值变量本身，改了影响最终结果 |
| 6 | `recover` 必须在 defer 函数里才有效 | 其它任何位置调用 recover 都返回 nil |
| 7 | panic 只在本 goroutine 内传播 | goroutine 的 panic 必须由自己 recover，别人救不了 |

---

## 快问快答

### Q1：defer 一般用来做什么？

最常见的是资源释放——关闭文件、释放锁、关闭数据库连接、关闭 HTTP response body。也常用于配合 recover 处理 panic，或者在函数退出时打印日志、上报 metrics。

### Q2：多个 defer 的执行顺序是什么？

后进先出（LIFO），最后注册的最先执行。Go 内部用 `_defer` 链表串联，函数退出时从头部开始倒序遍历。

### Q3：`defer fmt.Println(num)` 里的 `num` 什么时候确定？

注册 `defer` 那一刻就确定了。参数在 `defer` 语句执行时就已经被求值、复制、保存好了，之后 `num` 怎么变都和 defer 持有的那个副本没有关系。

### Q4：defer 能修改返回值吗？

能，但仅限于有名返回值。有名返回值时，defer 里改的就是返回值变量本身。无名返回值时，return 先把值复制到编译器的临时变量 `~r0`，然后 defer 改的是局部变量，`~r0` 不受影响。

### Q5：为什么 for 循环里用 defer 要小心？

因为 defer 只在**函数退出时**执行，不是循环退出时。如果在 for 循环里 defer 关闭文件，所有文件会堆积到函数返回时才一起关闭——可能导致文件句柄耗尽。循环里应该直接用 `Close()` 或者把循环体抽成独立函数。

```go
// ❌ 错误：文件在函数结束时才关
for _, name := range files {
    f, _ := os.Open(name)
    defer f.Close()   // 所有文件堆到函数结束才关
}

// ✅ 正确：循环体抽成函数
for _, name := range files {
    func() {
        f, _ := os.Open(name)
        defer f.Close()
        // 处理 f...
    }()  // 每次迭代结束时 f 就关了
}
```

---

## 一句话总结

`defer` 负责把收尾动作放到函数退出前执行——记住后进先出、参数提前求值、有名返回值可被修改、panic 靠 recover 兜底。`return` 三步走（赋值 → defer → 真正返回），有名返回值在 defer 里改的是返回值本身，无名返回值改的是和返回值无关的局部变量。
