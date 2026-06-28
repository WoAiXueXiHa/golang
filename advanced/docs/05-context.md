# Context 上下文控制

## 这一章要记住什么

- Context 是 goroutine 之间的**控制信号管道**——一个请求里启动了多个 goroutine，超时或失败时，Context 负责通知所有 goroutine 统一退出
- Context 是一棵**单向不可变树**：父节点取消 → 所有子节点自动取消；子节点取消 → 父节点和兄弟节点不受影响
- 取消信号通过**关闭 channel 实现广播**：`ctx.Done()` 返回的 channel 被关闭后，所有监听它的 goroutine 同时收到零值，一个不漏
- 四种挂载能力：`WithCancel`（手动取消）、`WithDeadline`（定时点取消）、`WithTimeout`（定时长取消）、`WithValue`（挂请求范围数据）
- Context 作为函数的**第一个参数**逐层传递，不要存到 struct 里；拿到 cancel 后必须 `defer cancel()` 防止内存泄漏
- `cancel` 函数**幂等且并发安全**——多次调用只有第一次生效，不会 panic

---

## 0. 没有 Context 的世界：为什么要发明它

假设一个 HTTP 请求进来，你启动了 3 个 goroutine 分别查数据库、调下游 API、写缓存。如果客户端断开连接，这些 goroutine 还在跑——白白消耗 CPU、内存、数据库连接。

**自己造轮子：**

```go
stop := make(chan struct{})
go worker1(stop)
go worker2(stop)
go worker3(stop)

close(stop)  // 想停掉它们？得自己管理 stop channel，自己约定"关了就是停"
```

这个方案的问题：

| 问题 | 说明 |
|------|------|
| 超时怎么处理？ | 得自己写定时器 + channel 组合逻辑 |
| 不同 goroutine 不同超时？ | 得自己创建多套 stop + timer |
| 链条上有数据要传递？ | 得往每个函数签名里加参数 |
| goroutine 的子 goroutine？ | stop channel 得一层层手动传 |

**Context 把这些标准化了：**

```
一个典型的 Web 请求：

请求进入
   │
   ├── goroutine A: 查数据库（ctx 传递 3 秒超时）
   ├── goroutine B: 调下游 API（ctx 传递 5 秒超时）
   └── goroutine C: 写缓存（ctx 传递请求结束信号）

如果请求超时（总时限 2 秒），Context 自动通知 A、B、C 全部退出。
```

---

## 1. Context 接口定义：四个方法拆解

`context.Context` 是一个接口，只有 4 个方法：

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

### 1.1 Done() — 取消信号的载体

**函数原型：**

```go
Done() <-chan struct{}
```

| 项目 | 说明 |
|------|------|
| 参数 | 无 |
| 返回值 | `<-chan struct{}` — **只读** channel。未取消时阻塞，取消后立即返回零值 |
| 为什么是只读 | 调用方只能**听**信号，不能**发**信号——只有创建者手里的 `cancel()` 才能触发 |

**行为对比：**

```
正常运行时：
  ch := ctx.Done()  → 返回的 channel 开着，<-ch 会阻塞

取消后：
  ch := ctx.Done()  → channel 被关闭
  <-ch               → 立刻返回 struct{}{} 零值（从已关闭 channel 读，不阻塞）
```

**为什么用"关闭 channel"而不是"往 channel 里写一个值"？**

关闭 channel 的效果是**广播**——所有阻塞在这个 channel 上的 goroutine 全被唤醒。写值只能唤醒一个。

```
写值（错误思路）：
  主 goroutine → ch <- val → 只有一个 goroutine 收到

关闭（正确实现）：
  主 goroutine → close(ch) → goroutine 1 ✓ 收到
                            → goroutine 2 ✓ 收到
                            → goroutine 3 ✓ 收到
```

### 1.2 Deadline() — 还有多少时间

**函数原型：**

```go
Deadline() (deadline time.Time, ok bool)
```

