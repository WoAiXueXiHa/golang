# time 包：Timer 与 Ticker

## 这一章要记住什么

- `time.Timer` 是**一次性定时器**——到点触发一次就结束；`time.Ticker` 是**周期性定时器**——每隔固定时间触发一次，循环往复
- Timer 和 Ticker 都通过 **channel `C`（`<-chan time.Time`）** 传递触发信号；阻塞读 `<-timer.C` 或 `<-ticker.C` 就是等触发
- `Stop()` 停掉定时器后，**不能直接再读 `C`**——需要先判断 `Stop()` 的返回值，返回 `false` 时手动排空 channel
- `Reset()` 只对**已经停止或已触发**的 Timer 有效；对活跃的 Timer 要先 `Stop()` 再 `Reset()`
- `time.After` 是 `NewTimer(d).C` 的语法糖——简单但**不能在循环里用**，会导致 Timer 堆积无法及时回收
- `time.AfterFunc` 不走 channel，而是**到点直接在新 goroutine 里执行回调函数**，返回的 Timer 可以 `Stop()` 取消
- Ticker **不会自动停止**——不用了必须调 `Stop()`，否则内部 goroutine 永远不释放

---

## 0. Timer 和 Ticker 的本质区别

```
Timer（一次性闹钟）：
  NewTimer(3s) → ...3 秒后... → C 收到一个 time.Time 值 → 结束，不再触发

Ticker（周期心跳）：
  NewTicker(2s) → ...2s... → C 收到值 → ...2s... → C 收到值 → ...一直循环
  直到 Stop() 才会停
```

| 维度 | Timer | Ticker |
|------|-------|--------|
| 触发次数 | 1 次 | 无限次（直到 Stop） |
| 创建函数 | `NewTimer`、`AfterFunc`、`After` | `NewTicker` |
| 触发后行为 | 自动停止 | 继续周期触发 |
| 停止方式 | `Stop()` 或等它触发 | **必须** `Stop()` |
| 典型场景 | 超时控制、延迟执行 | 心跳上报、定时轮询、周期刷新 |

---

## 1. time.NewTimer —— 一次性定时器

### 类型定义与函数原型

Timer 结构体：

```go
type Timer struct {
    C <-chan Time     // 定时器触发时，当前时间发送到这个 channel
    // 内部字段（不导出）
}
```

**创建函数：**

```go
func NewTimer(d Duration) *Timer
```

| 项目 | 说明 |
|------|------|
| 参数 `d` | 超时时长，如 `2*time.Second` |
| 返回值 | `*Timer`，通过 `.C` 字段接收触发信号 |

**方法签名：**

```go
func (t *Timer) Stop() bool                // 停止定时器
func (t *Timer) Reset(d Duration) bool     // 重置定时器（Go 1.23+ 返回值改为 bool）
```

### 1.1 基本用法：创建并等待触发

**代码案例（对应 timer.go `create()`）：**

```go
timer := time.NewTimer(2 * time.Second)
<-timer.C                       // 阻塞等待，直到 2 秒后 channel 收到值
fmt.Println("after 2s Time out")
```

```
时间轴：
0s    NewTimer(2s) 创建，内部开始倒计时
      <-timer.C 阻塞等待...
2s    定时器触发 → 当前时间 time.Time 发送到 timer.C
      <-timer.C 收到值，解除阻塞 → 打印 "after 2s Time out"
```

`timer.C` 的类型是 `<-chan time.Time`（**只读** channel），这意味着调用方只能读，不能写——触发权完全在 Timer 内部。

### 1.2 Stop —— 提前停止定时器

**函数原型：**

```go
func (t *Timer) Stop() bool
```

| 项目 | 说明 |
|------|------|
| 参数 | 无 |
| 返回值 `bool` | `true` = 还没触发，停止成功；`false` = 已经触发或已停止，停止失败 |

**代码案例（对应 timer.go `cancel()`）：**

