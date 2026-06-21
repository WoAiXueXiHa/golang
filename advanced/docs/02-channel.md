# Channel 并发通信

## 这一章要记住什么

- channel 是 goroutine 之间**通信的管道**，Go 的哲学是"**通过通信来共享内存，而不是通过共享内存来通信**"
- **无缓冲 channel 是同步的**——发送和接收必须同时就绪，否则阻塞等待对方
- **有缓冲 channel 能异步**——缓冲区不满时发送不阻塞，缓冲区不空时接收不阻塞；满了或空了就变回同步等待
- 关闭 channel 后**还能读**，先读完剩余数据，再读就是零值；用 `v, ok := <-ch` 或 `for range` 来感知关闭
- 利用有缓冲 channel **满了就阻塞**的特性，`make(chan bool, 1)` 可以实现轻量级的锁

---

## 1. 无缓冲 channel：同步通信

无缓冲 channel 就像两个人在打电话——必须**同时在线**，一方说话另一方听，缺一个就等。

```go
ch := make(chan int)  // 无缓冲：容量为 0
```

```
发送方 goroutine                    接收方 goroutine
     │                                    │
     │  ch <- 42                          │
     │  ─────────── 阻塞等待 ──────→      │  val := <-ch
     │  悬停在这里，直到有人来接收          │  拿到 42，继续执行
     │                                    │
```

**关键规则：**

- **发送方先到** → 发送方阻塞，直到有接收方来取
- **接收方先到** → 接收方阻塞，直到有发送方来发
- **双方同时就绪** → 数据直接传递，谁也不等

```go
// ❌ 死锁：main 里单独发送，没有其他 goroutine 接收
ch <- 45  // 永远等不到接收方 → fatal error: all goroutines are asleep - deadlock!

// ✅ 正确：发送前先启动一个接收 goroutine
go func() {
    val := <-ch
    fmt.Println("接收到：", val)
}()
ch <- 42  // OK，有接收方在等着
```

注意：`go func()` 启动 goroutine 后，虽然代码是串行写下来的，但 goroutine 调度是非确定性的——谁先到谁先等，Go 运行时自动协调。

### 总结一下

无缓冲 channel 是**强同步**的：每次发送都必须有一个接收方在对面等，反之亦然。发送和接收两个 goroutine 在 channel 上"汇合"，数据直接传递，不经过中间缓冲区。

---

## 2. 有缓冲 channel：异步发送 + 同步阻塞

有缓冲 channel 就像快递柜——柜子有格子的时候，放快递的人不用等收件人亲自来拿；格子满了就得等人来取走才能继续放。

### 核心机制：不空不堵，满了才堵

```
make(chan int, 3)  容量 = 3

发送 ch <- A  →  [ A ][   ][   ]  不阻塞 ✓（还有空位）
发送 ch <- B  →  [ A ][ B ][   ]  不阻塞 ✓（还有空位）
发送 ch <- C  →  [ A ][ B ][ C ]  不阻塞 ✓（刚好满）

发送 ch <- D  →  ⛔ 阻塞！缓冲区满，必须等有人取出一个才能放进
                 此时发送方 goroutine 挂起等待

接收 <-ch     →  [ A ][ B ][   ]  取出一个，释放一个位置
                 被阻塞的 ch <- D 被唤醒 → [ B ][ C ][ D ]
```

反过来也一样：

```
缓冲区空 → 接收方阻塞，等发送方放入数据
```

所以有缓冲 channel 的行为可以总结为一句话：

> **缓冲区内是异步的（不阻塞），超出缓冲区边界就变同步（阻塞等待）**

| 操作 | 缓冲区有空间/有数据时 | 缓冲区满/空时 |
|------|----------------------|--------------|
| 发送 `ch <- v` | 不阻塞，数据放入缓冲 | 阻塞，等待接收方取走 |
| 接收 `<-ch` | 不阻塞，从缓冲取数据 | 阻塞，等待发送方放入 |

### 异步模式：缓冲区充当解耦层

