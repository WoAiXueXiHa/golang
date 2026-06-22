# select：channel 多路复用

## 这一章要记住什么

- `select` 是 Go 的 **channel 多路复用语句**——同时监听多个 channel，哪个就绪就执行哪个
- **没有 default 的 select 会阻塞**——直到某个 case 可以执行；所有 case 都无法执行且没有 default → 死锁
- **有 default 的 select 是非阻塞的**——所有 case 都不就绪时立刻走 default，不会等
- **多个 case 同时就绪时，Go 随机选一个执行**——不是按代码顺序，是伪随机均匀分布
- `select {}`（空 select）**永久阻塞**——常用于 main 函数里阻止程序退出
- Go 的 select 和 Linux 的 select 不是一回事——前者管 channel，后者管文件描述符

---

## 1. select 的语法和基本概念

`select` 长得像 `switch`，但它的 case 不是值比较，而是 **channel 操作**：

```go
select {
case <-ch1:              // 从 ch1 读成功 → 执行
    fmt.Println("from ch1")
case ch2 <- 2:           // 向 ch2 写成功 → 执行
    fmt.Println("to ch2")
default:                 // 所有 case 都不就绪 → 执行
    fmt.Println("nothing ready")
}
```

```
select 的决策过程：

                    ┌─────────────┐
                    │ select 开始  │
                    └──────┬──────┘
                           ↓
              ┌────────────────────────┐
              │ 检查所有 case，有就绪的吗？ │
              └──────────┬─────────────┘
                     ↙        ↘
                   有            没有
                   ↓              ↓
           ┌──────────────┐  ┌──────────────┐
           │ 只有一个就绪？ │  │ 有 default 吗？ │
           └──┬────────┬──┘  └──┬───────┬───┘
           是         否       有        没有
           ↓          ↓        ↓         ↓
      执行那个   随机选一个  执行default  阻塞等待
       case      执行                       ↓
                                    ┌──────────────┐
                                    │ 一直等到有    │
                                    │ case 就绪     │
                                    │ （如果没有    │
                                    │  其他goroutine │
                                    │  能让case就绪  │
                                    │  → 死锁）     │
                                    └──────────────┘
```

Go 的 `select` 和 Linux 的 `select()` 系统调用虽然都叫 select，区别很大：

| 维度 | Go select | Linux select |
|------|-----------|-------------|
| 监视对象 | channel | 文件描述符（socket、文件等） |
| 就绪语义 | channel 可读/可写 | fd 可读/可写/异常 |
| 阻塞方式 | goroutine 阻塞 | 内核阻塞调用线程 |
| 层级 | 语言级语句 | OS 系统调用 |

### 总结一下

`select` 是 Go 在语言层面提供的 channel 多路复用——同时监听多个 channel，哪个先就绪就先执行哪个。有 default 时不阻塞，没 default 时阻塞到某个 case 就绪。多个 case 同时就绪则随机选一个。

---

## 2. 阻塞式 select：没有 default 时的行为

### 有 goroutine 在写 → 阻塞等到数据来

```go
ch := make(chan int)

go func() {
    time.Sleep(1 * time.Second)
    ch <- 42                // 1 秒后往 ch 发数据
}()

select {
case v := <-ch:             // 阻塞 1 秒，等到数据后执行
    fmt.Println(v)          // 42
}
```

select 会阻塞在这里，直到子 goroutine 往 ch 里写了数据。

### 没有 goroutine 能让 case 就绪 → 死锁

代码里的 `demo1`：

```go
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)

select {
case <-ch1:                 // ch1 没数据，且没有 goroutine 会往里写
    fmt.Println("from ch1")
case num := <-ch2:          // ch2 也没数据，也没有 goroutine 会往里写
    fmt.Printf("num is %d\n", num)
}
// fatal error: all goroutines are asleep - deadlock!
```

```
执行流程：

select 检查 case：
  case <-ch1 → ch1 空，且没有发送方 → 不可就绪
  case <-ch2 → ch2 空，且没有发送方 → 不可就绪
  
没有 default → 当前 goroutine 阻塞，等待 case 就绪

Go 运行时检测：所有 goroutine 都在睡、没有任何 goroutine 能唤醒对方
  → panic: deadlock
```

**死锁的根本原因**：select 在等 ch1 或 ch2 可读，但没有 goroutine 会往这两个 channel 写数据——等不到的东西，永远等不到。

### 空 select：永久阻塞

代码注释里提到了 `空 select 永久阻塞`：

```go
select {}  // 永远阻塞，不会 panic
```

这和上面不同——上面是有 case 但 case 不可执行（死锁），这里是根本没有 case。Go 把空 select 视为"我愿意永久阻塞"，不会触发死锁检测。常见用法是放在 `main` 函数最后，阻止程序退出：