```go
timer := time.NewTimer(2 * time.Second)
res := timer.Stop()
fmt.Println(res)  // true：还没到 2 秒，停止成功
                  // false：已经过了 2 秒，timing 来晚了
```

**Stop 返回 false 时的铁律：必须排空 channel**

```go
if !timer.Stop() {
    <-timer.C  // 把已经触发的值从 channel 里读出来
}
```

为什么需要这步？

```
Stop 返回 true:
  Timer 还在倒计时 → Stop 成功 → C 里没东西 → 不用排空 ✓

Stop 返回 false:
  Timer 已经触发 → Stop 晚了 → C 里已经有一个值（或即将有）
  → 如果不 <-timer.C 排空：
    1. 下次代码读到 timer.C 时，可能会读到"上一轮"的旧值
    2. channel 里的值没人读，Timer 底层 runtime 可能无法回收
```

### 1.3 Reset —— 重置定时器

**函数原型：**

```go
func (t *Timer) Reset(d Duration) bool
```

| 项目 | 说明 |
|------|------|
| 参数 `d` | 新的超时时长 |
| 返回值 `bool` | Go 1.23+ 返回是否成功；之前版本无返回值 |

**代码案例（对应 timer.go `reset()`）：**

```go
timer := time.NewTimer(time.Second * 2)

<-timer.C                           // 等 2 秒，触发
fmt.Println("time out1")

res1 := timer.Stop()                // 已经触发过了 → 返回 false
fmt.Printf("res1 is %t\n", res1)

timer.Reset(time.Second)            // 重新倒计时 1 秒

res2 := timer.Stop()                // 还没到 1 秒 → 返回 true
fmt.Printf("res2 is %t\n", res2)
```

```
时间轴：
0s    NewTimer(2s)
2s    <-timer.C 收到值 → "time out1"
      timer.Stop() → false（已过期）
      timer.Reset(1s) → 重新倒计时
      timer.Stop() → true（还没到 1 秒，停止成功）
```

**Reset 的使用规则：**

```go
// ✅ 标准写法：先确保 Timer 停下来，再 Reset
if !timer.Stop() {
    <-timer.C        // 排空已触发的值
}
timer.Reset(d)        // 然后重新倒计时

// ❌ 错误：Timer 还在跑就直接 Reset
timer.Reset(d)        // 可能造成 channel 数据竞争
```

### 1.4 Timer 实战：select 超时控制

Timer 最经典的用法是和 `select` 配合，实现"等结果，但最多等 N 秒"：

```go
timer := time.NewTimer(3 * time.Second)
defer timer.Stop()

select {
case result := <-resultCh:       // 结果先回来了
    fmt.Println("成功:", result)
case <-timer.C:                  // 3 秒到了还没结果
    fmt.Println("超时!")
}
```

```
场景 1：API 提前返回（0.5s）
  0s    创建 Timer(3s)，启动 API goroutine
  0.5s  resultCh 收到数据 → select 走 result 分支
        timer 被 defer Stop() 停掉，不会在 3s 时触发

场景 2：超时
  0s    创建 Timer(3s)，启动 API goroutine
  3s    timer.C 触发 → select 走 timer.C 分支 → "超时!"
```

---

## 2. time.AfterFunc —— 到点执行回调，不走 channel

**函数原型：**

```go
func AfterFunc(d Duration, f func()) *Timer
```

| 项目 | 说明 |
|------|------|
| 参数 `d` | 延迟时长 |
| 参数 `f` | 到点后在新 goroutine 中执行的回调函数 |
| 返回值 | `*Timer`，可以用 `Stop()` 在回调执行前取消 |

`AfterFunc` 和 `NewTimer` 的本质区别：**它不往 channel 发信号，而是到点直接执行回调**。

**代码案例（对应 timer.go `dur()`）：**

```go
duration := time.Duration(1) * time.Second

f := func() {
    fmt.Println("f has been called after 1s by time.AfterFunc")
}

timer := time.AfterFunc(duration, f)
defer timer.Stop()            // 如果 1 秒内函数返回了，取消回调

time.Sleep(2 * time.Second)   // 等 2 秒确保回调执行完
```

