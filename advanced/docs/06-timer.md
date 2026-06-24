# time 包：定时器与周期任务

## 这一章要记住什么

- `time.Timer` 是**一次性定时器**——到点触发一次就结束；`time.Ticker` 是**周期性定时器**——每隔固定时间触发一次，循环往复。
- Timer 和 Ticker 都通过 **channel `C`** 传递触发信号，典型写法是 `<-timer.C` 或 `<-ticker.C`。
- `Stop()` 停掉定时器后，**不能直接再读 `C`**——需要先判断 `Stop()` 的返回值，返回 `false` 时手动排空 channel。
- `Reset()` 只对**已经停止或已触发**的 Timer 有效，对活跃的 Timer 要先 `Stop()` 再 `Reset()`。
- `time.After` 是 `NewTimer(d).C` 的语法糖——简单但**不能在循环里用**，会导致 Timer 堆积无法及时回收。
- `time.AfterFunc` 不会往 channel 发信号，而是**到点直接执行回调函数**，返回的 Timer 可以 `Stop()` 取消执行。
- Ticker 是**不会自动停止**的——不用了必须调 `Stop()`，否则内部 goroutine 永远不释放。

---

## 1. Timer 和 Ticker 的本质区别

Timer 和 Ticker 都是 Go 的定时器，但它们的触发方式完全不同：

```text
Timer（一次性）：
  创建 → 等待 → ⏰ 触发一次 → 结束
  NewTimer(3s) → ...3 秒后... → C 收到信号 → Timer 不再触发

Ticker（周期性）：
  创建 → 等待 → ⏰ 触发 → 等待 → ⏰ 触发 → 等待 → ...
  NewTicker(2s) → ...2s... → C 收到信号 → ...2s... → C 收到信号 → ...
  直到 Stop() 才会停
```

用一句话区分：**Timer 是闹钟，响一次就歇了；Ticker 是心跳，一直跳到你让它停。**

| 维度 | Timer | Ticker |
|------|-------|--------|
| 触发次数 | 1 次 | 无限次（直到 Stop） |
| 典型场景 | 超时控制、延迟执行 | 心跳上报、定时轮询 |
| 停止方式 | Stop() 或等它触发 | 必须 Stop() |
| 创建函数 | `NewTimer`、`AfterFunc`、`After` | `NewTicker` |

---

## 2. time.NewTimer —— 一次性定时器

### 2.1 基本用法

```go
timer := time.NewTimer(3 * time.Second)  // 创建 3 秒定时器
<-timer.C                                 // 阻塞等待触发
fmt.Println("3 秒到了")
```

`NewTimer` 返回一个 `*time.Timer`，它有一个字段 `C`（`<-chan time.Time` 类型）。定时器到期时，当前时间会被发送到 `C`。

**案例：API 超时控制**

```go
timer := time.NewTimer(3 * time.Second)

go func() {
    resultCh <- callExternalAPI()  // 异步调用 API
}()

select {
case result := <-resultCh:
    timer.Stop()  // API 提前返回，停止定时器
    fmt.Println("成功:", result)
case <-timer.C:
    fmt.Println("超时！")
}
```

```text
时间线（API 提前返回的情况）：
  0s:  创建 Timer(3s)，启动 API goroutine
  0.5s: API 返回 → resultCh 可读 → select 走 result 分支
        timer.Stop() → Timer 被停掉，不会在 3s 时触发

时间线（超时的情况）：
  0s:  创建 Timer(3s)，启动 API goroutine
  3s:  Timer 触发 → timer.C 可读 → select 走 timer.C 分支
       API goroutine 仍在跑，但结果被丢弃
```

### 2.2 Stop —— 提前停止 Timer

```go
stopped := timer.Stop()
// stopped == true  → Timer 还没触发，停止成功
// stopped == false → Timer 已经触发了，停止失败
```

**关键坑：Stop 返回 false 时必须排空 channel**

```go
if !timer.Stop() {
    <-timer.C  // 把已经触发的值从 channel 里读出来
}
```

为什么需要这步？因为 Timer 触发后，它的值可能还在 `C` 里排队。如果不排空，后续代码如果又去读 `timer.C`，可能会读到**上一次**的触发值，造成逻辑错误。更严重的是，这个值没人读的话，Timer 底层的 goroutine 会因为 channel 阻塞而无法被 GC 回收。

