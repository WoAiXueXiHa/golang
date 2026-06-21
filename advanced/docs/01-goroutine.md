# Goroutine 并发基础

## 这一章要记住什么

- `go` 关键字启动协程，但 **main 协程退出会杀死所有子协程**，等待要用 `sync.WaitGroup`
- goroutine 调度是**非确定性**的，每次运行顺序可能不同，不要依赖执行顺序
- Go 1.22 起，`for` 循环变量**每次迭代创建新副本**，闭包捕获不再共享同一变量
- 每个 goroutine 初始栈约 **2KB**，可以轻松创建上万个，远轻于 OS 线程（~1MB）
- `WaitGroup` 传参必须传**指针**，否则每个协程操作的是独立副本，`Wait()` 永远等不到归零

---

## 1. goroutine 是什么

Go 里启动一个并发执行流，不需要像传统多线程那样手动创建线程，而是用 `go` 关键字：

```go
go fmt.Println("Hello from goroutine")
```

goroutine 是**用户态的轻量执行流**，由 Go 运行时调度到 OS 线程上执行，不是直接对应一个线程。它的初始栈空间只有约 2KB，而 OS 线程通常是 1MB，所以 goroutine 的创建和切换开销非常小。

```
程序（磁盘上的静态文件）
  └── 进程（加载到内存，资源分配的容器）
        └── 线程（内核执行流，CPU 调度的实体）
              └── goroutine（用户态执行流，Go 运行时调度）
```

**调度关系：** 多个 goroutine 被 Go 运行时的调度器（GMP 模型）映射到少量 OS 线程上执行，goroutine 的阻塞不会导致线程阻塞，运行时会把其他 goroutine 切换到空闲线程。

### 总结一下

goroutine 是 Go 的并发执行单元，用户态调度，2KB 起步栈空间。用 `go` 关键字就能启动，不直接对应 OS 线程，创建和切换成本远低于线程。

---

## 2. 主协程退出，子协程立即被杀

```go
go func() {
    time.Sleep(500 * time.Millisecond)
    fmt.Println("这句话大概率不会被打印....")
}()

fmt.Println("主协程结束 -> 子协程被杀死...")
```

Go 的 `main()` 函数本身运行在主 goroutine 中。当 `main()` 返回时，**整个程序退出，所有还在运行的子 goroutine 会立刻被杀死，不会有任何通知或清理机会**。

这点跟很多语言的线程模型不同——没有"等待所有线程结束"的默认行为，必须显式同步。

### 总结一下

主 goroutine 退出 → 程序结束 → 所有子 goroutine 直接被杀，不等待、不通知。想让子 goroutine 跑完，必须显式等它。

---

## 3. WaitGroup：等待一组协程完成

`sync.WaitGroup` 就像一个计数器，用来等待一组 goroutine 全部完成：

```
初始状态: counter = 0
  
wg.Add(2)           counter = 2     ← 登记 2 个任务
  go worker(A)         goroutine A 启动
  go worker(B)         goroutine B 启动

                      goroutine A 完成 → wg.Done() → counter = 1
                      goroutine B 完成 → wg.Done() → counter = 0

wg.Wait() 解除阻塞 ← counter == 0
```

三个核心方法：

```go
var wg sync.WaitGroup

wg.Add(n)   // 计数器 +n，在启动 goroutine 之前调用
wg.Done()   // 计数器 -1，goroutine 内部用 defer 调用
wg.Wait()   // 阻塞直到计数器归零
```

**关键：传指针，不传值。**

```go
// ✅ 正确：传指针
go worker("worker A", &wg)

// ❌ 错误：传值，每个 goroutine 拿到的是 WaitGroup 的副本
// 副本的 Done() 不影响原来的计数器，Wait() 永远等不到归零 → 死锁
go worker("worker A", wg)
```

`WaitGroup` 结构体内部有 `noCopy` 标记，用 `go vet` 能检测到值拷贝的问题，但编译器不会直接报错。

### 总结一下