```
时间轴：
0s    AfterFunc(1s, f) → timer 创建，开始倒计时
      defer timer.Stop()（函数退出时执行，这里等到 2s 后才退出）
1s    倒计时到 → f 在新 goroutine 中执行 → "f has been called..."
2s    主 goroutine 睡醒 → 函数退出 → defer timer.Stop()（此时回调已执行，Stop 无影响）
```

**和 NewTimer 的对比：**

| 方式 | 通知机制 | 如何取消 | 适合场景 |
|------|---------|---------|---------|
| `NewTimer` | `<-timer.C` channel 通知 | `Stop()` | 需要和 select 配合的超时控制 |
| `AfterFunc` | 直接在新 goroutine 执行回调 | `Stop()` 可阻止回调执行 | "到点做某事"，不需要 select |

**AfterFunc 最大的优势：** 它比 `time.Sleep` + goroutine 更可控——因为返回的 `*Timer` 可以 `Stop()`：

```go
// ❌ time.Sleep 方式：没法取消
go func() {
    time.Sleep(5 * time.Second)
    cleanup()  // 即使主流程已经不需要了，也会执行
}()

// ✅ AfterFunc 方式：可以取消
timer := time.AfterFunc(5*time.Second, cleanup)
// 如果中途想取消：
timer.Stop()  // cleanup 不会执行了
```

---

## 3. time.After —— 最方便的语法糖，也有最隐蔽的坑

**函数原型：**

```go
func After(d Duration) <-chan Time
```

| 项目 | 说明 |
|------|------|
| 参数 `d` | 延迟时长 |
| 返回值 | `<-chan Time`，d 时间后收到一个 time.Time 值 |

**源码等价：**

```go
func After(d Duration) <-chan Time {
    return NewTimer(d).C
}
```

只返回 channel，**不返回 Timer 对象**——所以没法 `Stop()`。

**代码案例（对应 timer.go `after()`）：**

```go
ch := make(chan string)

go func() {
    time.Sleep(time.Second * 3)
    ch <- "test"
}()

select {
case val := <-ch:
    fmt.Printf("val is %s\n", val)
case <-time.After(time.Second * 2):    // 2 秒超时
    fmt.Println("timeout!!!")
}
```

```
时间轴：
0s    启动 goroutine（3 秒后才往 ch 写数据）
      select 等待 ch 或 time.After(2s)
2s    time.After 的 channel 触发 → select 走 case 2 → "timeout!!!"
      （ch 的数据在 3s 时才就绪，但 select 已经结束了）
```

**单次使用没问题。循环中反复调用会内存泄漏：**

```go
// ❌ 循环中反复调用 time.After
for {
    select {
    case <-ch:
        // 处理数据
    case <-time.After(timeout):  // ← 每次循环创建一个新 Timer！
        // 如果 ch 一直有数据，After 分支永远走不到
        // → 每个 Timer 留在内存里直到超时 → 严重堆积
    }
}
```

```
假设每秒处理 1000 次循环，timeout=2s：

迭代 1: 创建 Timer(2s)，ch 有数据 → 走 ch 分支 → Timer 留着等 2 秒
迭代 2: 创建 Timer(2s)，ch 有数据 → 走 ch 分支 → Timer 留着等 2 秒
...
迭代 2000: 创建 Timer(2s) → 同时第一个 Timer 刚超时被回收
→ 内存里始终挂着 ~2000 个 Timer，每个都占着 runtime 资源
```

**✅ 正确做法：循环中复用 NewTimer + Reset**

```go
timer := time.NewTimer(timeout)
defer timer.Stop()

for {
    if !timer.Stop() {
        <-timer.C           // 排空已触发的值
    }
    timer.Reset(timeout)    // 重新倒计时
    select {
    case <-ch:
        // 处理
    case <-timer.C:
        // 超时
    }
}
```