```text
Stop 返回 true 的情况：
  Timer 还在倒计时 → Stop 成功 → C 里没东西 → 不用排空

Stop 返回 false 的情况：
  Timer 已经触发 → Stop 晚了 → C 里已经有一个值（或即将有）
  → 必须 <-timer.C 排空
```

### 2.3 Reset —— 重置 Timer

```go
timer.Reset(5 * time.Second)  // 重新倒计时 5 秒
```

**案例：超时后"再给一次机会"**

在订单支付场景中，第一次超时后可能还想再等等：

```go
timer := time.NewTimer(3 * time.Second)

select {
case result := <-resultCh:
    timer.Stop()
    // 处理成功

case <-timer.C:
    // 第一次超时了，但还是想等等看
    timer.Reset(2 * time.Second)  // 再给 2 秒

    select {
    case result := <-resultCh:
        timer.Stop()
        // 延迟恢复，处理成功
    case <-timer.C:
        // 最终超时
    }
}
```

```text
时间线：
  0s:  NewTimer(3s)
  3s:  超时触发 → Reset(2s)
  5s:  如果这时 resultCh 还没数据 → 再次超时，彻底放弃
```

**Reset 的使用规则：**

```go
// ✅ 正确：Timer 已触发或已 Stop，直接 Reset
if !timer.Stop() {
    <-timer.C
}
timer.Reset(d)

// ❌ 错误：Timer 还在运行就 Reset（Go 1.23+ 已修复，但仍推荐先 Stop）
timer.Reset(d)  // 可能造成 channel 竞争
```

### 总结一下

`NewTimer(d)` 创建一个一次性定时器，到点后往 `C` 发信号。`Stop()` 提前终止，返回 `false` 时要排空 `C`。`Reset(d)` 给已触发的 Timer 重新倒计时。三个方法组合起来可以灵活控制超时逻辑。

---

## 3. time.AfterFunc —— 到点执行回调

`AfterFunc` 和 `NewTimer` 的本质区别：**它不走 channel，而是直接在新 goroutine 里执行你给的函数。**

```go
timer := time.AfterFunc(5*time.Second, func() {
    fmt.Println("5 秒后执行清理")
})

// 如果中途想取消：
timer.Stop()  // 回调不会执行了
```

**案例：用户会话空闲超时清理**

用户 5 秒不操作 → 自动踢下线。但如果用户中途操作了 → 取消倒计时，重新开始。

```go
var idleTimer *time.Timer

renewTimer := func() {
    if idleTimer != nil {
        idleTimer.Stop()  // 取消旧倒计时
    }
    idleTimer = time.AfterFunc(5*time.Second, func() {
        fmt.Println("用户被踢下线")
    })
}

renewTimer()           // 开始倒计时
time.Sleep(2 * time.Second)
renewTimer()           // 用户操作了 → 取消旧的，重新倒计时 5 秒
```

```text
操作时间线：
  0s:  AfterFunc(5s, 踢人) → 倒计时开始
  2s:  用户操作 → Stop() 旧的 → AfterFunc(5s, 踢人) → 倒计时重新从 2s 开始
  7s:  5 秒到了 → 踢人回调执行
```

**和 NewTimer 的对比：**

| 方式 | 通知机制 | 适合场景 |
|------|---------|---------|
| `NewTimer` | `<-timer.C` channel 通知 | 需要和 select 配合的超时控制 |
| `AfterFunc` | 自动调用回调函数 | "到点做某事"，不需要 select |

**AfterFunc 最关键的优势**：它让你拿到一个可以 `Stop()` 的 Timer 对象。如果你只是想在 N 秒后执行一个函数，用 `AfterFunc` 比 `time.Sleep` + goroutine 更可控——因为可以取消。

### 总结一下

`AfterFunc(d, f)` 在 d 时间后在独立 goroutine 里调用 f，返回的 Timer 可以被 Stop 取消。它适合"到点执行清理、过期处理"这类不参与 select 的场景。最大的优势是可以随时取消，这是 `time.Sleep` 做不到的。

---

## 4. time.After —— 最方便的语法糖，也有最隐蔽的坑

`time.After(d)` 等价于 `time.NewTimer(d).C`，但它只返回一个 `<-chan time.Time`，**不返回 Timer 对象**。这意味着你没法 Stop 它。

```go
select {
case result := <-resultCh:
    fmt.Println("成功")
case <-time.After(2 * time.Second):
    fmt.Println("超时")
}
```