| 项目 | 说明 |
|------|------|
| 参数 | 无 |
| 返回值 `deadline` | 绝对截止时间点（如 `2026-06-28 15:30:00`） |
| 返回值 `ok` | `true` = 设了截止时间；`false` = 没设（WithCancel 或 Background 创建的） |

```go
deadline, ok := ctx.Deadline()
if ok {
    fmt.Println("还剩", time.Until(deadline))
} else {
    fmt.Println("没有截止时间")
}
```

只有 `WithDeadline` / `WithTimeout` 创建的 context 才会返回 `ok=true`。

### 1.3 Err() — 为什么被取消了

**函数原型：**

```go
Err() error
```

| 项目 | 说明 |
|------|------|
| 参数 | 无 |
| 返回值 | 取消原因：`nil`（未取消）、`context.Canceled`、`context.DeadlineExceeded` |

两种标准错误值：

```go
// 手动调了 cancel()
var Canceled = errors.New("context canceled")

// 到了 deadline / timeout 自动取消
var DeadlineExceeded error = deadlineExceededError{}
```

典型用法：

```go
select {
case <-ctx.Done():
    switch ctx.Err() {
    case context.Canceled:
        fmt.Println("手动取消")
    case context.DeadlineExceeded:
        fmt.Println("超时了")
    }
}
```

**重要：`ctx.Err()` 在取消前返回 `nil`，取消后返回非 nil 错误。** 不能拿它来"提前预判"会不会取消——它是事后诊断。

### 1.4 Value() — 从 context 里取数据

**函数原型：**

```go
Value(key any) any
```

| 项目 | 说明 |
|------|------|
| 参数 `key` | 查找键，**必须可比较**（`==` 能比较的类型） |
| 返回值 | 查到的值（`any` 类型，需要自己类型断言）；没找到返回 `nil` |

后文 5.5 节会详细展开。

### 小结

四个方法各司其职：`Done()` 负责"通知你停"，`Deadline()` 告诉你"什么时候停"，`Err()` 告诉你"为什么停了"，`Value()` 负责"顺路带点数据"。接口很小，但组合起来刚好够用。

---

## 2. Context 的树状结构：取消为什么能传播

### 2.1 创建方式一览

| 函数 | 完整签名 | 父节点要求 | 说明 |
|------|---------|-----------|------|
| `Background()` | `func Background() Context` | 无 | 根节点，永远不会取消 |
| `TODO()` | `func TODO() Context` | 无 | 占位符，语义同 Background |
| `WithCancel` | `func WithCancel(parent Context) (ctx Context, cancel CancelFunc)` | 非 nil | 创建可手动取消的子节点 |
| `WithDeadline` | `func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)` | 非 nil | 创建定时点自动取消的子节点 |
| `WithTimeout` | `func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)` | 非 nil | 创建定时长自动取消的子节点 |
| `WithValue` | `func WithValue(parent Context, key, val any) Context` | 非 nil | 创建带键值对的子节点 |
| `WithCancelCause` | `func WithCancelCause(parent Context) (ctx Context, cancel CancelCauseFunc)` |非 nil | Go 1.20+ 取消时附带原因 |
| `AfterFunc` | `func AfterFunc(ctx Context, f func()) (stop func() bool)` | 非 nil | Go 1.21+ 取消后自动执行回调 |

`CancelFunc` 类型定义：

```go
type CancelFunc func()  // 调用后关闭 Done channel，级联取消所有子节点
```

`CancelCauseFunc` 类型定义（Go 1.20+）：

```go
type CancelCauseFunc func(cause error)  // 取消时附带原因，ctx.Err() 返回该原因
```

### 2.2 树的结构

每个 context 都是通过**包装父 context** 产生的，形成一棵树：

```
              Background()               ← 根节点，永不取消
                   │
        ┌──────────┼──────────┐
        │          │          │
   WithCancel  WithTimeout  WithValue   ← 各自挂载不同能力
        │          │
    ┌───┴───┐  WithCancel(子)
    │       │       │
  gor1    gor2    gor3                 ← 叶节点：goroutine 监听它们
```

### 2.3 取消传播规则