```go
ch := make(chan string, 5)
ch <- "hello"  // 不阻塞，放进缓冲区
ch <- ","      // 不阻塞
close(ch)      // 关闭后依然可以读

// 接收方在另一个 goroutine，可以稍后再读
go func() {
    for i := 0; i < 5; i++ {
        v := <-ch      // 前 2 次读到真实数据，后 3 次读到零值 ""
        fmt.Printf("收到:%s\n", v)
    }
}()
```

发送方和接收方**不需要同时就绪**——缓冲区把它们的执行节奏解耦了。发送方塞完就走，接收方有空再来取。

### 总结一下

有缓冲 channel = 异步 + 同步的混合体。缓冲区有空间时发送不阻塞（异步），缓冲区满了发送就阻塞（同步）。你可以用缓冲区来解耦生产者和消费者的执行节奏。缓冲大小决定了能容忍多少"生产超前消费"。

---

## 3. 关闭 channel 后的读取行为

关闭 channel 不是销毁它，而是告诉接收方"**不会再有人往里面写新数据了**"。

```
ch := make(chan string, 5)
ch <- "hello"
ch <- ","
close(ch)

读取过程：
第 1 次 <-ch  →  "hello"  (缓冲区的数据还在)
第 2 次 <-ch  →  ","      (缓冲区的数据还在)
第 3 次 <-ch  →  ""       (空字符串，零值！缓冲区已空)
第 4 次 <-ch  →  ""
...
```

**两个方式正确判断是否读到了真实数据：**

```go
// 方式 1：comma ok 模式
v, ok := <-ch
if ok {
    fmt.Printf("真实数据：%d\n", v)
} else {
    fmt.Printf("channel 已关闭且读完，v 是零值：%d\n", v)
}

// 方式 2：for range（自动感知关闭）
for v := range ch {
    fmt.Printf("v = %s\n", v)  // 读完缓冲区数据自动退出循环
}
// 循环结束 = channel 已关闭且数据已读完
```

`for range` 是最常用的读取模式——有数据就读，关闭就退出，不用手动判断 `ok`。

### 总结一下

关闭 channel 后缓冲区里的数据**还在**，可以先读完。读完后再读就是零值。用 `v, ok := <-ch` 区分"真实数据"和"关闭后的零值"，或者直接用 `for range`，读完自动退出。

---

## 4. 单向 channel：限制方向，防止误用

单向 channel 是一种**类型约束**，用于在函数签名里声明"这个 channel 只能读"或"只能写"，编译器帮你检查。

```go
type RecvChannel = <-chan int   // 只能从 channel 接收
type SendChannel = chan<- int   // 只能向 channel 发送

ch := make(chan int)

var send SendChannel = ch  // ch 被当作只写 channel
send <- 100                // ✅ 可以发送

var recv RecvChannel = ch  // ch 被当作只读 channel
num := <-recv              // ✅ 可以接收
// recv <- 200             // ❌ 编译错误：send to receive-only type
```

**为什么需要单向 channel？**
比如有一个生产者函数，它只需要写，不需要读——限定为 `chan<-` 后，函数内部就不可能误读，调用方看到签名也知道这个函数只负责生产数据。更常见的用途是 Go 标准库里的 `context.Done()` 返回 `<-chan struct{}`——只让你等，不让你往里塞东西。

### 总结一下

单向 channel 是编译器层面的约束：`<-chan` 只读，`chan<-` 只写。在函数参数和返回值上使用，能防止误操作，也让意图更清晰。

---

## 5. 用通信来共享内存

这是 Go 并发哲学最核心的一句话，代码里通过对比 C++ 和 Go 的写法体现出来：

```
C++ 思路（共享内存来通信）          Go 思路（通信来共享内存）

 全局 sum                            局部 sum（每个 goroutine 自己的）
    ↑   ↑                               ↓           ↓
  线程1 线程2                         goroutine1  goroutine2
    ↑   ↑                               ↓           ↓
   加锁保护                            各自算完，通过 channel 发送
                                      ↓           ↓
                                      c <- sum    c <- sum
                                          ↓       ↓
                                      主 goroutine 从 channel 收
```

