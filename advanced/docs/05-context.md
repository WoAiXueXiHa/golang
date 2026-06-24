# Context 上下文控制

## 这一章要记住什么

- Context 是 goroutine 之间的**控制信号管道**，解决了"一个请求里启动多个 goroutine，超时或失败时怎么让它们全部停掉"的问题。
- Context 是一棵**单向不可变树**：父节点取消，所有子节点自动取消；但子节点取消不影响父节点和兄弟节点。
- 取消信号通过**关闭 channel** 广播：`ctx.Done()` 返回的 channel 被关闭后，所有监听它的 goroutine 立刻收到零值，一个都不漏。
- 四种挂载能力：`WithCancel`（手动取消）、`WithDeadline`（定时点取消）、`WithTimeout`（定时长取消）、`WithValue`（挂请求范围数据）。
- Context 要作为函数的**第一个参数**逐层传递，不要存到 struct 里长期持有；拿到 cancel 后必须 `defer cancel()` 释放资源。

---

## 1. 为什么需要 Context —— 如果没有它会怎样

假设你在写一个 HTTP 服务：每个请求进来，你启动 3 个 goroutine 分别去查数据库、调下游 API、写缓存。如果客户端断开连接了（比如用户关了浏览器），这 3 个 goroutine 还在跑，白白消耗资源。

**没有 Context 的时候**，你得自己造轮子：

```go
stop := make(chan struct{})
go worker1(stop)
go worker2(stop)
go worker3(stop)

// 想停掉它们？
close(stop)  // 得自己管理这个 channel，自己约定"关了就是停"
```

但这个方案问题很多：你怎么统一处理超时？怎么给不同的 goroutine 设不同的截止时间？怎么在调用链上传一些数据（比如 trace id），而不改动每一层函数的签名？

Context 把这些问题标准化了。

```
一个典型的 Web 请求：

请求进入
   │
   ├── goroutine A: 查数据库（需要 3 秒超时）
   ├── goroutine B: 调下游 API（需要 5 秒超时）
   └── goroutine C: 写缓存（请求结束时要停掉）

如果请求超时（总时限 2 秒），Context 自动通知 A、B、C 全部退出。
```

### 总结一下

Context 就是把"让一群 goroutine 统一停掉"这个需求标准化了 —— 不用自己造 channel 约定，不用每个函数都加一个 `stop chan` 参数。

---

## 2. Context 到底是什么 —— 接口定义

`context.Context` 是一个接口，定义在标准库 `context` 包里。它只有 4 个方法：

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)  // 返回截止时间，没有则 ok=false
    Done() <-chan struct{}                     // 返回一个 channel，取消时关闭
    Err() error                                // 返回取消原因：Canceled 或 DeadlineExceeded
    Value(key interface{}) interface{}         // 取值
}
```

每个方法都小，但组合起来刚好够用：

### 2.1 Done() —— 取消信号的载体

```text
正常运行时：
  ctx.Done() → 返回一个 channel，这个 channel 一直阻塞，没人往里面写东西

取消后：
  ctx.Done() → 返回的 channel 被关闭
  → 所有在 <-ctx.Done() 上等待的 goroutine 全部被唤醒
  → 因为是从已关闭的 channel 读取，立刻返回零值（不阻塞）
```

**为什么用"关闭 channel"而不是"往 channel 里写一个值"？**

因为关闭一个 channel 是**广播**：所有阻塞在这个 channel 上的 goroutine 全部被唤醒。如果只是往里写一个值，只有一个 goroutine 能收到。

```
写值方式（错误）：
  主 goroutine → ch <- struct{}{} → 只有一个 goroutine 收到，其它的收不到

关闭方式（正确）：
  主 goroutine → close(ch) → goroutine1 ✓收到
                            → goroutine2 ✓收到
                            → goroutine3 ✓收到