---

## 4. time.NewTicker —— 周期性定时器

### 类型定义与函数原型

Ticker 结构体：

```go
type Ticker struct {
    C <-chan Time     // 每隔 d 时间，当前时间发送到这个 channel
    // 内部字段（不导出）
}
```

**创建函数：**

```go
func NewTicker(d Duration) *Ticker
```

| 项目 | 说明 |
|------|------|
| 参数 `d` | 触发间隔，如 `1*time.Second` |
| 返回值 | `*Ticker`，通过 `.C` 字段接收周期性触发信号 |

**方法签名：**

```go
func (t *Ticker) Stop()                   // 停止 Ticker，关闭内部 goroutine
func (t *Ticker) Reset(d Duration)        // Go 1.15+ 重置间隔（先停再设）
```

### 4.1 基本用法：创建、监听、停止

**代码案例（对应 ticker.go `demo1()`）：**

```go
Watch := func() chan struct{} {
    ticker := time.NewTicker(time.Second)

    ch := make(chan struct{})

    go func(ticker *time.Ticker) {
        defer ticker.Stop()           // ① defer 保证退出时停掉 Ticker
        for {
            select {
            case <-ticker.C:          // ② 每秒收到一次信号
                fmt.Println("watch!!!")
            case <-ch:                // ③ 收到停止信号
                fmt.Println("Ticker Stop!!!")
                return
            }
        }
    }(ticker)
    return ch
}

ch := Watch()
time.Sleep(5 * time.Second)             // 让 ticker 触发 5 次左右
ch <- struct{}{}                        // 发送停止信号
close(ch)
```

```
时间轴：
0s    NewTicker(1s) → ticker 创建，开始计时
1s    ticker.C 触发 → "watch!!!"
2s    ticker.C 触发 → "watch!!!"
3s    ticker.C 触发 → "watch!!!"
4s    ticker.C 触发 → "watch!!!"
5s    ticker.C 触发 → "watch!!!"
      ch <- struct{}{} → select 走 ch 分支 → "Ticker Stop!!!"
      defer ticker.Stop() 执行 → Ticker 停止
```

**关于 "Ticker Stop!!! 打印不稳定" 的注释：**

代码注释指出了这个现象。原因：`time.Sleep(5 * time.Second)` 不能精确保证刚好 5 个 tick——第 5 秒时 ticker 的触发和 ch 的写入几乎同时发生，select 随机选一个就绪的 case。有时 ticker.C 先抢到，多打印一次 "watch!!!"；有时 ch 先抢到，打印 "Ticker Stop!!!"。

### 4.2 Ticker 不会自动停止

Ticker 一旦创建就开始周期触发，**直到你调用 `Stop()`**。即使你不再读 `ticker.C`，它内部的 goroutine 和 channel 仍然活着——这就是 goroutine 泄漏。

```go
// ❌ goroutine 泄漏：Ticker 永远在后台跑
ticker := time.NewTicker(time.Second)
// 用了几次后忘了 Stop，函数返回 → ticker 成为垃圾
// 但它的内部 goroutine 还活着，永远不会被 GC

// ✅ 正确：defer Stop
ticker := time.NewTicker(time.Second)
defer ticker.Stop()
```

### 4.3 Ticker 不能暂停

Ticker 没有 `Pause` / `Resume` 方法。需要"暂停一段时间再继续"的场景：

```go
// 方案：Stop + 之后重新 NewTicker
ticker.Stop()
// ... 暂停期间 ...
ticker = time.NewTicker(time.Second)  // 恢复
```

### 4.4 Ticker 的触发间隔不保证精确

如果接收端处理慢，还没读完上一次的值，Ticker 会**跳过这次触发**：

```
理想情况（处理 < 1s）：
  tick 1s → 处理(200ms) → 等 800ms → tick 2s → 处理(200ms) → ...

实际情况（处理 > 1s）：
  tick 1s → 处理(1.5s) → tick 2s 被跳过！→ 直接收到 tick 3s 附近的值
  实际间隔 = max(设定间隔, 处理耗时)
```