```go
Sum := func(s []int, c chan int) {
    sum := 0
    for _, v := range s {
        sum += v
    }
    c <- sum  // 把结果发出去，不写全局变量
}

go Sum(s[:len(s)/2], c)   // 算前半段，结果发到 channel
go Sum(s[len(s)/2:], c)   // 算后半段，结果发到 channel

x, y := <-c, <-c          // 从 channel 接收两个结果
fmt.Println(x, y, x+y)    // 汇总
```

**没有全局变量，没有锁，没有共享状态。** 每个 goroutine 操作自己的局部变量，算完通过 channel 把结果传给下一个环节。数据的所有权在 goroutine 之间通过 channel 传递。

### 总结一下

Go 的风格是"不要用共享变量来通信，用通信来传递数据"。每个 goroutine 管好自己的局部数据，结果通过 channel 发送，避免了锁和竞态。

---

## 6. 用有缓冲 channel 实现锁 🔒

### 核心思路

`make(chan bool, 1)` 创建一个**容量为 1** 的 channel：

```
容量 = 1 的 channel 就像一个令牌桶，里面最多只能放一个令牌。

谁拿到了令牌（往里放了 true），谁就能操作临界区；
其他人想放令牌但桶满了 → 阻塞等待；
操作完取出令牌（<-ch），让出位置 → 唤醒下一个等待者。
```

```go
ch := make(chan bool, 1)  // 容量 1，一次只能有一个人拿到"令牌"

add := func(ch chan bool, num *int) {
    ch <- true      // ① 往 channel 放令牌：如果满了就阻塞等待
    *num = *num + 1 // ② 临界区：同一时刻只有一个 goroutine 能执行到这里
    <-ch            // ③ 取出令牌：释放位置，让给下一个等待的 goroutine
}

for i := 0; i < 100; i++ {
    go add(ch, &num)
}
```

### 执行流程分解

```
channel: [ ? ]  (容量1，初始空)

goroutine 1:  ch <- true  →  channel 为空，放入成功  →  [true]  →  执行 *num++
goroutine 2:  ch <- true  →  channel 已满，阻塞等待  →  挂起 ⏸
goroutine 3:  ch <- true  →  channel 已满，阻塞等待  →  挂起 ⏸
...

goroutine 1:  <-ch        →  取出 true，腾出空位      →  [ ] 
                          →  唤醒 goroutine 2

goroutine 2:  ch <- true  →  channel 为空，放入成功  →  [true]  →  执行 *num++
...
```

这样 100 个 goroutine 就是排队串行执行 `*num = *num + 1`，最终 `num` 一定是 100，**不会出现竞态导致结果小于 100**。

### 这和 sync.Mutex 的对应关系

| 概念 | `sync.Mutex` | `make(chan bool, 1)` |
|------|-------------|---------------------|
| 加锁 | `mu.Lock()` | `ch <- true` |
| 解锁 | `mu.Unlock()` | `<-ch` |
| 容量 | 不适用 | 容量 = 1（互斥锁） |
| 阻塞语义 | Lock 在已锁定时阻塞 | send 在 channel 满时阻塞 |

**用 channel 实现锁只是展示 channel 的特性**，理解原理就好。实际写代码时直接使用 `sync.Mutex`——语义更明确，性能更好，`defer mu.Unlock()` 也不容易忘。

### 总结一下

有缓冲 channel 容量为 1 时，`ch <- true` 相当于"获取锁"（满则等），`<-ch` 相当于"释放锁"（取走令牌让位）。本质是把 channel 的**满阻塞**特性当成了互斥锁的排队机制。理解这个能帮你更深入地掌握 channel 的阻塞语义，但生产代码直接用 `sync.Mutex` 更合适。

---

## 易错点

1. **无缓冲 channel 单方操作会死锁**——`ch <- 42` 写在 main goroutine 里又没有另一个 goroutine 来接收，程序直接 panic：`all goroutines are asleep - deadlock!`
2. **向已关闭的 channel 发送数据会 panic**——关闭后只能读不能写；关闭一个已经关闭的 channel 也会 panic
3. **`demo2_after_close_recv` 里循环读 5 次太粗暴**——数据只有 2 个却硬读 5 次，后 3 次都是零值。正确的做法是用 `for range` 或判断 `ok`
4. **代码里 `for range` 注释说"不关闭也正常"**——实际上不关闭 channel 时 `for range` 永远不会退出，goroutine 会一直阻塞在 range 上，直到 main 退出才被杀掉。虽然这个 demo 里因为 main 退出被清理了，但如果 main 不退出就是 goroutine 泄漏
5. **`demo7_mutex` 用 `time.Sleep(2 * time.Second)` 等 goroutine 完成**——这只是演示写法，生产代码必须用 `sync.WaitGroup` 确保所有 goroutine 都执行完成