```

### 2.2 Deadline() —— 这个 context 什么时候到期

```go
deadline, ok := ctx.Deadline()
if ok {
    fmt.Println("还剩", time.Until(deadline))
}
```

如果 context 是通过 `WithDeadline` 或 `WithTimeout` 创建的，`Deadline()` 返回那个截止时间。如果是 `WithCancel` 或 `Background()` 创建的，`ok` 就是 `false`。

### 2.3 Err() —— 为什么被取消了

```go
select {
case <-ctx.Done():
    err := ctx.Err()  // context.Canceled 或 context.DeadlineExceeded
}
```

只有两种错误值：
- `context.Canceled`：手动调了 `cancel()`
- `context.DeadlineExceeded`：到了 deadline 自动取消

### 2.4 Value() —— 从 context 里取数据

```go
traceID := ctx.Value(traceKey).(string)
```

返回 `interface{}`，需要自己断言类型。后面会详细讲。

### 总结一下

Context 接口就四个方法，但设计得很克制：`Done()` 负责"通知你该停了"，`Deadline()` 告诉你"什么时候停"，`Err()` 告诉你"为什么停了"，`Value()` 负责"顺路带点数据"。各司其职，不越界。

---

## 3. Context 的树状结构 —— 为什么取消会传播

每个 context 都是通过**包装父 context** 产生的。这形成了一棵树：

```text
              Background()               ← 根节点，永远不会取消
                   │
        ┌──────────┼──────────┐
        │          │          │
   WithCancel  WithTimeout  WithValue   ← 各自挂载不同能力
        │          │
    ┌───┴───┐  WithCancel(子)
    │       │       │
  gor1    gor2    gor3                 ← 叶节点：goroutine 监听它们
```

**传播规则：**

```text
父节点取消 → 所有后代节点的 Done() channel 全部关闭
子节点取消 → 父节点不受影响
兄弟节点取消 → 互不影响

例子：
  root := Background()
  parent, cancelParent := WithCancel(root)
  child1, cancel1 := WithCancel(parent)
  child2, _ := WithTimeout(parent, 5*time.Second)

  cancel1()   → 只取消 child1 和它的子节点，parent 和 child2 不受影响
  cancelParent() → 取消 parent、child1、child2，root 不受影响
```

用图表示：

```text
调用 cancelParent() 之前：
  root          ○（活跃）
  parent        ○（活跃）
  child1        ○（活跃）
  child2        ○（活跃）


调用 cancelParent() 之后：
  root          ○（活跃，不受影响）
  parent        ✕（已取消）
  child1        ✕（被父节点牵连取消）
  child2        ✕（被父节点牵连取消）

调用 cancel1() 之后（不调 cancelParent）：
  root          ○（活跃）
  parent        ○（活跃，不受子节点影响）
  child1        ✕（已取消）
  child2        ○（活跃，兄弟节点取消不影响它）
```

### 总结一下

Context 的取消信号是**从上往下单向传播**的。这棵树里，父节点是"控制者"，子节点是"被控制者"。父取消，子孙全灭；子取消，父毫不知情。

---

## 4. 四种创建方式 —— 代码详解

### 4.1 根节点：Background() 和 TODO()

```go
ctx := context.Background()  // 最常用，作为整个 context 树的根
ctx := context.TODO()        // 占位符，表示"还没想好怎么用"
```

两者返回的都是同一个类型 `context.emptyCtx`，永远不会取消、没有 deadline、不存值。

它们的区别**只在语义上**：`Background()` 说"这就是根"，`TODO()` 说"这里以后要改"。

```text
使用场景：

Background() → main() 函数、请求入口、测试入口
TODO()       → 重构时暂时不确定用什么 context，先占个位
```

**实际用的时候，很少直接传 `Background()` 到业务函数里，而是在入口处基于它创建一个有实际能力的 context。**

### 4.2 WithCancel —— 手动取消

```go
ctx, cancel := context.WithCancel(parent)
// ...
cancel()  // 调用后 ctx.Done() channel 关闭
```

**创建过程拆解：**

```text
parent (Background)
   │
   │  WithCancel(parent)
   │  内部做了什么：
   │  1. 创建一个新的 cancelCtx 结构体，记录父节点是谁
   │  2. 创建一个新的 channel（Done channel）
   │  3. 如果父节点已经取消了，直接关闭新 channel
   │  4. 如果父节点还没取消，把自己注册到父节点的"子节点列表"里
   │
   ▼