### 4.5 单 Ticker 驱动多频率任务

不需要多个 Ticker 各管各的，一个 Ticker + 计数器就能分频：

```go
ticker := time.NewTicker(1 * time.Second)  // 基础粒度 1 秒
defer ticker.Stop()
counter := 0

for {
    select {
    case <-ticker.C:
        counter++
        checkHealth()         // 每 1 秒

        if counter%3 == 0 {
            fetchConfig()     // 每 3 秒
        }
        if counter%10 == 0 {
            fullReport()      // 每 10 秒
        }
    case <-stopCh:
        return
    }
}
```

```
tick #1  (1s):  checkHealth
tick #2  (2s):  checkHealth
tick #3  (3s):  checkHealth + fetchConfig
...
tick #10 (10s): checkHealth + fullReport
```

### 4.6 Ticker 的 select 退出注意事项

Ticker 通过 select + done channel 退出时，**正在执行的任务不会被 select 中断**：

```
ticker 触发 → 开始执行 doWork()（耗时 800ms）
              ... 600ms 后收到停止信号 → stopCh 就绪
              但 select 还在等 doWork() 完成！
              → 当前任务执行完，下次循环 select 才走到 stopCh 分支
```

这是因为 Go 的 select 只关心 channel 操作是否就绪，不能中断正在执行的代码。如果需要可中断的耗时任务，在 `doWork` 里传入 context。

---

## 函数速查

| 函数/方法 | 完整签名 | 返回值说明 | 用途 |
|----------|---------|-----------|------|
| `NewTimer` | `func NewTimer(d Duration) *Timer` | `*Timer`，通过 `.C` 接收触发 | 创建一次性定时器 |
| `(*Timer).C` | `C <-chan Time`（字段） | 触发时收到一个 `time.Time` 值 | 阻塞等待触发或用于 select |
| `(*Timer).Stop` | `func (t *Timer) Stop() bool` | `true`=停止成功；`false`=已触发 | 提前停止；返回 false 必须排空 C |
| `(*Timer).Reset` | `func (t *Timer) Reset(d Duration) bool` | Go 1.23+: `true`=重置成功 | 重新倒计时 d |
| `AfterFunc` | `func AfterFunc(d Duration, f func()) *Timer` | `*Timer`，用 `Stop()` 取消回调 | 到点执行回调，不走 channel |
| `After` | `func After(d Duration) <-chan Time` | `<-chan Time`，等价 `NewTimer(d).C` | select 超时的语法糖；忌循环 |
| `NewTicker` | `func NewTicker(d Duration) *Ticker` | `*Ticker`，通过 `.C` 接收周期触发 | 创建周期性定时器 |
| `(*Ticker).C` | `C <-chan Time`（字段） | 每 d 时间收到一个 `time.Time` 值 | 阻塞等待触发或用于 select |
| `(*Ticker).Stop` | `func (t *Ticker) Stop()` | 无返回值 | 停止 Ticker；必须调用 |
| `(*Ticker).Reset` | `func (t *Ticker) Reset(d Duration)` | 无返回值 | Go 1.15+ 重置间隔 |

---

## 易错点

1. **`timer.Stop()` 返回 false 时不排空 channel**。返回 false 说明 Timer 已经触发，`C` 里有值（或即将有）。不排空 → 下次读到旧值、底层资源无法回收。必须 `if !timer.Stop() { <-timer.C }`。

2. **`time.After` 在循环中使用**。每次 `time.After(d)` 创建新 Timer，无法 Stop。如果 select 频繁走其他分支，这些 Timer 堆积在内存里直到超时。循环中必须用 `NewTimer` + `Reset` 复用。

3. **Ticker 忘了 Stop**。Ticker 创建后一直运行，即使不再读 `ticker.C`。不 Stop → goroutine 泄漏。**永远 `defer ticker.Stop()`**。