```
父节点取消 → 所有子孙节点的 Done() channel 全部关闭
子节点取消 → 父节点不受影响
兄弟节点取消 → 互不影响
```

**具体例子：**

```go
root := context.Background()
parent, cancelParent := context.WithCancel(root)
child1, cancel1 := context.WithCancel(parent)
child2, _ := context.WithTimeout(parent, 5*time.Second)

cancel1()         // → 只取消 child1 + 它的子节点，parent / child2 不受影响
cancelParent()    // → 取消 parent、child1、child2（全局洗白），root 不受影响
```

**图示：**

```
调用 cancelParent() 之前：        调用 cancelParent() 之后：
  root      ○ (活跃)               root      ○ (活跃，不受影响)
  parent    ○ (活跃)               parent    ✕ (已取消)
  child1    ○ (活跃)               child1    ✕ (被牵连)
  child2    ○ (活跃)               child2    ✕ (被牵连)

调用 cancel1()（不调 cancelParent）：
  root      ○ (活跃)
  parent    ○ (活跃，不受影响)
  child1    ✕ (已取消)
  child2    ○ (活跃，兄弟互不影响)
```

### 小结

Context 取消信号是**从上往下单向传播**的。父节点是控制者，子节点是被控制者。父取消，子孙全灭；子取消，父毫不知情。

---

## 3. 四种创建方式：代码详解

### 3.1 根节点：Background() 和 TODO()

**函数原型：**

```go
func Background() Context
func TODO() Context
```

| 项目 | `Background()` | `TODO()` |
|------|---------------|----------|
| 参数 | 无 | 无 |
| 返回值 | `Context` | `Context` |
| 底层类型 | `emptyCtx`（int 类型别名） | `emptyCtx`（同一个类型） |
| 能否取消 | 永不取消 | 永不取消 |
| Done() | 返回 nil（取不到值，永远阻塞） | 同 |
| Deadline() | 返回 ok=false | 同 |

两者底层是**同一个类型**，区别**只在语义上**：

```go
// Background：这就是根，请求的起点
ctx := context.Background()

// TODO：这里以后要改成正经的 context，先占个位
ctx := context.TODO()
```

使用场景：

```
Background() → main()、请求入口、测试入口、goroutine 树的起点
TODO()       → 重构时暂不确定用什么 context；第三方库的过渡期占位
```

### 3.2 WithCancel —— 手动取消

**函数原型：**

```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
```

| 项目 | 说明 |
|------|------|
| `parent` | 父 context，**不能为 nil**（传 nil 会 panic） |
| 返回值 `ctx` | 新建的子 context，Done() 返回的 channel 在 cancel 调用后关闭 |
| 返回值 `cancel` | 取消函数，调用后关闭 ctx 的 Done channel，级联取消所有子节点 |

**cancel 函数的特性：**

| 特性 | 说明 |
|------|------|
| 幂等性 | 多次调用只有第一次生效 |
| 并发安全 | 可以多个 goroutine 同时调用 |
| 级联取消 | 调用后该节点下的所有子孙节点全部取消 |
| 必须调用 | 即使不用，也要 `defer cancel()` 释放父节点中的注册信息 |

**代码案例（对应 demo1_WithCancel）：**

```go
Watch := func(ctx context.Context, name string) {
    for {
        select {
        case <-ctx.Done():              // 收到取消信号
            fmt.Printf("%s exit!\n", name)
            return
        default:                         // 还没取消，继续干活
            fmt.Printf("%s watching...\n", name)
            time.Sleep(time.Second)
        }
    }
}

ctx, cancel := context.WithCancel(context.Background())
go Watch(ctx, "goroutine1")
go Watch(ctx, "goroutine2")

time.Sleep(6 * time.Second)
fmt.Println("end watching!!!")
cancel()                                // 关闭 Done channel → 两个 Watch 同时退出
time.Sleep(time.Second)                 // 等 goroutine 执行完退出打印
```