child (cancelCtx)
  - Done channel: 新建的 chan struct{}
  - cancel 函数: 关闭 Done channel + 从父节点注销 + 递归取消所有子节点
```

**重要：cancel 函数有几个特性**

1. **可以多次调用**——只有第一次生效，后续调用是空操作
2. **并发安全**——可以被多个 goroutine 同时调用
3. **取消是级联的**——调用 cancel 后，这个节点下的所有子节点全部取消

**案例分析：**

```go
ctx, cancel := context.WithCancel(context.Background())
go Watch(ctx, "goroutine1")
go Watch(ctx, "goroutine2")
time.Sleep(6 * time.Second)
cancel()  // ← 这里！两个 Watch 同时收到信号退出
time.Sleep(1 * time.Second)  // 留 1 秒让 goroutine 执行完退出打印
```

时序图：

```text
时间  主 goroutine                goroutine1              goroutine2
0s    创建 ctx + cancel
      启动 goroutine1             watching...
      启动 goroutine2                                     watching...
1s                                watching...             watching...
2s                                watching...             watching...
3s                                watching...             watching...
4s                                watching...             watching...
5s                                watching...             watching...
6s    调用 cancel()
      → ctx.Done() channel 关闭
                                  <-ctx.Done() 返回       <-ctx.Done() 返回
                                  "goroutine1 exit!"      "goroutine2 exit!"