`WaitGroup` 是计数器：`Add` 登记任务数，`Done` 完成一个，`Wait` 等全部完成。**必须传指针**，值拷贝会导致死锁。

---

## 4. 闭包捕获循环变量（Go 1.22 前后）

这是 Go 并发里最经典的坑，但 Go 1.22 之后有重大变化。

### Go 1.22 之前：所有协程共享同一个变量

```go
for i := 1; i <= 5; i++ {
    go func() {
        fmt.Printf("协程 [%d]\n", i)  // 全都引用同一个 i
    }()
}
```

```
循环变量 i：内存里只有一份
    │
    ├── goroutine 1 ── 引用 ──→ i
    ├── goroutine 2 ── 引用 ──→ i
    ├── goroutine 3 ── 引用 ──→ i
    ├── goroutine 4 ── 引用 ──→ i
    └── goroutine 5 ── 引用 ──→ i

循环结束 i=6 → 所有 goroutine 读到 6（或读到中间值，完全看调度）
```

### Go 1.22 之后：每次迭代创建新变量

**Go 1.22（2024.02）改变了循环语义：每次迭代的 `i` 是一个全新的变量。**

```
迭代 1: i₁ = 1 → goroutine 1 捕获 i₁
迭代 2: i₂ = 2 → goroutine 2 捕获 i₂   ← 和 i₁ 不是同一块内存！
迭代 3: i₃ = 3 → goroutine 3 捕获 i₃
迭代 4: i₄ = 4 → goroutine 4 捕获 i₄
迭代 5: i₅ = 5 → goroutine 5 捕获 i₅
```

每个 goroutine 会打印自己迭代的值（1~5），不会全部一样。**但执行顺序仍然不确定**——调度是非确定性的，可能先打印 5 再打印 1。

为了**向后兼容旧版本**和代码清晰，有这两种写法：

```go
// 方式 1：循环内创建局部变量
for i := 1; i <= 5; i++ {
    i := i      // 遮蔽外层 i，创建新变量
    go func() {
        fmt.Printf("协程 [%d]\n", i)
    }()
}

// 方式 2：通过参数显式传入
for i := 1; i <= 5; i++ {
    go func(id int) {
        fmt.Printf("协程 [%d]\n", id)
    }(i)   // 把 i 的值拷贝给参数 id
}
```

### 总结一下

Go 1.22 之前，闭包捕获循环变量是所有协程共享同一个变量地址，典型 bug。Go 1.22 之后每次迭代创建新变量，闭包各自捕获不同的变量。但为了兼容旧版本和代码可读性，参数传入或 `i := i` 仍然是推荐写法。

---

## 5. goroutine 的交错执行

goroutine 的执行顺序是**非确定性**的——每次运行结果可能不同。

```go
go printer("a", 5)  // 打印 a1 a2 a3 a4 a5
go printer("b", 5)  // 打印 b1 b2 b3 b4 b5
go printer("c", 5)  // 打印 c1 c2 c3 c4 c5
```

可能的输出：`a1 b1 c1 a2 b2 c2 ...`，也可能 `a1 a2 b1 c1 b2 a3 ...`，完全取决于 Go 调度器当时的决策。

**`time.Sleep` 会让出时间片：**

```go
// 有 Sleep：每次循环主动让出，三个 goroutine 交错
time.Sleep(time.Millisecond)  // → a1 b1 c1 a2 b2 c2 ...

// 无 Sleep：一个 goroutine 可能连续跑完整轮
// → X1 X2 X3 X4 X5 ... Y1 Y2 Y3 Y4 Y5 ...
```

没有 Sleep 的情况下，如果任务很短，调度器可能没有机会在中间切换——一个 goroutine 一口气跑完，另一个才开始。

### 总结一下

goroutine 调度是非确定性的，**不要依赖执行顺序**。`time.Sleep` 等阻塞操作会让出执行权，给其他 goroutine 执行机会；短任务没有阻塞点时，调度器可能不会在中间切换。

---