```
时间轴：
0s   创建 ctx + cancel，启动 goroutine1、goroutine2
1s   goroutine1 watching...    goroutine2 watching...
2s   goroutine1 watching...    goroutine2 watching...
...   
5s   goroutine1 watching...    goroutine2 watching...
6s   主 goroutine: cancel() ──→ Done channel 关闭
     goroutine1: <-ctx.Done() 返回 → "goroutine1 exit!" → return
     goroutine2: <-ctx.Done() 返回 → "goroutine2 exit!" → return
7s   程序退出
```

**内部创建过程：**

```
parent (Background)
   │
   │  WithCancel(parent)
   │  内部做了什么：
   │  1. new(cancelCtx)，记录父节点
   │  2. 创建新的 Done channel（chan struct{}）
   │  3. 如果父节点已经取消 → 立刻关闭新 channel
   │  4. 如果父节点活跃 → 把自己注册到父节点的 children map 里
   │
   ▼
child (cancelCtx)
  - Done channel: 新创建的 chan struct{}
  - cancel 函数: 关闭 Done channel → 从父节点注销 → 递归取消所有子节点
```

### 3.3 WithDeadline —— 到点自动取消

**函数原型：**

```go
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)
```

| 项目 | 说明 |
|------|------|
| `parent` | 父 context，**不能为 nil** |
| `d` | **绝对时间点**（如 `time.Now().Add(4*time.Second)`），到了就自动取消 |
| 返回值 `ctx` | 带截止时间的子 context |
| 返回值 `cancel` | 手动取消函数（提前取消用）；必须 `defer cancel()` |

**注意：** 如果 `d` 已经是过去的时间（`d.Before(time.Now())`），创建的 context **立刻就是已取消状态**。

**WithDeadline 和 WithCancel 的关系：**

`WithDeadline` 内部是基于 `WithCancel` 实现的，多了一个定时器：

```
WithDeadline 内部做的事情：
  1. 先调用 WithCancel(parent) 创建基础 cancelCtx
  2. 用 time.AfterFunc 设一个定时器
  3. 定时器到期 → 自动调用 cancel()
  4. 同时把 cancel 返回给用户（用户也可以提前手动调）
```

**代码案例（对应 demo2_WithDeadline）：**

```go
ctx, cancel := context.WithDeadline(
    context.Background(),
    time.Now().Add(4*time.Second),  // 4 秒后自动取消
)
defer cancel()

go Watch(ctx, "goroutine1")
go Watch(ctx, "goroutine2")

time.Sleep(6 * time.Second)
fmt.Println("end watching!!!")
```

```
时间轴：
0s   创建 ctx（deadline = now + 4s）
1s   goroutine1 watching...    goroutine2 watching...
2s   goroutine1 watching...    goroutine2 watching...
3s   goroutine1 watching...    goroutine2 watching...
4s   ← deadline 到达，定时器自动调用 cancel()
     goroutine1 exit!          goroutine2 exit!
5s   (两个 goroutine 已退出)
6s   主 goroutine: "end watching!!!"
```

**`defer cancel()` 的作用：**

```
正常路径（没有提前返回）：
  4 秒后自动取消 → 函数结束时 defer cancel() 再调一次（无影响，cancel 幂等）

异常路径（函数提前返回）：
  函数因错误提前 return → defer cancel() 立刻执行 → 不用等到 deadline
  如果不 defer，ctx 一直挂在父节点上 → 内存泄漏
```

### 3.4 WithTimeout —— WithDeadline 的快捷方式

**函数原型：**

```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

| 项目 | 说明 |
|------|------|
| `parent` | 父 context，**不能为 nil** |
| `timeout` | **相对时长**（如 `1*time.Second`），从现在起多久后自动取消 |
| 返回值 | 同 WithDeadline |

**源码就是一行：**

```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
    return WithDeadline(parent, time.Now().Add(timeout))
}
```

所以这两行完全等价：

```go
ctx, cancel = context.WithDeadline(parent, time.Now().Add(1*time.Second))
ctx, cancel = context.WithTimeout (parent, 1*time.Second)
```

**选用原则：**

- 知道**具体的截止时间点**（"下午 3 点整停止"） → `WithDeadline`
- 知道**从现往起多久**（"3 秒后停止"） → `WithTimeout`

**代码案例（对应 demo3_WithTimeout）：**

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

go Watch(ctx, "goroutine1")
go Watch(ctx, "goroutine2")

time.Sleep(6 * time.Second)
fmt.Println("end watching!!!")
```