```go
func main() {
    // 启动各种 goroutine...
    select {}  // 主 goroutine 永久阻塞，程序不会退出
}
```

### 总结一下

没有 default 的 select 会阻塞直到某个 case 就绪。如果所有 case 都无法执行，且没有其他 goroutine 能改变这个状态，Go 运行时会检测到死锁并 panic。空 `select {}` 是合法用法，永久阻塞且不触发死锁。

---

## 3. 非阻塞式 select：有 default

代码里的 `demo2` 展示了最简单的非阻塞 select：

```go
ch := make(chan int, 1)   // 空的，没数据

select {
case <-ch:                 // ch 没数据 → 不满足
    fmt.Println("from ch")
default:                   // 立刻走这里
    fmt.Println("default....")
}
```

```
检查 case <-ch → 读不了
  → 有 default → 直接执行 default，不等待
```

**default 的作用就是把 select 从阻塞变成非阻塞**——"有就处理，没有拉倒"。

`demo3` 展示了更实际的用法——在循环里用 select 做轮询：

```go
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)

go func() {
    time.Sleep(1 * time.Second)
    for i := 0; i < 3; i++ {
        select {
        case v := <-ch1:
            fmt.Println("received from ch1: ", v)
        case v := <-ch2:
            fmt.Println("received from ch2: ", v)
        default:
            fmt.Println("default....")
        }
    }
}()

ch1 <- 1                      // 往 ch1 发数据
time.Sleep(1 * time.Second)
ch2 <- 2                      // 往 ch2 发数据
time.Sleep(1 * time.Second)
```

```
主 goroutine              子 goroutine
─────────────             ─────────────
                          Sleep(1s)  ← 先睡 1 秒
ch1 <- 1  (写入成功)
Sleep(1s)
                          醒来，进入 for 循环，i=0：
                            select → ch1 有数据 → 读走，打印 "received from ch1: 1"
                          i=1：
                            select → ch1 空，ch2 空 → default: "default...."
                          
ch2 <- 2  (写入成功)         i=2：
Sleep(1s)                     select → ch2 有数据 → 读走，打印 "received from ch2: 2"
                          循环结束
```

**注意**：两个 channel 都用了容量 1 的缓冲——这样主 goroutine 发送时不会阻塞等待接收方。如果换成无缓冲 channel，`ch1 <- 1` 会在那里等子 goroutine 来接收。

### 总结一下

select 加了 default 就变成非阻塞——所有 case 都不就绪时直接走 default，不等待。适合"有数据就处理、没数据做别的事"的场景。多个 case + default 可以同时监听多个 channel，哪个有数据就处理哪个。

---

## 4. 多 case 就绪时的随机选择

代码里的 `demo4` 展示了 select 的一个关键行为：**当多个 case 同时就绪时，Go 伪随机选一个执行**。

```go
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)

ch1 <- 66          // ch1 有数据
ch2 <- 11          // ch2 也有数据

select {
case v := <-ch1:                // ch1 可读 ✅
    fmt.Println("from ch1: ", v)
case v := <-ch2:                // ch2 也可读 ✅
    fmt.Println("from ch2: ", v)
}
```

```
两个 case 都就绪的状态：

  ch1 ┌─────────┐    select
      │   66    │◄───────── case1: <-ch1  ✅
      └─────────┘
                    随机选一个 →
      ┌─────────┐
  ch2 │   11    │◄───────── case2: <-ch2  ✅
      └─────────┘

多次运行输出不固定：
  第 1 次 → "from ch1: 66"
  第 2 次 → "from ch2: 11"
  第 3 次 → "from ch1: 66"
  ...
```

**为什么是随机而不是按代码顺序？** 如果总按代码顺序选第一个就绪的 case，后面的 case 可能永远得不到执行——这就是饥饿。Go 用伪随机选择保证公平：每个就绪的 case 都有同等机会被执行。

### 总结一下

多个 case 同时就绪时，Go 运行时随机选一个执行——不是按代码顺序。这是为了防止前面的 case 一直抢到执行导致后面的 case 饿死。你的代码不应该依赖 case 的先后顺序。

---

## 5. select 的经典使用模式

你的代码覆盖了 select 的基础用法，这里扩展几个在实际项目中几乎一定会用到的模式。

### 模式一：超时控制

```go
select {
case v := <-ch:
    fmt.Println("got:", v)
case <-time.After(3 * time.Second):
    fmt.Println("timeout!")
}
```

```
time.After 返回一个 channel，在指定时间后会收到一个值
  → 如果 ch 在 3 秒内没数据，timeout case 就绪
  → 不会永远等下去
```

### 模式二：非阻塞发送

```go
select {
case ch <- value:
    fmt.Println("sent")
default:
    fmt.Println("channel full, dropped")
}
```