4. **对正在运行的 Timer 直接 Reset**。Go 1.23 之前可能导致 channel 数据竞争。标准写法：先 `Stop()` → 排空 → 再 `Reset()`。

5. **`timer.Stop()` 返回 false 但不去读 `C` 直接 Reset**。同样会造成 channel 残留数据。正确顺序：Stop → 排空 → Reset。

6. **以为 Ticker 触发间隔精确**。如果接收端处理耗时 > 间隔，Ticker 会跳过中间触发。实际间隔 = max(设定间隔, 处理耗时)。

7. **在 Timer 的 `C` 上同时阻塞多个 goroutine**。`timer.C` 是普通的 channel，只有一个 goroutine 能读到触发值。如果多个 goroutine 都 `<-timer.C`，只有一个能收到。

8. **`time.AfterFunc` 的回调里做耗时操作**。回调在新 goroutine 执行，但它是同步的——如果回调很慢，不会影响定时器其他功能，但要留意 goroutine 数量。

---

## 快问快答

### Q1：Timer 和 Ticker 的本质区别是什么？

Timer 是一次的，到点触发就停；Ticker 是周期的，每隔固定时间触发，直到显式 Stop。Timer 像闹钟，Ticker 像心跳。

### Q2：`time.After` 和 `time.NewTimer` 的区别？

`time.After(d)` 就是 `NewTimer(d).C`，只返回 channel 不返回 Timer 对象，无法 Stop。单次 select 用没问题，**循环里绝对不能**——因为 Timer 无法提前回收。`NewTimer` 返回 Timer 对象，可以 Stop 和 Reset。

### Q3：`timer.Stop()` 返回 false 是什么意思？怎么处理？

返回 false = Timer 已经触发（或已 Stop 过了）。`C` 里可能有一个待读的值。必须 `if !timer.Stop() { <-timer.C }` 排空，否则后续代码会误读旧值，且底层资源无法释放。

### Q4：`AfterFunc` 和 `NewTimer` + goroutine 的区别？

`AfterFunc(d, f)` 到点自动在新 goroutine 执行 f，不走 channel，返回的 Timer 可以 Stop 取消。`NewTimer` 通过 `C` 发信号，由你决定信号到了干什么。前者适合"到点做某事"，后者适合"到点通知某个 select"。

### Q5：Ticker 能暂停吗？能重置间隔吗？

不能 Pause/Resume。需要暂停就 Stop，恢复时重新 NewTicker。Go 1.15+ 可以用 `ticker.Reset(d)` 重置间隔（内部先 Stop 再设新的）。

### Q6：`AfterFunc` 里 Stop 后回调还会执行吗？

如果在回调执行前调了 `Stop()`，回调**不会执行**。如果回调已经开始执行了，`Stop()` 无法中断——Go 没有中断 goroutine 的机制。

### Q7：Ticker 的触发间隔精确吗？

不精确。处理慢时会跳过中间触发。实际间隔 = max(设定的间隔, 处理耗时)。需要高精度定时请用 `time.Timer` 每次 Reset。

### Q8：为什么 ticker.go 里 "Ticker Stop!!!" 打印不稳定？

因为 `time.Sleep(5*time.Second)` 之后，第 5 秒的 ticker 触发和 `ch <- struct{}{}` 几乎同时发生。select 在多个 case 同时就绪时**随机选择**——有时 ticker.C 先被选中（多打一次 "watch!!!"），有时 ch 先被选中（打印 "Ticker Stop!!!"）。

---

## 一句话总结

Timer 是点到即止的闹钟（一次触发），Ticker 是永不停歇的心跳（周期触发）。用 `NewTimer` / `NewTicker` 创建，从 `.C` 读信号；`AfterFunc` 直接回调不经过 channel；`After` 是语法糖但忌循环。不管哪种，不用就 Stop：Timer 的 Stop 失败要排空 `C`，Ticker 的 Stop 忘调就是 goroutine 泄漏。