这段代码很简洁，但问题出在 resultCh 先就绪的路径上：

```text
正常情况（resultCh 先返回）：
  select 走 result 分支 → time.After 创建的 Timer 已经没人管了
  → Timer 在 2 秒后触发 → 底层 channel 里的值被 GC 回收（Go 1.23+）
  → 问题不大，单次使用可以接受

循环中反复调用（危险！）：
  for {
      select {
      case <-ch:
          // 处理
      case <-time.After(timeout):  // ← 每次循环创建一个新 Timer！
          // 如果 ch 一直有数据，After 分支永远走不到
          // → Timer 堆积，内存泄漏
      }
  }
```

```text
循环中 Timer 堆积：

迭代 1: 创建 Timer(2s)，ch 有数据 → 走 ch 分支 → Timer 留在内存里等 2 秒
迭代 2: 创建 Timer(2s)，ch 有数据 → 走 ch 分支 → Timer 留在内存里等 2 秒
迭代 3: 创建 Timer(2s)，ch 有数据 → 走 ch 分支 → Timer 留在内存里等 2 秒
...
每秒处理 1000 个请求 → 每 2 秒就有 2000 个 Timer 在内存里等待 GC
```

**正确做法：在循环里用 NewTimer + Reset 复用同一个 Timer**

```go
timer := time.NewTimer(timeout)
defer timer.Stop()

for {
    timer.Reset(timeout)  // 复用同一个 Timer
    select {
    case <-ch:
        // 处理
    case <-timer.C:
        // 超时
    }
}
```

### 总结一下

`time.After` 是 `NewTimer(d).C` 的语法糖，写起来方便但不能 Stop。单次 select 里用没问题，**循环里反复调用会导致 Timer 堆积**。循环场景用 `NewTimer` + `Reset` 复用。

---

## 5. time.NewTicker —— 周期性定时器

### 5.1 基本用法

```go
ticker := time.NewTicker(2 * time.Second)
defer ticker.Stop()  // 必须 Stop！否则内部 goroutine 泄漏

for {
    select {
    case <-ticker.C:
        fmt.Println("每 2 秒执行一次")
    case <-done:
        return
    }
}
```

`ticker.C` 每 2 秒收到一个值，循环往复，直到调用 `ticker.Stop()`。

**案例：服务心跳上报**

```go
ticker := time.NewTicker(2 * time.Second)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        sendHeartbeat()  // 每 2 秒向注册中心上报心跳
    case <-stopCh:
        return           // 服务关闭，Stop 由 defer 执行
    }
}
```

### 5.2 单 Ticker 驱动多个周期任务

不需要为每个周期创建不同的 Ticker。用一个 Ticker + 计数器取模就能区分不同频率：

```go
ticker := time.NewTicker(1 * time.Second)  // 基础粒度：1 秒
counter := 0

for {
    select {
    case <-ticker.C:
        counter++
        checkHealth()      // 每 1 秒执行

        if counter%3 == 0 {
            fetchConfig()  // 每 3 秒执行（每 3 次 tick）
        }
    }
}
```

```text
tick 驱动的时间线：
  tick #1 (1s): 健康检查
  tick #2 (2s): 健康检查
  tick #3 (3s): 健康检查 + 拉取配置  ← counter%3==0
  tick #4 (4s): 健康检查
  tick #5 (5s): 健康检查
  tick #6 (6s): 健康检查 + 拉取配置  ← counter%3==0
```

### 5.3 Ticker 不能暂停

Ticker 没有 Pause/Resume。如果需要"暂停一段时间再继续"，只能 Stop 当前 Ticker，之后重新 NewTicker。

### 总结一下

`NewTicker(d)` 创建一个每 d 时间触发一次的定时器，通过 `ticker.C` 接收信号。不用必须 `Stop()`。一个 Ticker 可以通过计数器分频来驱动多个不同频率的周期任务。Ticker 不可暂停，只能 Stop 后重建。

---

## 6. Ticker + select 的优雅退出

Ticker 最常见的死法就是 **Stop 没调，goroutine 泄漏**。标准写法：

```go
ticker := time.NewTicker(3 * time.Second)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        doWork()
    case <-stopCh:
        return  // defer 执行 ticker.Stop()
    }
}
```

**如果在 doWork() 执行过程中收到了停止信号怎么办？**