和 demo2 的区别是超时只有 1 秒，所以两个 goroutine 大概各打印 1 次 "watching..." 就被掐掉了。

```
时间轴：
0s   创建 ctx（1 秒后超时）
1s   ← 超时！Done channel 关闭
     goroutine1 exit!          goroutine2 exit!
2~6s (主 goroutine 还在睡，但子 goroutine 已退出)
6s   "end watching!!!"
```

### 3.5 WithValue —— 在 Context 上挂数据

**函数原型：**

```go
func WithValue(parent Context, key, val any) Context
```

| 项目 | 说明 |
|------|------|
| `parent` | 父 context，**不能为 nil** |
| `key` | 查找键，**必须可比较**（`==` 可比），通常用自定义类型避免冲突 |
| `val` | 要存储的值，类型为 `any` |
| 返回值 | 新的 context，不返回 cancel 函数（WithValue 不涉及取消） |

**底层数据结构：**

WithValue 就像链表节点，每个节点存一对 key-value：

```
Background()  →  valueCtx{k1, v1}  →  valueCtx{k2, v2}  →  valueCtx{k3, v3}
  (空)              ↑                    ↑                    ↑
                 当前节点             当前节点             当前节点
```

**Value() 的查找过程（沿着链表往上回溯）：**

```go
ctx3.Value(k2)
  → 先在 valueCtx{k3, v3} 里找 → 没有 k2
  → 往父节点 valueCtx{k2, v2} 里找 → 找到了！返回 v2
  → 如果一直找到 Background() 都没有 → 返回 nil
```

查找方向是从**当前节点往父节点方向**。这意味着：

```go
parent := context.WithValue(ctx, "key", "parent")
child  := context.WithValue(parent, "key", "child")

parent.Value("key")  // → "parent"
child.Value("key")   // → "child"（覆盖了父节点的同名 key）
// 子节点覆盖父节点，但父节点看不到子节点的值
```

**代码案例（对应 demo4）：**

```go
func1 := func(ctx context.Context) {
    fmt.Printf("name is %s\n", ctx.Value("name").(string))
}

ctx := context.WithValue(context.Background(), "name", "vect")
go func1(ctx)
time.Sleep(time.Second)
```

**案例的问题和改进：**

demo4 里 `func1` 没有取消机制——如果它是长任务，没办法通过 context 让它停。实际工程里，通常组合多个 context：

```go
// ✅ 更完整的写法：取消 + 数据一起挂
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)  // 先挂超时保护
defer cancel()
ctx = context.WithValue(ctx, "name", "vect")             // 再挂数据
go func1(ctx)  // 既有超时保护，又能拿到数据
```

### 小结

| 函数 | 返回 cancel？ | 取消触发方式 | 用途 |
|------|-------------|------------|------|
| `WithCancel` | 是 | 手动调用 `cancel()` | 精确控制取消时机 |
| `WithDeadline` | 是 | 到达绝对时间点 | "下午 3 点截止" |
| `WithTimeout` | 是 | 经过相对时长 | "3 秒后超时" |
| `WithValue` | 否 | 不涉及取消 | 传递请求范围数据 |

四种方法都是包装父 context、返回子 context，从不修改父节点。

---

## 4. Context 消费者标准写法

### 4.1 经典模式：select + Done + default

代码里的 Watch 函数展示了 context 消费者最经典的模式：

```go
Watch := func(ctx context.Context, name string) {
    for {
        select {
        case <-ctx.Done():          // 收到取消信号 → 退出
            fmt.Printf("%s exit!\n", name)
            return
        default:                     // 还没取消 → 继续干活
            fmt.Printf("%s watching...\n", name)
            time.Sleep(time.Second)
        }
    }
}
```