7s    程序退出
```

### 4.3 WithDeadline —— 到点自动取消

```go
deadline := time.Now().Add(4 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

**和 WithCancel 的关系：**

`WithDeadline` 内部是基于 `WithCancel` 实现的。它多了一个定时器：

```text
WithDeadline 内部做的事情：
1. 基于父 context 创建一个 cancelCtx（和 WithCancel 一样）
2. 启动一个后台 goroutine，用 time.AfterFunc 设一个定时器
3. 定时器到了 → 自动调用 cancel()
4. 同时把 cancel 函数返回给用户（用户也可以提前手动取消）
```

**案例分析：**

```go
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
defer cancel()
go Watch(ctx, "goroutine1")
go Watch(ctx, "goroutine2")
time.Sleep(6 * time.Second)
```

时序图：

```text
时间  主 goroutine                goroutine1              goroutine2
0s    创建 ctx（4 秒后自动取消）
      启动 goroutine1             watching...
      启动 goroutine2                                     watching...
1s                                watching...             watching...
2s                                watching...             watching...
3s                                watching...             watching...
4s    ← ctx 自动到期！
      ctx.Done() channel 关闭
                                  <-ctx.Done() 返回       <-ctx.Done() 返回
                                  "goroutine1 exit!"      "goroutine2 exit!"
5s                                已退出                   已退出（空闲）
6s    "end watching!!!" 打印
```

注意：4 秒后 goroutine 退出了，但主 goroutine 要睡到 6 秒才结束。这期间两个子 goroutine 已经不再运行了。

### 4.4 WithTimeout —— WithDeadline 的快捷方式

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()
```

它和 `WithDeadline` 的区别**仅仅是参数形式**：

```go
// 这两行等价：
ctx, cancel = context.WithDeadline(parent, time.Now().Add(1*time.Second))
ctx, cancel = context.WithTimeout (parent, 1*time.Second)

// WithTimeout 的源码就是一行：
//   return WithDeadline(parent, time.Now().Add(timeout))
```

**案例分析（1 秒超时）：**

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()
```

和 `demo2` 用 4 秒不同，这里只有 1 秒。所以两个 goroutine 只打印一次 "watching..." 就被超时取消了。

```text
时间  主 goroutine                goroutine1              goroutine2
0s    创建 ctx（1 秒后自动取消）
1s    ← ctx 超时！
      ctx.Done() 关闭
                                  exit!                   exit!
```

### 4.5 WithValue —— 挂数据

```go
ctx := context.WithValue(context.Background(), "name", "vect")
name := ctx.Value("name").(string)  // "vect"
```

**底层数据结构：**

```text
WithValue 就像一个链表节点：

Background()  →  valueCtx{A}  →  valueCtx{B}  →  valueCtx{C}
  (空)             key=k1         key=k2         key=k3
                   val=v1         val=v2         val=v3

ctx.Value(k2) 的查找过程：
  1. 先在 valueCtx{C} 里找 → 没有 k2
  2. 往父节点 valueCtx{B} 里找 → 找到了！返回 v2
  3. 如果一直找到 Background() 都没找到，返回 nil
```

**查找方向是从当前节点往上（往父节点方向），不是往下。** 这意味着：
- 子节点可以"覆盖"父节点的值（因为先在当前节点找）
- 父节点看不到子节点的值（因为只往上找）

```go
parent := context.WithValue(ctx, "key", "parent")
child  := context.WithValue(parent, "key", "child")

parent.Value("key")  // → "parent"
child.Value("key")   // → "child"（覆盖了父节点的值）
```

**案例分析：**

```go
ctx := context.WithValue(context.Background(), "name", "vect")
fmt.Printf("name is %s\n", ctx.Value("name").(string))  // → "vect"
```

### 总结一下

四种创建方式都是包装父 context 生成子 context。`WithCancel` 是基础——后面两个（`WithDeadline`、`WithTimeout`）内部都基于它加了定时器。`WithValue` 是独立的一类，用链表存键值对，查找从当前节点往上回溯。

---

## 5. 代码逐块解析

### 5.1 Watch 函数 —— context 消费者的标准写法

```go
Watch := func(ctx context.Context, name string) {
    for {
        select {
        case <-ctx.Done():       // 收到取消信号
            fmt.Printf("%s exit!\n", name)
            return
        default:                 // 还没取消，继续工作
            fmt.Printf("%s watching...\n", name)
            time.Sleep(time.Second)
        }
    }
}
```

这是 context 消费者最经典的模式：**`select` + `ctx.Done()` + `default`**。

```text
select 的分支逻辑：

每次循环：
  ┌─ ctx.Done() channel 被关闭了？
  │   YES → 执行 case <-ctx.Done()：打印 exit，return
  │   NO  → 走 default：打印 watching...，睡 1 秒
  └─ 下一轮循环
```

**这里的 `select` 不是阻塞的**，因为 `default` 分支的存在。如果去掉 `default`：

```go
select {
case <-ctx.Done():
    return
}
// 没有 default → select 会阻塞在这里，直到 ctx.Done() 关闭
// 这样就不能持续打印 "watching..." 了
```

### 5.2 demo1_WithCancel —— 手动控制取消时机

```go
ctx, cancel := context.WithCancel(context.Background())
go Watch(ctx, "goroutine1")
go Watch(ctx, "goroutine2")
time.Sleep(6 * time.Second)
fmt.Println("end watching!!!")
cancel()
time.Sleep(time.Second)
```

**执行流程逐行解释：**

1. `WithCancel(Background())` → 创建一个可以手动取消的 context，同时拿到 `cancel` 函数
2. `go Watch(...)` × 2 → 启动两个 goroutine，把**同一个 ctx** 传进去。两个 goroutine 监听的是**同一个 Done channel**
3. `time.Sleep(6s)` → 主 goroutine 等 6 秒，让两个 Watch 各打印 6 次 "watching..."
4. `cancel()` → **核心操作**：关闭 ctx 内部的 Done channel。两个 Watch 的 `<-ctx.Done()` 同时返回
5. 两个 Watch 分别打印 "exit!" 并 return
6. 最后 `time.Sleep(1s)` → 给 goroutine 1 秒时间执行完退出逻辑

### 5.3 demo2_WithDeadline —— 定时自动取消

```go
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
defer cancel()
```

这里设置 4 秒后自动取消。`defer cancel()` 的作用：

```text
正常路径（没有提前返回）：
  4 秒后自动取消 → 函数结束时 defer cancel() 再调一次（无影响，cancel 幂等）

异常路径（函数提前返回）：
  函数因为错误提前 return → defer cancel() 立刻执行 → 不需要等到 4 秒
```

**如果没有 `defer cancel()`：** 如果函数提前返回了（比如某个初始化出错），context 要等到 4 秒后才被垃圾回收，这 4 秒内它持有的定时器 goroutine 一直占着资源。

### 5.4 demo3_WithTimeout —— 演示短超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
```

只有 1 秒，所以两个 goroutine 大概只能打印 1 次 "watching..." 就被超时掐掉了。

**为什么主 goroutine 睡 6 秒？**

因为主 goroutine 不知道自己启动的子 goroutine 什么时候退出。这里用 `Sleep(6s)` 只是一种粗糙的等待方式。**更正确的做法**是用 `sync.WaitGroup` 或 channel 来同步，但那是并发控制的范畴了。

### 5.5 demo4 —— WithValue 传递数据

```go
func1 := func(ctx context.Context) {
    fmt.Printf("name is %s\n", ctx.Value("name").(string))
}
ctx := context.WithValue(context.Background(), "name", "vect")
go func1(ctx)
time.Sleep(time.Second)
```

这里 `.(string)` 是类型断言。`Value()` 返回 `interface{}`，必须断言成具体类型才能用。

**代码里的小问题：**

这个函数没有取消机制。如果 `func1` 里不只是打印，而是个长任务，就没办法通过 context 让它停下来。实际工程里，通常把 `WithValue` 和 `WithCancel`/`WithTimeout` 组合使用：

```go
// 更完整的写法
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)  // 先挂超时
defer cancel()
ctx = context.WithValue(ctx, "name", "vect")             // 再挂数据
go func1(ctx)  // 既有超时保护，又能拿到数据
```

---

## 6. Context 使用规范

### 6.1 永远是第一个参数

```go
// ✅ 正确
func DoSomething(ctx context.Context, arg string) error