---

## 快问快答

### Q1：无缓冲 channel 和有缓冲 channel 的根本区别是什么？

答：无缓冲 channel 每次收发都要求双方同时就绪，是**完全同步**的——发送方和接收方在 channel 上汇合，数据直接传递。有缓冲 channel 里，缓冲区内有空间时发送不阻塞（异步），缓冲区满时发送阻塞（同步）；缓冲区有数据时接收不阻塞，空时接收阻塞。简单说：无缓冲 = 全同步，有缓冲 = 半异步半同步。

### Q2：有缓冲 channel 的同步模式和异步模式怎么切换？

答：不需要手动切换，它**自动**切换。缓冲区不空不堵时工作在异步模式（发送/接收都不阻塞），缓冲区满了/空了就自动切到同步模式（阻塞等待对方）。缓冲容量本身就决定了"允许异步处理多少数据"。

### Q3：`make(chan bool, 1)` 怎么就能当锁用了？

答：容量为 1 意味着 channel 最多存一个值。`ch <- true` 往里放——如果 channel 为空，放入成功，相当于拿到锁；如果已经有人放了、channel 满了，就阻塞等待，相当于等锁。`<-ch` 取出来腾空间，相当于释放锁，下一个等待的 goroutine 就能放进去了。本质上就是把 channel 的"满了阻塞"当成了互斥机制。

### Q4：channel 关闭后能做什么，不能做什么？

答：能读——先把缓冲区里剩余的数据读完，之后每次读返回零值。能用 `v, ok := <-ch` 或 `for range` 感知到关闭。**不能写**——往关闭的 channel 写数据会 panic。**不能重复关**——`close` 一个已关闭的 channel 也会 panic。

### Q5：什么时候用 `for range` 读 channel？

答：当你明确知道有**一个发送方**会在合适的时候 `close` channel，就用 `for range`。它会在 channel 关闭且数据读完后自动退出循环，不用手动判断。如果发送方永远不关闭 channel，`for range` 就永远不会退出，接收方 goroutine 会一直阻塞。

### Q6："通过通信来共享内存"具体是什么意思？

答：就是每个 goroutine 操作自己的局部变量，不搞全局共享变量。算完结果通过 channel 发给下一个 goroutine，数据的所有权跟着 channel 消息走。这样就不需要锁、不需要考虑竞态——因为同一时刻只有拿到数据的那个 goroutine 在操作它。

---

## Channel 特性总结

| 特性 | 说明 |
|------|------|
| **类型安全** | `chan int` 只能传 int，`chan string` 只能传 string，编译时检查 |
| **并发安全** | 多个 goroutine 同时读写同一个 channel，Go 运行时保证不会数据竞争 |
| **阻塞语义** | 发送/接收在特定条件下自动阻塞（无缓冲总是阻塞；有缓冲满/空时阻塞） |
| **关闭语义** | 关闭后只能读不能写；读完后返回零值；`for range` 自动感知关闭 |
| **方向约束** | `<-chan` 只读，`chan<-` 只写，编译器强制执行 |
| **容量固定** | 缓冲大小在 `make` 时就确定了，不能动态扩缩 |
| **FIFO 顺序** | channel 内部是先进先出队列，先发的数据先被接收 |
| **零值 nil** | 未初始化的 channel 是 `nil`，对它读写都会**永久阻塞**（不是 panic） |

---

## 一句话总结

channel 是 goroutine 之间安全通信的管道：无缓冲是同步握手机制，有缓冲是异步队列 + 满/空时自动同步；关闭后可读不可写；容量为 1 的有缓冲 channel 能当轻量锁用；Go 的并发哲学是"用通信来传递数据的所有权，而不是用共享变量加锁来通信"。