```
每次循环：
  ┌─ ctx.Done() channel 被关闭了？
  │   YES → 执行 case <-ctx.Done()：打印 exit，return
  │   NO  → 走 default：打印 watching...，睡 1 秒
  └─ 下一轮循环
```

**这里的 `select` 不是阻塞的**——因为有 `default` 分支。如果去掉 `default`：

```go
select {
case <-ctx.Done():
    return
}
// 没有 default → select 阻塞，直到 ctx.Done() 关闭
// 不能持续打印 "watching..." 了
```

### 4.2 其他常见模式

**阻塞等待取消：**

```go
// 模式 A：只等取消，不做别的
<-ctx.Done()
return ctx.Err()
```

**带超时的子操作：**

```go
// 模式 B：子操作有自己的超时（继承父 ctx，再加限制）
childCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
defer cancel()
result, err := doSomething(childCtx)  // 父取消或子超时，任一个触发就停
```

**周期性检查：**

```go
// 模式 C：干活前先看一眼"还要不要继续"
for _, item := range items {
    select {
    case <-ctx.Done():
        return ctx.Err()   // 被取消了，剩下的不做了
    default:
    }
    process(item)          // 还没取消，继续做
}
```

---

## 5. Context 使用规范

### 5.1 永远是第一个参数

```go
// ✅ 正确
func DoSomething(ctx context.Context, arg string) error

// ❌ 错误：ctx 不在第一个
func DoSomething(arg string, ctx context.Context) error
```

官方约定参数名就叫 `ctx`，不要用 `c`、`context` 等其他名字。

### 5.2 不要存到 struct 里

```go
// ❌ 错误
type Worker struct {
    ctx context.Context  // 不要这样存
}

// ✅ 正确：每次调用时传入
func (w *Worker) Run(ctx context.Context) error
```

原因：context 的生命周期是**一次请求的范围**，存到 struct 里意味着把 context 和 struct 的生命周期绑定了——比如 struct 是全局单例，一个请求取消了它，另一个请求的 context 也被连带，这就乱套了。

### 5.3 不要传 nil

如果你不确定传什么，用 `context.TODO()` 占位，不要传 `nil`。传 `nil` 在调用 `ctx.Done()` 等方法时可能 panic。

### 5.4 Value 的 key 用自定义类型，不要用 string

```go
// ❌ 不好：string 可能被别的包覆盖
ctx = context.WithValue(ctx, "traceID", traceID)

// ✅ 好：自定义类型不会冲突
type contextKey string
const traceKey contextKey = "traceID"
ctx = context.WithValue(ctx, traceKey, traceID)
```

原因：Go 的 `any` 比较用的是 `==`，不同包定义的 `contextKey` 即使底层都是 `string`，也不相等——避免了命名冲突。

### 5.5 不要用 Value 传业务参数

```go
// ❌ 滥用
ctx = context.WithValue(ctx, "userID", 123)
ctx = context.WithValue(ctx, "pageSize", 20)

// ✅ 正确：只传横切关注点
ctx = context.WithValue(ctx, traceKey, "abc-123")   // trace id
ctx = context.WithValue(ctx, requestIDKey, "xyz")   // request id
```

业务参数应该直接通过函数参数传递。Value 只适合传**横切关注点**（cross-cutting concerns）：trace ID、request ID、logger、认证 token 等贯穿整个调用链但不属于业务逻辑的数据。

---

## 6. Go 1.20~1.21 新增方法

### 6.1 WithCancelCause（Go 1.20+）

**函数原型：**

```go
func WithCancelCause(parent Context) (ctx Context, cancel CancelCauseFunc)
type CancelCauseFunc func(cause error)
```

| 项目 | 说明 |
|------|------|
| 返回值 `cancel` | 接受一个 `error` 参数，说明取消原因 |
| `ctx.Err()` | 取消后返回传入的 cause，而不是通用的 `context.Canceled` |

```go
ctx, cancel := context.WithCancelCause(parent)
cancel(fmt.Errorf("数据库连接失败"))

// 配合 context.Cause() 取出原因
err := context.Cause(ctx)  // → "数据库连接失败"
```