## 6. goroutine 的轻量级特性

```go
const N = 10000
for i := 0; i < N; i++ {
    go func() {
        time.Sleep(10 * time.Millisecond)
    }()
}
```

创建 1 万个 goroutine 是瞬间完成的：

|          | goroutine | OS 线程 |
|----------|-----------|---------|
| 初始栈   | ~2KB      | ~1MB    |
| 1 万个   | ~20MB     | ~10GB   |
| 创建成本 | 极低      | 高      |
| 切换成本 | 用户态    | 内核态  |

如果这 1 万个任务串行执行，需要 `10000 × 10ms = 100s`；而并发执行只需要约 10ms（所有 goroutine 同时 sleep）。这就是 goroutine 的威力——**大量 IO 等待任务并发时，时间几乎只取决于最慢的那个任务**。

### 总结一下

goroutine 栈空间极小（2KB）、创建极快，可以轻松开上万个。对于 IO 密集型任务（如网络请求），并发效果显著，整体耗时接近最慢的那个任务而不是所有任务之和。

---

## 易错点

1. **`WaitGroup` 传值不传指针**——`go worker(&wg)` 必须传 `&wg`，传值会导致每个 goroutine 操作副本，`Wait()` 死锁
2. **`wg.Add()` 放在 goroutine 内部**——应该在外层调用 `wg.Add()`，在 goroutine 内部 `defer wg.Done()`；如果 `Add` 也写在 goroutine 里，`Wait()` 可能在 `Add` 之前就执行了
3. **Go 1.22 后闭包捕获不再是坑，但不能依赖版本**——如果代码可能跑在 <1.22 的环境，仍然需要 `i := i` 或参数传入
4. **不要用 `time.Sleep` 做同步**——`demo1_basic` 里用 `time.Sleep(10 * time.Millisecond)` 等 goroutine，这只是演示，生产代码里必须用 `WaitGroup` 或 channel 来同步

---

## 快问快答

### Q1：`go` 关键字做了什么？

答：`go f()` 会创建一个新的 goroutine，把 `f()` 放到这个 goroutine 里执行。调用方不会阻塞，立即继续往下走。goroutine 由 Go 运行时调度到 OS 线程上执行。

### Q2：主 goroutine 退出会发生什么？

答：整个程序退出，所有还在运行的子 goroutine 立刻被杀掉，没有任何通知或清理的机会。所以必须用 `WaitGroup` 或 channel 等机制让主 goroutine 等待。

### Q3：`sync.WaitGroup` 为什么要传指针？

答：`WaitGroup` 是一个结构体，传值会拷贝整个结构体。goroutine 里调 `Done()` 操作的是副本的计数器，不影响原版的计数器，原版的 `Wait()` 永远等不到归零，造成死锁。

### Q4：Go 1.22 之后闭包捕获循环变量还有什么问题吗？

答：Go 1.22 修复了循环变量共享的问题——每次迭代创建新变量，闭包各自捕获自己的那份。所以不会再出现"所有 goroutine 打印同一个值"的 bug。但如果代码要兼容 <1.22 的版本，仍然建议用参数传入或 `i := i`。

### Q5：为什么 goroutine 能开上万个？

答：每个 goroutine 初始栈只有大约 2KB，1 万个才约 20MB 内存。而一个 OS 线程的初始栈约 1MB，1 万个线程就是 10GB，根本开不起。goroutine 的创建和切换都在用户态完成，不需要陷入内核。

### Q6：goroutine 的执行顺序可以预测吗？

答：不能。goroutine 的调度是非确定性的，每次运行结果都可能不同。不要依赖任何特定的执行顺序，需要顺序控制时用同步原语（channel、Mutex、WaitGroup 等）。

---

## 一句话总结

goroutine 是 Go 并发的核心——轻量、便宜、调度由运行时负责；用 `WaitGroup` 等子协程完成；闭包捕获循环变量在 Go 1.22 后不再是坑但兼容写法不亏；永远不要依赖 goroutine 的执行顺序。