// ❌ 错误：ctx 不在第一个
func DoSomething(arg string, ctx context.Context) error
```

官方建议参数名就叫 `ctx`，不要用别的名字。

### 6.2 不要存到 struct 里

```go
// ❌ 错误
type Worker struct {
    ctx context.Context  // 不要这样
}

// ✅ 正确：每次调用时传入
func (w *Worker) Run(ctx context.Context) error
```

原因：context 的生命周期通常是**一次请求的范围**。存到 struct 里意味着你把它和 struct 的生命周期绑定了，这不符合 context 的设计意图。

### 6.3 不要传 nil

如果你不确定传什么 context，用 `context.TODO()`，不要传 `nil`。传 `nil` 会在调用 `ctx.Done()` 等地方触发 panic。

### 6.4 Value 的 key 用自定义类型

```go
// ❌ 不好：string 可能被别的包覆盖
ctx = context.WithValue(ctx, "traceID", traceID)

// ✅ 好：自定义类型不会冲突
type contextKey string
const traceKey contextKey = "traceID"
ctx = context.WithValue(ctx, traceKey, traceID)
```

Go 的 `interface{}` 比较用的是 `==`，自定义类型不会和其他包的 key 相等，避免了命名冲突。

### 6.5 不要用 context.Value 传业务参数

```go
// ❌ 滥用
ctx = context.WithValue(ctx, "userID", 123)
ctx = context.WithValue(ctx, "pageSize", 20)