**Cause 函数（Go 1.20+）：**

```go
func Cause(ctx Context) error
```

递归向上查找，返回第一个 `WithCancelCause` 设置的原因。如果不是 `WithCancelCause` 取消的，返回 `context.Canceled` 或 `context.DeadlineExceeded`。

### 6.2 AfterFunc（Go 1.21+）

**函数原型：**

```go
func AfterFunc(ctx Context, f func()) (stop func() bool)
```

| 项目 | 说明 |
|------|------|
| `f` | context 取消后自动执行的清理函数 |
| 返回值 `stop` | 调用 stop() 可以取消执行 f；返回 true 表示成功拦截，false 表示 f 已执行 |

```go
stop := context.AfterFunc(ctx, func() {
    conn.Close()  // ctx 取消后自动关闭连接
})
// 如果不需要了：
stop()  // 取消执行清理函数
```

---

## 7. 标准库中的典型用法

### 7.1 net/http：每个请求自带 context

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // 客户端断开时自动取消

    result, err := doSomething(ctx)  // 把 ctx 传下去
    // ...
}
```

`r.Context()` 返回的 context 在客户端断开连接时自动取消。这就是为什么 HTTP handler 里启动的子 goroutine 应该在 ctx 取消时退出。

### 7.2 database/sql：查询超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

rows, err := db.QueryContext(ctx, "SELECT * FROM users")
```

查询超过 3 秒，数据库驱动收到 ctx 取消，停止等待数据库响应。

---

## 易错点

1. **忘记 defer cancel()**。`WithCancel`、`WithDeadline`、`WithTimeout` 返回的 `cancel` 必须调用。即使超时后会自动取消，`defer cancel()` 也不能省——它负责把 context 从父节点的 children map 里移除。不调会一直在那挂着，父节点永远引着子节点。

2. **以为函数返回后 context 会被 GC**。context 内部启动的定时器 goroutine（`AfterFunc`）会阻止它被回收，直到定时器到期或 cancel 被调用。所以不能依赖 GC。

3. **在 select 的 Done 分支里做耗时操作**。Done channel 关闭后，应该尽快退出。如果在 Done 分支里做同步 IO 或复杂计算，goroutine 不能及时释放，context 的取消效果就打折扣了。

4. **Value 的 key 用 string 导致冲突**。不同包都可能用 `"name"` 或 `"id"` 做 key，互相覆盖。用自定义类型避免。

5. **把 Context 存到 struct 里**。Context 是"本次调用的上下文"，不是对象的属性。存到 struct 意味着生命周期被延长，与设计意图背离。

6. **在主 goroutine 里阻塞等 ctx.Done() 但不给 default**。如果主 goroutine 自己阻塞在 `<-ctx.Done()` 上，谁来调 `cancel()`？→ 死锁。

7. **`ctx.Err()` 返回 nil 不代表不会取消**——`Err()` 返回 nil 表示"截至目前还没取消"，不等同于"以后也不会取消"。不能在取消前用它来"预判"。

8. **父 context 传 nil**——所有 `With*` 函数遇到 nil parent 都会 panic。不确定时用 `context.Background()` 或 `context.TODO()`。

---

## 快问快答

### Q1：Context 是什么？一句话说清楚。

Context 是 Go 里在 goroutine 之间传递取消信号、截止时间和请求范围数据的标准机制。它是一棵树，父节点取消则子节点全停。

### Q2：ctx.Done() 为什么返回 `<-chan struct{}` 而不是 `bool`？

因为 channel 可以同时被多个 goroutine 监听。关闭 channel 产生**广播效应**——所有阻塞在这个 channel 上的 goroutine 同时被唤醒。返回 bool 只能靠轮询。

### Q3：WithDeadline 和 WithTimeout 有什么区别？

`WithDeadline` 传入**绝对时间点**（"下午 3:00 截止"），`WithTimeout` 传入**相对时长**（"3 秒后截止"）。`WithTimeout` 源码就一行：`return WithDeadline(parent, time.Now().Add(timeout))`，本质完全一样，只是参数形式不同。