```text
场景：备份任务正在执行（耗时 800ms），这时收到 SIGTERM

  ticker 触发 → 开始备份（800ms）
                ...600ms 后...
                收到 SIGTERM → stopCh 关闭
                但 select 还在等 doWork() 完成！
                
  → 下一次 select 循环才会走到 stopCh 分支
  → 当前正在执行的备份会完成，不会被中途打断
```

这是 Go 的 select 机制决定的——select 是等 case 里的 channel 操作就绪，而不是中断正在执行的代码。如果需要能中断的长时间任务，应该把上下文 context 传到 doWork 里。

### 总结一下

Ticker 的退出模式就是 `defer Stop()` + `select { case ticker.C: ... case done: return }`。当前在执行的任务不会被 select 中断——如果需要可中断的耗时任务，需要结合 context。

---

## 易错点

1. **`timer.Stop()` 返回 false 时不排空 channel**。`Stop` 返回 false 说明 Timer 已经触发，值已经（或即将）在 `C` 里。不排空会导致下次读 `C` 时读到旧值，且 Timer 无法被 GC。

2. **`time.After` 在循环中用**。每次 `time.After(d)` 创建一个新 Timer，如果 select 频繁走其他分支，这些 Timer 堆积在内存里，直到超时才能被回收。循环中复用 `NewTimer` + `Reset`。

3. **Ticker 忘了 Stop**。Ticker 一旦创建就会一直运行，即使你不再读 `ticker.C`。不 Stop 意味着它底层的 goroutine 和 channel 永远不会被释放。

4. **对正在运行的 Timer 直接 Reset**。Go 1.23 之前这是未定义行为，可能造成 channel 数据竞争。应先 `Stop()` 再 `Reset()`。

5. **把 `ticker.C` 当成"准时触发"**。Ticker 不保证精确的周期——如果接收端处理慢，触发间隔会被拉长。Ticker 是"每隔 d 时间尝试发送"，如果上次的值还没被读走，这次会跳过。

6. **`AfterFunc` 的回调里做耗时操作**。回调在独立的 goroutine 中执行，不要在里面做阻塞操作——如果要做的活比较重，在回调里再启动一个 goroutine。

---

## 快问快答

### Q1：Timer 和 Ticker 的本质区别是什么？

答：Timer 是一次性的，到点触发一次就结束；Ticker 是周期性的，每隔固定时间触发一次，直到显式 Stop。Timer 像闹钟，Ticker 像心跳。

### Q2：`time.After` 和 `time.NewTimer` 的区别？

答：`time.After(d)` 等同 `NewTimer(d).C`，只返回 channel 不返回 Timer，无法 Stop。适合一次性的 select 超时，但**不能在循环里用**——因为无法 Stop 导致 Timer 堆积。`NewTimer` 返回 Timer 对象，可以 Stop 和 Reset，更可控。

### Q3：`timer.Stop()` 返回 false 是什么意思？怎么处理？

答：返回 false 说明 Timer 已经触发（或已经被 Stop 过了），`C` 里可能已经有一个值。必须用 `<-timer.C` 把值排空，否则后续代码可能误读旧值，且 Timer 底层资源无法释放。

### Q4：`AfterFunc` 和 `NewTimer` + goroutine 的区别？

答：`AfterFunc(d, f)` 内部自动启动 goroutine，到点执行 f，没有 channel 可以 select。`NewTimer` 通过 `C` 发信号，由你决定收到信号后做什么。前者适合"到点做某事"，后者适合"到点通知某处"。

### Q5：Ticker 能暂停吗？

答：不能。没有 Pause/Resume 方法。需要暂停就 Stop 掉，需要恢复时重新 NewTicker。

### Q6：`time.AfterFunc` 里 Stop 后，回调还会执行吗？

答：如果在回调执行前调了 `Stop()`，回调**不会执行**。如果回调已经开始执行了，`Stop()` 无法中断它——Go 没有提供中断 goroutine 的机制。

### Q7：Ticker 的触发间隔精确吗？

答：不精确。如果接收端处理太慢，还没有来得及读上一次的值，Ticker 会跳过这次触发。实际间隔 = max(设定间隔, 处理耗时)。

---

## 一句话总结

Timer 是点到即止的闹钟（一次触发），Ticker 是永不停歇的心跳（周期触发）——用 `NewTimer`/`NewTicker` 创建，通过 channel 接收信号，`AfterFunc` 直接回调，`After` 是语法糖但忌循环。不管哪种，不用就 Stop，Stop 失败就排空 C。