发送方不想被阻塞——channel 满了就丢掉或走其他逻辑。常用于有缓冲 channel 的"尽力发送"场景。

### 模式三：done channel 优雅退出

```go
select {
case v := <-workCh:
    // 处理工作
case <-doneCh:           // 收到退出信号
    // 清理资源、return
}
```

配合 context 或手动 close 的 channel，让 goroutine 能响应外部退出信号。

### 总结一下

select 不止是"等多 channel"——配合 `time.After` 做超时、配合 default 做非阻塞、配合 done channel 做优雅退出，是 Go 并发编程里的瑞士军刀。

---

## 易错点

1. **所有 case 都不就绪 + 没 default → 死锁**——代码 `demo1` 里两个 channel 都没数据、没有 goroutine 会写，select 永远等不到，Go 运行时检测到死锁 panic
2. **把 select 的 case 顺序当成执行顺序**——多个 case 就绪时是随机选的，不要依赖 case 的先后顺序来保证逻辑
3. **带 default 的 select 在循环里可能空转**——如果循环里 select 大部分时候走 default，CPU 会空转浪费。需要在 default 里加点 sleep 或用其他机制限制频率
4. **`time.After` 在循环里用会内存泄漏**——每次 `time.After` 创建新的 timer，select 结束后如果没有超时，timer 要等到超时才会被 GC。循环里应该用 `time.NewTimer` + `Reset` 复用
5. **有缓冲 channel 发送时不阻塞**——代码里 `ch1 <- 1` 因为通道容量为 1 所以不阻塞。如果换成无缓冲 channel，主 goroutine 会卡在发送上等接收，这会影响理解 `demo3` 的执行时序

---

## 快问快答

### Q1：select 和 switch 有什么区别？

答：switch 是值匹配——拿一个值和多个 case 比，哪个相等执行哪个。select 是 channel 就绪判断——检查多个 channel 操作，哪个能执行就执行哪个。switch 里的 case 是顺序比较的；select 里的 case 如果多个都就绪，是随机选的。

### Q2：有 default 和没 default 的 select，核心区别是什么？

答：没 default 的 select 是阻塞的——必须等到某个 case 就绪才继续。有 default 的 select 是非阻塞的——所有 case 都不就绪时立刻走 default。前者是"等"，后者是"看一眼就走"。

### Q3：为什么 select 多个 case 就绪时要随机选？

答：防止饥饿。如果总是按顺序选第一个就绪的 case，排在后面的 case（比如低优先级的 channel）可能永远得不到执行。随机选择给每个就绪的 case 同等的执行机会。

### Q4：空 select `select {}` 有什么用？

答：永久阻塞当前 goroutine。最常见的是放在 main 函数末尾，让程序在后台 goroutine 运行时不退出。和死锁不同——空 select 没有 case，Go 认为你是故意阻塞，不会触发死锁检测。

### Q5：Go 的 select 和 Linux 的 select 有什么关系？

答：名字都叫 select，但本质上没关系。Linux select 是系统调用，监视文件描述符（socket、文件等）的可读/可写状态。Go select 是语言语句，监视 channel 操作的可行性。Go 的运行时在底层可能用类似的 I/O 多路复用机制（epoll 等）来实现网络相关的 channel 操作，但语言层面的 select 是独立的抽象。

### Q6：select 里 case 的 channel 操作会阻塞吗？

答：select 的 case 本身不会阻塞——如果 channel 操作不能立刻完成（比如向满的 channel 写、从空的 channel 读），这个 case 就是"不就绪"，select 会跳过它（去其他 case 或 default），或者整个 select 阻塞等待。select 阻塞的是整个语句，不是单个 case。

---

## select 速查

| 写法 | 含义 | 行为 |
|------|------|------|
| `select { case <-ch: }` | 等 ch 可读 | 阻塞到有数据 |
| `select { case ch <- v: }` | 等 ch 可写 | 阻塞到能写入 |
| `select { case <-ch: default: }` | 非阻塞读 | 有数据就读，没数据走 default |
| `select { case <-ch1: case <-ch2: }` | 等多 channel | 阻塞到任意一个可读，同时可读则随机 |
| `select { case <-ch: case <-time.After(d): }` | 超时等待 | ch 有数据或超时，哪个先来执行哪个 |
| `select {}` | 永久阻塞 | 不会 panic，常用于 main 阻止退出 |

---

## 一句话总结

`select` 是 Go 在语言层面提供的 channel 多路复用——同时监听多个 channel，哪个就绪就执行哪个；没 default 时阻塞等待，有 default 时非阻塞；多个 case 同时就绪时随机选；配合超时和 done channel，能写出既简洁又强大的并发控制逻辑。