### Q4：为什么要 defer cancel()？不调会怎样？

两个原因：（1）`cancel()` 负责把子 context 从父节点的 children map 里移除——不调用则父节点一直持有子节点的引用，造成内存泄漏；（2）如果函数有多个 return 路径（错误提前返回），`defer` 保证不管走哪条路径 cancel 都会执行。

### Q5：context.WithValue 的查找方向是什么？

从当前节点**往上（父节点方向）**回溯。当前节点找不到就找父节点，一直找到根。子节点可以"覆盖"父节点同 key 的值（因为先查当前节点），但父节点看不到子节点的值（因为不往下找）。

### Q6：cancel 可以被调用多次吗？

可以。`cancel` 函数内部用了原子操作保证幂等——只有第一次调用生效，后续调用是空操作。重复调用不会 panic，也不会重复关闭 channel。

### Q7：select 里的 default 分支是什么意思？

`select { case <-ch: ...; default: ... }`：如果 ch 没准备好（没数据、没关闭），**不阻塞**，立刻执行 default。不带 default 的 `select { case <-ch: ... }`：没有其他 case，就**一直阻塞**，直到 ch 有数据或被关闭。

Watch 函数用带 default 的 select，就是为了"看一眼有没有取消信号，没有就继续干活"，而不是"等着取消"。

### Q8：Context 不能存 struct，那我想在 middleware 里设置值怎么办？

在 middleware 外层用 `WithValue` 创建新 context，然后传给内层：

```go
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), traceKey, generateTraceID())
        next.ServeHTTP(w, r.WithContext(ctx))  // 把新 ctx 放回 request
    })
}
```

关键：不是把 context 存到 middleware struct 里，而是每次请求时**创建新的 context**。

---

## Context 包工具速查

| 函数 | 完整签名 | cancel 返回值 | 核心行为 |
|------|---------|-------------|---------|
| `Background()` | `func Background() Context` | 无 | 创建永不取消的根 context |
| `TODO()` | `func TODO() Context` | 无 | 占位 context，语义同 Background |
| `WithCancel` | `func WithCancel(parent Context) (Context, CancelFunc)` | `CancelFunc`（手动调） | 创建一个可手动取消的子 context |
| `WithDeadline` | `func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)` | `CancelFunc`（手动/自动） | 创建定时点自动取消的子 context |
| `WithTimeout` | `func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)` | `CancelFunc`（手动/自动） | WithDeadline 的快捷方式 |
| `WithValue` | `func WithValue(parent Context, key, val any) Context` | 无 | 创建带键值对数据的子 context |
| `WithCancelCause` | `func WithCancelCause(parent Context) (Context, CancelCauseFunc)` | `CancelCauseFunc`（附带原因） | Go 1.20+ 取消时指定原因 |
| `Cause` | `func Cause(ctx Context) error` | / | Go 1.20+ 递归查找取消失败原因 |
| `AfterFunc` | `func AfterFunc(ctx Context, f func()) (stop func() bool)` | / | Go 1.21+ 取消后自动执行回调 |

| Context 接口方法 | 签名 | 行为 |
|-----------------|------|------|
| `Deadline()` | `Deadline() (deadline time.Time, ok bool)` | 返回截止时间；未设置时 ok=false |
| `Done()` | `Done() <-chan struct{}` | 返回只读 channel；取消时该 channel 被关闭 |
| `Err()` | `Err() error` | 返回 nil（未取消）、Canceled、DeadlineExceeded |
| `Value(key)` | `Value(key any) any` | 沿树往上查找 key，返回匹配的值或 nil |

---

## 一句话总结

Context 是 goroutine 的生命周期遥控器：`Background()` 是根，往上挂 `WithCancel`（手动停）、`WithDeadline`/`WithTimeout`（定时停）、`WithValue`（顺路带数据），然后作为第一个参数传下去，所有子 goroutine 通过 `<-ctx.Done()` 统一响应退出信号。cancel 必须 defer，key 要用自定义类型，Done 分支要快速退出。