// ✅ 正确用法：只传横切关注点的数据
ctx = context.WithValue(ctx, traceKey, "abc-123")   // trace id
ctx = context.WithValue(ctx, requestIDKey, "xyz")   // request id
```

业务参数应该直接通过函数参数传递，这样类型安全、编译器能检查、代码也可读。

---

## 易错点

1. **忘记 defer cancel()**。`WithCancel`、`WithDeadline`、`WithTimeout` 返回的 `cancel` 函数必须调用。即使超时后会自动取消，`defer cancel()` 也不能省——它负责释放 context 在父节点那里注册的信息。不调会一直挂在父节点的子节点列表里，造成内存泄漏。

2. **拿到 cancel 不调用，等 GC 回收**。有人以为函数返回后 context 就会被 GC 回收。但实际上，context 内部启动的定时器 goroutine（`WithDeadline`/`WithTimeout` 用的 `time.AfterFunc`）会阻止它被回收，直到定时器到期。

3. **在 select 的 Done 分支里做耗时操作**。Done channel 关闭后，应该尽快退出。如果在这个分支里做同步 IO 或复杂计算，goroutine 不能及时释放，context 的取消效果就打折扣了。

4. **Value 的 key 用 string 导致冲突**。不同的包都可能用 `"name"` 或 `"id"` 做 key，互相覆盖。用自定义类型可以避免。

5. **把 Context 存到 struct 里**。Context 应该是函数参数，传递的是"本次调用的上下文"，不是对象的属性。如果你想在 struct 初始化时缓存一个 context，大概率是设计上想歪了。

6. **在主 goroutine 里 select ctx.Done() 但不给 default**。如果主 goroutine 自己阻塞在 `<-ctx.Done()` 上，那谁来调 `cancel()`？这种写法会导致死锁。

---

## 快问快答

### Q1：Context 是什么？一句话说清楚。

答：Context 是 Go 里用来在 goroutine 之间传递取消信号、截止时间和请求范围数据的标准机制。它是一棵树，父节点取消，子节点全停。

### Q2：ctx.Done() 为什么返回 channel 而不是直接返回一个 bool？

答：因为 channel 可以同时被多个 goroutine 监听。关闭一个 channel 会产生广播效应——所有阻塞在这个 channel 上的 goroutine 同时被唤醒。如果只是返回一个 bool，你只能在某个 goroutine 里轮询它，效率低而且不可靠。

### Q3：WithDeadline 和 WithTimeout 的区别是什么？

答：`WithDeadline` 传入绝对时间点（"下午 3:00 截止"），`WithTimeout` 传入相对时长（"3 秒后截止"）。`WithTimeout` 内部就是 `WithDeadline(parent, time.Now().Add(timeout))`，两者本质一样，只是参数形式不同。

### Q4：为什么要 defer cancel() 而不是直接 cancel()？

答：两个原因。一是 `cancel()` 负责把当前 context 从父节点的子节点列表里移除，不调用会内存泄漏。二是如果函数有多个 return 路径（提前出错返回），`defer` 保证不管走哪条路径 `cancel()` 都会被执行。

### Q5：context.WithValue 的查找方向是什么？

答：从当前节点**往父节点方向**回溯，一层一层往上找。当前节点找不到就找父节点，一直找到根。这意味着子节点可以"遮蔽"父节点同 key 的值，但父节点看不到子节点的值。

### Q6：一个 context 可以被取消两次吗？

答：可以调用 `cancel()` 多次，但只有第一次生效。`cancel` 函数内部用的是 `sync.Once` 或等效机制，重复调用不会 panic，也不会重复关闭 channel。

---

## 一句话总结

Context 是 goroutine 的生命周期遥控器：`Background()` 是根，往上挂 `WithCancel`（手动停）、`WithDeadline`/`WithTimeout`（定时停）、`WithValue`（顺路带数据），然后作为第一个参数传下去，所有子 goroutine 通过 `<-ctx.Done()` 统一响应退出信号。
