# sync 包：并发安全与同步原语

## 这一章要记住什么

- `sync` 包提供 Go 并发编程的核心工具箱——**WaitGroup**（等协程完成）、**Mutex/RWMutex**（保护临界区）、**Once**（保证只执行一次）、**Cond**（条件变量）、**Pool**（临时对象复用）
- **原生 map 不是并发安全的**——多 goroutine 同时写会直接 fatal；小并发用 Mutex 保护，特定场景用 sync.Map
- `sync.Mutex` **绝不能值拷贝**——拷贝会把锁的内部状态一起复制，导致副本里的 Lock 永远等不到解锁 → 死锁
- **原子操作**（`sync/atomic`）比 Mutex 更快，不涉及 OS 调度，适合简单计数器和状态标志
- 锁的获取顺序不一致会导致**循环等待死锁**——所有 goroutine 必须按相同顺序加锁
- 读写锁把读者和写者分开，**读读共享、读写互斥、写写互斥**，适合读多写少的场景
- **sync.Pool 是临时对象池**——Get/Put 复用高频临时对象，减少 GC 压力；但对象随时可能被 GC 回收，不能当持久缓存
- **sync.Cond** 是条件变量，让 goroutine 在某个条件满足前挂起等待，适合"等信号再干活"的场景

---

## 0. chan struct{}：零成本信号通知

在进入 sync 包之前，先看一种不依赖 sync 包的轻量级同步方式——用 channel 发信号。

### 函数原型

```go
// 创建无缓冲信号 channel
ch := make(chan struct{})      // 无缓冲：发送方必须等接收方就绪
ch := make(chan struct{}, 10)  // 有缓冲：最多缓存 10 个信号，发送方不阻塞
```

`struct{}` 是 Go 里**唯一零内存占用**的类型（`unsafe.Sizeof(struct{}{}) == 0`），专门用于表达"我不传数据，只传信号"。

### 使用方式

```go
done := make(chan struct{})  // 无缓冲信号 channel

go func() {
    // 干活...
    done <- struct{}{}        // 发送：塞一个空结构体，意思是"搞定了"
}()

<-done                        // 接收：阻塞等待，直到收到信号
```

```
子 goroutine: 干活中... → done <- struct{}{}  ← 发送信号
                                            ↓
主 goroutine: <-done 阻塞等待 ───────────────→ 收到信号，继续执行
```

### 参数和返回值

| 操作 | 签名 | 说明 |
|------|------|------|
| 发送 | `ch <- struct{}{}` | 无返回值；无缓冲时阻塞直到有人接收 |
| 接收 | `<-ch` | 返回 `struct{}{}`（零字节），通常不赋值直接用 |

### 注意事项

- 无缓冲 channel 必须两端同时就绪，否则发送/接收方都会阻塞
- 如果接收方先退出而发送方还在等，发送方会永久阻塞 → goroutine 泄漏
- 单纯"等 N 个协程完成"的场景，用 `sync.WaitGroup` 更直观

---

## 1. sync.WaitGroup：等待协程完成

### 类型原型

```go
type WaitGroup struct {
    // 内部包含计数器 + 信号量（不导出）
}
```

### 方法签名

```go
func (wg *WaitGroup) Add(delta int)   // 计数器 +delta（可正可负）
func (wg *WaitGroup) Done()           // 等价于 Add(-1)
func (wg *WaitGroup) Wait()           // 阻塞直到计数器归零
```

| 方法 | 参数 | 返回值 | 行为 |
|------|------|--------|------|
| `Add(delta int)` | `delta`: 正数增加计数，负数减少 | 无 | 原子地修改内部计数器 |
| `Done()` | 无 | 无 | 计数器 -1，等价 `Add(-1)` |
| `Wait()` | 无 | 无 | 阻塞当前 goroutine，直到计数器 == 0 |

### 标准使用模式

```go
var wg sync.WaitGroup
const workers = 5

for i := 1; i <= workers; i++ {
    wg.Add(1)                       // ① 启动前登记：又有一个要等的
    go func(id int) {
        defer wg.Done()             // ② defer 保证不管怎么退出都 -1
        // 干活...
    }(i)
}

wg.Wait()                           // ③ 阻塞，直到所有 worker Done
// 全部完成，继续执行
```

```
主 goroutine:
  wg.Add(1) x N  → 计数器 = N
  wg.Wait()      → 阻塞等待...
                    ← 计数器逐个 -1，直到归零
                    ← Wait() 返回，继续执行

每个子 goroutine:
  defer wg.Done() → 退出时计数器 -1
```

### 注意事项

- **必须传指针**——`func(wg *sync.WaitGroup)`，值拷贝会导致每个 goroutine 操作的是独立的副本，主 goroutine 的 `Wait()` 永远等不到
- **Add 必须在 goroutine 外部调用**——如果把 `wg.Add(1)` 写在 goroutine 里，可能出现主 goroutine 已经跑到 `Wait()` 时 Add 还没执行，计数器为 0 直接返回
- **Done 必须用 defer**——防止函数中途 return 或 panic 导致 Done 没执行，Wait 永远等不到
- **计数器不能为负数**——`Add(-1)` 过多会 panic
- **WaitGroup 可以复用**——上一轮 Wait 返回后，可以再 Add 开始新一轮

### 选型：WaitGroup vs channel

| 场景 | 推荐 |
|------|------|
| 单纯等 N 个 goroutine 完成 | `WaitGroup`，语义更直接 |
| 需要在协程间**传数据**同时同步 | `channel` |
| 生产者-消费者模型（协程数量动态变化） | `channel` |
| 限流/令牌桶（控制并发数量） | 有缓冲 `channel` |

---

## 2. sync.Mutex：互斥锁

### 类型原型

```go
type Mutex struct {
    state int32   // 锁状态：是否锁定、是否有等待者、是否饥饿
    sema  uint32  // 信号量，用于阻塞/唤醒 goroutine
}
```

### 方法签名

```go
func (m *Mutex) Lock()     // 获取锁；若已被别人持有则阻塞等待
func (m *Mutex) Unlock()   // 释放锁；若当前未锁定则 panic
func (m *Mutex) TryLock() bool  // Go 1.18+ 非阻塞尝试获取；拿到 true，否则 false
```

| 方法 | 参数 | 返回值 | 行为 |
|------|------|--------|------|
| `Lock()` | 无 | 无 | 加锁；若锁空闲则获取并返回；否则阻塞等待 |
| `Unlock()` | 无 | 无 | 解锁；若锁未被锁定则 **panic** |
| `TryLock()` | 无 | `bool` | **非阻塞**尝试加锁，成功返回 `true`，失败返回 `false` |

### Lock/Unlock 标准写法

```go
var mu sync.Mutex
var counter int

// ✅ 标准写法：Lock 后立刻 defer Unlock
mu.Lock()
defer mu.Unlock()
counter++

// ✅ 也可用 TryLock 做"拿不到就走"逻辑
if mu.TryLock() {
    defer mu.Unlock()
    // 拿到锁，做临界区操作
} else {
    // 没拿到，走别的逻辑
}
```

### 不加锁的后果：数据竞争

`counter++` 在 CPU 层面是**三条指令**：

```
LOAD  counter → 寄存器     (读)
ADD   寄存器, 1            (改)
STORE 寄存器 → counter     (写)
```

两个 goroutine 同时执行时：

```
goroutine A: LOAD counter(100) → ADD(101) → STORE counter(101)
goroutine B: LOAD counter(100) → ADD(101) → STORE counter(101)
                                  ↑ 覆盖！两次 +1，结果只多了 1
```

### 锁的配对是铁律

```go
// ❌ 忘了 Unlock → 下一个 Lock 永远等不到 → 死锁
mu.Lock()
counter++
// 没有 Unlock！

// ❌ Unlock 不在 defer 里，中间 return 导致 Unlock 没执行
mu.Lock()
if something {
    return  // ← 提前返回，Unlock 跳过了！
}
mu.Unlock()

// ❌ 对未锁定的 Mutex 调用 Unlock → panic
mu.Unlock()  // panic: sync: unlock of unlocked mutex
```

**永远用 `Lock()` 后紧跟 `defer Unlock()`**。

### 锁的值拷贝陷阱

```go
func bad(mu sync.Mutex) {   // ❌ 值拷贝！
    mu.Lock()               // 锁的是副本的 mu
    defer mu.Unlock()       // 释放的也是副本的
}

var mu sync.Mutex
mu.Lock()                    // 锁的是原版 mu
bad(mu)                      // 把 mu 的状态拷贝进去
mu.Unlock()                  // 释放原版
// → bad 里的 Lock 永远等不到
```

```
原版 mu（栈上）：                  副本 mu（函数参数）：
┌──────────────┐                 ┌──────────────┐
│ state: 1     │  Lock() 之后    │ state: 1     │  ← 值拷贝！锁状态被复制
│ （已加锁）   │ ── 值拷贝 ──→  │ （显示已加锁）│
└──────────────┘                 └──────────────┘

副本的 Lock() 检查副本 state → 已锁 → 等待副本 Unlock()
但原版的 Unlock() 只影响原版 → 副本永远等不到 → 死锁
```

**`sync.Mutex` 有 `noCopy` 标记，`go vet` 能检测值拷贝。永远传指针**：

```go
func good(mu *sync.Mutex) {  // ✅ 传指针
    mu.Lock()
    defer mu.Unlock()
}
```

### 底层原理简要

Mutex 分**正常模式**和**饥饿模式**。正常模式下，新到达的 goroutine 通过自旋（spin）尝试抢锁——如果抢到了，等待队列里的 goroutine 继续等。当等待者等太久（>1ms），Mutex 切换到饥饿模式，直接把锁交给队首的等待者，防止尾部 goroutine 饿死。

更详细的底层原理见 [并发控制底层原理文档](./03a-concurrency-os-principles.md) 第 3 节。

---

## 3. sync.RWMutex：读写锁

### 类型原型

```go
type RWMutex struct {
    w           Mutex   // 互斥锁，保护写者
    writerSem   uint32  // 写者等待信号量
    readerSem   uint32  // 读者等待信号量
    readerCount int32   // 当前持有读锁的数量（为正）或写者等待时的负值标记
    readerWait  int32   // 写者等待时，还需要等多少读者释放
}
```

### 方法签名

```go
func (rw *RWMutex) Lock()       // 获取写锁（独占）
func (rw *RWMutex) Unlock()     // 释放写锁
func (rw *RWMutex) RLock()      // 获取读锁（共享）
func (rw *RWMutex) RUnlock()    // 释放读锁
func (rw *RWMutex) TryLock()    bool  // Go 1.18+ 非阻塞尝试获取写锁
func (rw *RWMutex) TryRLock()   bool  // Go 1.18+ 非阻塞尝试获取读锁
```

| 方法 | 参数 | 返回值 | 行为 |
|------|------|--------|------|
| `Lock()` | 无 | 无 | 获取**写锁**；与读锁、写锁都互斥，阻塞等所有读者/写者释放 |
| `Unlock()` | 无 | 无 | 释放写锁；若写锁未被锁定则 panic |
| `RLock()` | 无 | 无 | 获取**读锁**；与读锁可共存，有写锁时阻塞 |
| `RUnlock()` | 无 | 无 | 释放读锁；若读锁未被锁定则 panic |
| `TryLock()` | 无 | `bool` | 非阻塞尝试写锁 |
| `TryRLock()` | 无 | `bool` | 非阻塞尝试读锁 |

### 互斥矩阵

```
           │ 读锁  │ 写锁
───────────┼───────┼───────
  读锁请求 │  ✅   │  ❌    ← 读读共享，读写互斥
  写锁请求 │  ❌   │  ❌    ← 写写互斥，读写互斥
```

### 使用示例

```go
var (
    book string
    rw   sync.RWMutex
)

// 读者：可以很多人同时读
reader := func(id int) {
    rw.RLock()                                 // 读锁：不阻塞其他读者
    fmt.Printf("读者 %d 正在看: %s\n", id, book)
    rw.RUnlock()
}

// 写者：独占
writer := func(id int) {
    rw.Lock()                                  // 写锁：等所有读者/写者释放
    book = fmt.Sprintf("第 %d 版", id)
    rw.Unlock()
}
```

### 运行过程

```
5 个读者先启动，全部拿到读锁（读读共享）：

Reader 1: RLock() ✅
Reader 2: RLock() ✅  ← 和 Reader 1 同时持有
Reader 3: RLock() ✅
Reader 4: RLock() ✅
Reader 5: RLock() ✅

Writer 1: Lock() ⏸   ← 有读者在，阻塞！
Writer 2: Lock() ⏸   ← 阻塞！

Reader 1~5: RUnlock() → 全部释放 → 唤醒 Writer 1
Writer 1: Lock() ✅ → 独占 → Unlock()
Writer 2: Lock() ✅ → 独占 → Unlock()
```

### 写者偏向

Go 的 RWMutex 有**写者偏向**：当写者已经在等待时，新来的读者会被阻塞，防止写者饿死。

```
Reader 1~3 持有读锁 → Writer 1 Lock() 等待
此时 Reader 4 请求 RLock() → ⏸ 阻塞！
（虽然读读共享，但写者在等了，新读者先排队）
```

### 注意事项

- **RLock 和 RUnlock 必须成对**——与 Lock/Unlock 同样的铁律
- **Lock 和 Unlock 也必须成对**
- **不能把读锁升级为写锁**——先 RLock 再 Lock 会死锁（Lock 要等所有读者释放，包括自己）
- **RWMutex 本身比 Mutex 重**——内部有更复杂的状态管理；读写差不多时不如直接用 Mutex

### 选型

| 场景 | 推荐 | 理由 |
|------|------|------|
| 读多写少（配置、缓存） | RWMutex | 读者可以并发，性能好很多 |
| 读写差不多 | Mutex | RWMutex 本身有额外开销 |
| 写多读少 | Mutex | 读锁的开销纯浪费 |

---

## 4. sync.Once：保证只执行一次

### 类型原型

```go
type Once struct {
    done uint32    // 原子标志位：0=未执行，1=已执行
    m    Mutex     // 互斥锁，保证执行期间的排他性
}
```

### 方法签名

```go
func (o *Once) Do(f func())   // 执行 f；多次调用也只执行第一次
```

| 方法 | 参数 | 返回值 | 行为 |
|------|------|--------|------|
| `Do(f func())` | `f`: 无参无返回的函数 | 无 | 保证 `f` 只执行一次；后续调用直接返回 |

### 执行流程

```go
var (
    once   sync.Once
    config string
)

loadConfig := func() {
    time.Sleep(600 * time.Millisecond)
    config = "配置加载完毕"
}

// 10 个 goroutine 同时调用
for i := 1; i <= 10; i++ {
    go func(id int) {
        once.Do(loadConfig)       // 只有第一个 goroutine 真正执行 loadConfig
        fmt.Println(config)       // 其他 9 个等第一个执行完就继续
    }(i)
}
```

```
goroutine 1: once.Do(f) → 拿到锁 → 执行 f → config = "..." → 释放锁
goroutine 2: once.Do(f) → done==true → 直接返回（不执行 f）
goroutine 3: once.Do(f) → done==true → 直接返回（不执行 f）
...
goroutine 10: once.Do(f) → done==true → 直接返回（不执行 f）
```

### 底层保证

`Once` 内部用**原子操作 + 互斥锁**双层检查（double-check locking）：

1. 先用 `atomic.LoadUint32(&o.done)` 快速检查——若已执行，直接返回（零开销）
2. 若未执行，`Lock()` 进入慢路径——拿到锁后再次检查，防止多个 goroutine 同时进入慢路径
3. 执行 `f()`，执行完后 `atomic.StoreUint32(&o.done, 1)` 标记为完成

### 注意事项

**① 嵌套调用同一个 Once 会死锁**

```go
var once sync.Once
once.Do(func() {
    once.Do(func() {   // ❌ 死锁！外层 Do 持有锁，内层 Do 等待外层释放
        fmt.Println("never reach")
    })
})
```

**② f 不能有参数和返回值**

`Do` 接受的函数签名是 `func()`。如果需要传参或拿返回值，用闭包捕获外部变量：

```go
var db *sql.DB
once.Do(func() {
    db, _ = sql.Open("mysql", dsn)   // 闭包捕获 db 变量
})
// 用 db...
```

**③ 如果 f 执行过程中 panic，Once 仍标记为已执行**

后续调用不会重试 f——如果初始化失败了，Once 不会给你第二次机会。需要重试的场景，加一层自己的错误处理：

```go
var once sync.Once
var initErr error

func Init() error {
    once.Do(func() {
        initErr = doInit()   // 把错误存到外部变量
    })
    return initErr
}
```

**④ f 执行期间，后续调用阻塞**

所有调用 `once.Do(f)` 的 goroutine 都会等第一个执行完才返回——不是"后面的直接跳过"，而是"后面的等前面执行完，然后跳过"。所以如果 f 很慢，所有调用方都会卡住。

---

## 5. sync.Map：并发安全的 Map

### 类型原型

```go
type Map struct {
    mu     Mutex           // 保护 dirty map 的互斥锁
    read   atomic.Value    // 只读 map（无锁读）
    dirty  map[any]*entry  // 需要加锁访问的脏 map
    misses int             // 从 read 未命中次数，触发 dirty 晋升
}
```

### 方法签名

```go
func (m *Map) Store(key, value any)           // 存储键值对
func (m *Map) Load(key any) (value any, ok bool)  // 读取；ok 表示 key 是否存在
func (m *Map) LoadOrStore(key, value any) (actual any, loaded bool)  // 读或存
func (m *Map) LoadAndDelete(key any) (value any, loaded bool)        // 读后删除（Go 1.15+）
func (m *Map) Delete(key any)                 // 删除键值对
func (m *Map) Swap(key, value any) (previous any, loaded bool)       // 交换值（Go 1.20+）
func (m *Map) CompareAndSwap(key, old, new any) (swapped bool)       // CAS（Go 1.20+）
func (m *Map) CompareAndDelete(key, old any) (deleted bool)          // 比较并删除（Go 1.20+）
func (m *Map) Range(f func(key, value any) bool)  // 遍历
```

| 方法 | 参数 | 返回值 | 行为 |
|------|------|--------|------|
| `Store(key, value)` | `key, value`: 任意类型 | 无 | 存储；覆盖已有值 |
| `Load(key)` | `key`: 任意类型 | `(value any, ok bool)` | 读取；`ok=false` 表示 key 不存在 |
| `LoadOrStore(key, value)` | `key, value`: 任意类型 | `(actual any, loaded bool)` | key 存在返回已有值+`loaded=true`；不存在则存储+`loaded=false` |
| `LoadAndDelete(key)` | `key`: 任意类型 | `(value any, loaded bool)` | 读取并删除；不存在时 `loaded=false` |
| `Delete(key)` | `key`: 任意类型 | 无 | 删除；key 不存在也不报错 |
| `Swap(key, value)` | `key, value`: 任意类型 | `(previous any, loaded bool)` | 设置新值返回旧值；
| `CompareAndSwap(key, old, new)` | `key, old, new`: 任意类型 | `swapped bool` | 仅当当前值==old 时替换为 new |
| `CompareAndDelete(key, old)` | `key, old`: 任意类型 | `deleted bool` | 仅当当前值==old 时删除 |
| `Range(f)` | `f`: `func(key, value any) bool` | 无 | 遍历；`f` 返回 `false` 时提前终止 |

### 使用示例

```go
var sm sync.Map

// 存储
sm.Store("🌏", "Earth")
sm.Store("🌙", "Moon")

// 读取
if v, ok := sm.Load("🌏"); ok {
    fmt.Println(v)  // Earth
}

// 有就取，没有就存
v, loaded := sm.LoadOrStore("🪐", "Saturn")
// loaded=false: 之前不存在，Saturn 存进去了
// loaded=true:  之前存在，v 是已有值（Saturn 被忽略）

// 遍历
sm.Range(func(key, value any) bool {
    fmt.Printf("%v → %v\n", key, value)
    return true   // 返回 false 会提前终止遍历
})

// 删除
sm.Delete("🌙")

// 读后删除（Go 1.15+）
v, loaded := sm.LoadAndDelete("🌏")

// 交换（Go 1.20+）
old, loaded := sm.Swap("☀️", "Sol")
```

### 内部设计：读写分离

```
sync.Map 内部有两层 map：

┌─────────────────────────────────┐
│  read (atomic.Value)            │  ← 只读 map，无锁访问
│  {🌏→Earth, 🌙→Moon, ...}       │     大部分 Load 命中这里，零开销
└─────────────────────────────────┘
            ↑ miss 次数够了就晋升
┌─────────────────────────────────┐
│  dirty (map[any]*entry)         │  ← 脏 map，需要 mu 保护
│  {☀️→Sun, 🪐→Saturn, ...}       │     Store 写入这里
└─────────────────────────────────┘
```

- `Load` 先查 read（无锁），命中就直接返回；miss 多了触发 dirty 晋升为新的 read
- `Store` 写入 dirty（需要加锁）
- 这个设计使得 sync.Map 在**读多写少**场景下性能远超 Mutex+map

### 适用场景（官方文档原文）

sync.Map 针对两种场景优化：

1. **key 只写一次但读很多次**（如缓存）
2. **多个 goroutine 读写不同的 key**（减少锁竞争）

### 注意事项

**① 没有 Len() 方法**

```go
// ❌ sync.Map 没有 Len()
// ✅ 用 Range 自己数
count := 0
sm.Range(func(k, v any) bool { count++; return true })
```

**② Key/Value 是 `any`（`interface{}`），丢失类型安全**

```go
v, _ := sm.Load("key")
s := v.(string)  // 必须做类型断言，断错了会 panic
```

**③ Range 不保证强一致性**

遍历过程中有并发写入，可能看到一部分新数据、一部分旧数据。不保证"某个时间点的快照"。

**④ Range 回调返回 false 提前终止**

```go
sm.Range(func(k, v any) bool {
    if k == "stop" {
        return false  // 停止遍历，类似 break
    }
    return true       // 继续遍历
})
```

**⑤ 不适合频繁读写同一 key**

sync.Map 的读写分离设计在"不同 key 被不同 goroutine 操作"时表现好；同一个 key 被频繁读写时，不如 Mutex+map。

### 选型

| 场景 | 推荐 | 理由 |
|------|------|------|
| 写入后读很多次（缓存） | sync.Map | 读操作几乎不用锁 |
| 多个 goroutine 读写不同的 key | sync.Map | 减少锁竞争 |
| 频繁读写同一个 key | Mutex + map | sync.Map 的额外开销反而更大 |
| 非并发场景 | 原生 map | 没有同步开销，最快 |
| 需要 Len() / 类型安全 | Mutex + map | sync.Map 没有这些 |

---

## 6. sync.Cond：条件变量

### 类型原型

```go
type Cond struct {
    noCopy noCopy        // go vet 检测值拷贝
    L      Locker        // 关联的锁（通常是 *Mutex 或 *RWMutex）
    notify notifyList    // 等待 goroutine 列表
    checker copyChecker  // 运行时检测值拷贝
}
```

### 构造函数

```go
func NewCond(l Locker) *Cond   // 创建 Cond，传入一个已存在的锁
```

`Locker` 接口定义：

```go
type Locker interface {
    Lock()
    Unlock()
}
```

所以 `*sync.Mutex` 和 `*sync.RWMutex` 都可以传入，但绝大多数情况用 `*sync.Mutex`。

### 方法签名

```go
func (c *Cond) Wait()        // 原子地解锁 → 挂起 → 被唤醒后重新加锁
func (c *Cond) Signal()      // 唤醒一个等待的 goroutine
func (c *Cond) Broadcast()   // 唤醒所有等待的 goroutine
```

| 方法 | 参数 | 返回值 | 行为 |
|------|------|--------|------|
| `Wait()` | 无 | 无 | **调用前必须持有 c.L 锁**。原子的三个动作：释放锁 → 挂起等待 → 被唤醒→重新获取锁 |
| `Signal()` | 无 | 无 | 唤醒**一个**等待中的 goroutine（选哪个不保证） |
| `Broadcast()` | 无 | 无 | 唤醒**所有**等待中的 goroutine |

### Wait 的三步原子操作

```
goroutine 调用 cond.Wait() 的过程：

① 释放 cond.L 锁         ← 让别的 goroutine 有机会改条件
② 挂起当前 goroutine      ← 进入等待队列，睡着了
③ 被 Signal/Broadcast 唤醒 → 重新拿锁 → Wait() 返回

这三步对调用方来说是一个原子操作——你不会看到中间状态
```

### 标准使用模式：发令枪

```go
var (
    mu    sync.Mutex
    ready bool
    cond  = sync.NewCond(&mu)
)

// 等待方
go func() {
    cond.L.Lock()              // ① 先拿到锁
    for !ready {               // ② 用 for 不是 if！防止虚假唤醒
        cond.Wait()            // ③ 释放锁+挂起，被唤醒后重新拿锁
    }
    cond.L.Unlock()            // ④ 退出
    // 条件满足了，干活！
}()

// 通知方
cond.L.Lock()
ready = true                   // 改变条件
cond.L.Unlock()
cond.Broadcast()               // 叫醒所有等待者
```

```
选手 1: Lock() → ready==false → Wait() → 释放锁，睡觉
选手 2: Lock() → ready==false → Wait() → 释放锁，睡觉
选手 3: Lock() → ready==false → Wait() → 释放锁，睡觉
...
裁判:   Lock() → ready=true → Unlock() → Broadcast() 🔫
选手 1: 被唤醒 → 重新拿锁 → for 检查 ready==true → 跳出循环 → Unlock() → 冲！
选手 2: 被唤醒 → 重新拿锁 → for 检查 ready==true → 跳出循环 → Unlock() → 冲！
选手 3: 被唤醒 → 重新拿锁 → for 检查 ready==true → 跳出循环 → Unlock() → 冲！
```

### 注意事项

**① Wait 必须用 for 不能用 if：防止虚假唤醒**

```go
// ❌ 错误写法
cond.L.Lock()
if !ready {
    cond.Wait()  // 虚假唤醒后直接继续，ready 可能还是 false！
}
cond.L.Unlock()

// ✅ 正确写法
cond.L.Lock()
for !ready {
    cond.Wait()  // 被唤醒后重新检查条件，不满足就继续等
}
cond.L.Unlock()
```

**虚假唤醒（spurious wakeup）** 是操作系统层面的现象——线程可能在没有收到明确信号时被唤醒。Go 的 runtime 做了一些防护，但官方文档仍然明确要求用 `for` 循环。

**② Signal vs Broadcast 的选择**

```go
cond.Signal()    // 叫醒 1 个——适合"只有一个能干活"的场景（如单消费者）
cond.Broadcast() // 叫醒全部——适合"所有人都可以干活"的场景（如发令枪）
```

**③ Wait 前必须持有锁**

```go
// ❌ 没拿锁就 Wait → panic
cond.Wait()  // panic: sync: Wait on unheld Cond

// ✅ 先 Lock 再 Wait
cond.L.Lock()
cond.Wait()
```

**④ 不能在持有锁的情况下 Broadcast**

Broadcast 可以在持有锁时调用，也可以不持有。但通常推荐先解锁再 Broadcast，避免"惊群"时被唤醒的 goroutine 立即堵塞在 Lock 上：

```go
// ✅ 推荐：先解锁，再广播
cond.L.Lock()
ready = true
cond.L.Unlock()
cond.Broadcast()  // 被唤醒的 goroutine 可以直接拿到锁

// ⚠️ 可以但不推荐：持有锁广播
cond.L.Lock()
ready = true
cond.Broadcast()  // 所有被唤醒的 goroutine 会立刻堵在 Lock 上
cond.L.Unlock()   // 此时才释放
```

---

## 7. sync.Pool：临时对象池

### 类型原型

```go
type Pool struct {
    noCopy noCopy      // go vet 检测值拷贝
    local     unsafe.Pointer  // 每个 P 的本地池（无锁访问）
    localSize uintptr
    victim     unsafe.Pointer  // 上一轮 GC 的幸存池（GC 时 local → victim）
    victimSize uintptr
    New       func() any       // 池空时创建新对象的工厂函数
}
```

### 方法签名

```go
func (p *Pool) Get() any    // 从池中取一个对象；池空则调用 New（若 nil 则返回 nil）
func (p *Pool) Put(x any)   // 把对象放回池中
```

| 方法 | 参数 | 返回值 | 行为 |
|------|------|--------|------|
| `Get()` | 无 | `any`（需类型断言） | 从池取对象；池空时调用 `New()` 创建；`New` 为 nil 时返回 `nil` |
| `Put(x any)` | `x`: 要放回的对象 | 无 | 放回池中供后续复用；**不能放 nil** |

### 使用示例

```go
type Bullet struct {
    ID   int
    Used bool
}

var bulletID int32

pool := sync.Pool{
    New: func() any {                     // 工厂函数：池空时调用
        id := atomic.AddInt32(&bulletID, 1)
        return &Bullet{ID: int(id)}
    },
}

// 借
b1 := pool.Get().(*Bullet)               // 池空 → 调用 New
b2 := pool.Get().(*Bullet)               // 池空 → 调用 New

// 还
pool.Put(b1)                              // b1 回到池中

// 再借
b3 := pool.Get().(*Bullet)               // 池中有 b1 → 直接返回，不用 New！
```

```
第一次 Get：
  检查本地池 → 空 → 检查 victim → 空 → 调用 New() → 创建 #1 → 返回

Put(b1)：
  b1 放回本地池 → 池：[#1]

第二次 Get：
  检查本地池 → 有 #1 → 直接取出返回（内存复用！）
```

### GC 对 Pool 的影响：两轮回收机制

sync.Pool 的**关键特性**：GC 时池内对象会被回收。

```
GC 第 N 轮之前：
  local:  [#1, #2, #3]    ← 当前活跃对象
  victim: []               ← 空

GC 第 N 轮发生：
  victim = local           ← local 降级为 victim
  local = []                ← local 清空

GC 第 N 轮之后到第 N+1 轮之前：
  local:  []               ← 新一轮 Get 会 New 出新对象
  victim: [#1, #2, #3]     ← 上一轮的幸存者，还能再用一次！

GC 第 N+1 轮发生：
  victim = []               ← 上一轮的 victim 彻底回收
```

Get 的执行顺序：先查 local → 没命中查 victim → 还没就调用 New。

这意味着：对象至少能存活一轮 GC（一个 GC 周期内反复复用），跨两轮 GC 后被回收。

### 注意事项

**① Get 拿到的对象状态不确定——必须自己初始化关键字段**

```go
st := pool.Get().(*Bullet)
// st 可能是：
//   - New 刚创建的（ID 是新分配的，Used=false）
//   - 别人 Put 回来的（ID 和 Used 是上一次的值）

// ✅ 正确：Get 后重置关键字段
st.Used = false
// 不要依赖 Get 拿到"干净的默认值"
```

**② Put 之后不能再碰这个对象**

```go
st := pool.Get().(*Bullet)
pool.Put(st)
// ❌ Put 之后不要再访问 st！它可能已被另一个 goroutine Get 走
fmt.Println(st.ID)  // 数据竞争！
```

Put 意味着放弃对象所有权——它现在是池的财产。

**③ Pool 不能当持久缓存**

持久缓存用 `map` + `Mutex`。Pool 只能存"丢了也无所谓"的临时对象。

**④ Pool 本身是并发安全的**

多个 goroutine 可以同时 `Get` 和 `Put`，无需额外加锁。Pool 内部为每个 P 分配了独立的本地池，大多数 Get/Put 操作是无锁的。

**⑤ New 字段可以为 nil**

如果不设置 `New`，`Get()` 在池空时返回 `nil`。这种用法适合"有就有，没有就算了"的场景（但很少见）。

### 典型应用场景

| 场景 | 例子 | 为什么适合 Pool |
|------|------|----------------|
| 网络/IO 缓冲区 | `[]byte` 频繁分配和释放 | 每次读写都分配新 buffer，GC 压力巨大 |
| 字符串拼接 | `bytes.Buffer` 的复用 | `fmt` 包内部就用 Pool 管理 Buffer |
| 序列化中间对象 | JSON/Protobuf 编解码的临时结构体 | 每个请求都 new→用完→GC，复用后 GC 压力骤降 |
| 高频临时对象 | 秒杀系统中的订单结构体 | 高并发下减少 GC 暂停 |

---

## 8. sync/atomic：无锁原子操作

### 和 sync.Mutex 的本质区别

| 维度 | atomic | Mutex |
|------|--------|-------|
| 实现层级 | CPU 硬件指令（LOCK 前缀） | OS 调度器（信号量、睡眠/唤醒） |
| 操作范围 | 单个整数/指针的简单操作 | 任意代码块 |
| 阻塞 | **不阻塞**，立即完成 | 等不到锁时**阻塞** |
| 性能 | 几个 CPU 周期 | 几百个 CPU 周期起 |
| 适用场景 | 简单计数器、状态标志、指针交换 | 复杂逻辑、多变量事务 |

### 整数原子操作

支持类型：`int32`, `int64`, `uint32`, `uint64`, `uintptr`, `unsafe.Pointer`

```go
// Add —— 原子加
func AddInt32(addr *int32, delta int32) (new int32)
func AddInt64(addr *int64, delta int64) (new int64)
func AddUint32(addr *uint32, delta uint32) (new uint32)
func AddUint64(addr *uint64, delta uint64) (new uint64)
func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)

// Load —— 原子读
func LoadInt32(addr *int32) (val int32)
func LoadInt64(addr *int64) (val int64)
func LoadUint32(addr *uint32) (val uint32)
func LoadUint64(addr *uint64) (val uint64)
func LoadUintptr(addr *uintptr) (val uintptr)
func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer)

// Store —— 原子写
func StoreInt32(addr *int32, val int32)
func StoreInt64(addr *int64, val int64)
func StoreUint32(addr *uint32, val uint32)
func StoreUint64(addr *uint64, val uint64)
func StoreUintptr(addr *uintptr, val uintptr)
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)

// Swap —— 原子交换（返回旧值）
func SwapInt32(addr *int32, new int32) (old int32)
func SwapInt64(addr *int64, new int64) (old int64)
func SwapUint32(addr *uint32, new uint32) (old uint32)
func SwapUint64(addr *uint64, new uint64) (old uint64)
func SwapUintptr(addr *uintptr, new uintptr) (old uintptr)
func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)

// CompareAndSwap —— CAS 乐观锁
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)
func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)
func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool)
func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)
func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)
```

### 使用示例

```go
var counter int64

// Add：原子加
atomic.AddInt64(&counter, 1)   // 并发安全地 counter++

// Load + Store：原子读写
var flag int32
atomic.StoreInt32(&flag, 42)         // 写入
val := atomic.LoadInt32(&flag)       // 读出 → 42

// Swap：原子交换，返回旧值
old := atomic.SwapInt32(&flag, 100)  // flag 变成 100，old = 42

// CAS：如果当前值 == 期望值，就替换为新值
swapped := atomic.CompareAndSwapInt32(&flag, 100, 200)
// swapped = true（flag 是 100，等于期望值 → 换成 200）

swapped = atomic.CompareAndSwapInt32(&flag, 999, 300)
// swapped = false（flag 是 200，不等于 999 → 不换）
```

### CAS 的 ABA 问题

```
时刻 1: goroutine A 读到 flag = 100
时刻 2: goroutine B 把 flag 改成 999
时刻 3: goroutine C 把 flag 改回 100
时刻 4: goroutine A 执行 CAS(flag, 100, 200) → 成功！

A 以为 flag 从 100 变成了 200，中间没被人动过
但实际上 flag 经历了 100 → 999 → 100 的变化
```

CAS 只看"当前值是不是期望值"，不关心中间是否变化过。在指针场景下这可能很危险——指针指向的对象可能已被回收。

**缓解方案：带版本号的 tagged pointer**（高位存版本号，低位存指针），确保版本变化也能被检测到。

### atomic.Value：原子地存取任意类型

```go
type Value struct {
    // 内部字段（不导出）
}
```

**方法签名：**

```go
func (v *Value) Load() (val any)                                // 原子加载
func (v *Value) Store(val any)                                  // 原子存储
func (v *Value) Swap(new any) (old any)                         // 原子交换（Go 1.17+）
func (v *Value) CompareAndSwap(old, new any) (swapped bool)     // CAS（Go 1.17+）
```

| 方法 | 参数 | 返回值 | 行为 |
|------|------|--------|------|
| `Load()` | 无 | `any` | 原子读取；若从未 Store 则返回 nil |
| `Store(val)` | `val`: **不能为 nil** | 无 | 原子写入 |
| `Swap(new)` | `new`: 新值 | `old any` | 设置新值，返回旧值 |
| `CompareAndSwap(old, new)` | `old, new`: 比较值和替换值 | `swapped bool` | 仅当当前值 == old 时替换 |

**使用示例：**

```go
type Config struct{ DB string }
var cfg atomic.Value

// 存储
cfg.Store(&Config{DB: "mysql://localhost"})

// 加载（需要类型断言）
c := cfg.Load().(*Config)
fmt.Println(c.DB)  // mysql://localhost
```

**注意事项：**

**① 类型必须一致——第一次 Store 就锁定了类型**

```go
var v atomic.Value
v.Store(1)
v.Store("hello")  // ❌ panic！之前存的是 int，不能换成 string
```

**② Store 不能传 nil**

```go
var v atomic.Value
v.Store(nil)  // ❌ panic
```

**③ 适合读多写少**

`atomic.Value` 的 Load 操作极快，但 Store 需要分配内存（interface{} 装箱）。如果写很频繁，考虑用 Mutex 保护。

**④ 并发写同一个 Value 也不会有数据竞争**

多个 goroutine 同时 Store 会内部排队，不会出现"一半新一半旧"的数据。

---

## 9. 综合案例：限量秒杀系统

展示了多个原语协同工作的真实模式：

```go
var (
    stock   int32 = 10        // atomic 保护：无锁扣库存
    success int32             // atomic 保护：统计成交
    wg      sync.WaitGroup    // 等所有用户抢完
    once    sync.Once         // 初始化只执行一次
)

const users = 100

initSale := func() { /* 初始化，只跑一次 */ }

for i := 1; i <= users; i++ {
    wg.Add(1)
    go func(uid int) {
        defer wg.Done()
        once.Do(initSale)                // Once：初始化只跑一次

        for {
            cur := atomic.LoadInt32(&stock)
            if cur <= 0 { return }       // 库存没了，退出

            if atomic.CompareAndSwapInt32(&stock, cur, cur-1) {
                atomic.AddInt32(&success, 1)  // CAS 成功 = 抢到了
                return
            }
            // CAS 失败 → 库存被别人改了，重试
        }
    }(i)
}

wg.Wait()  // 等所有人抢完
```

```
技术栈使用：
  sync.Once   → 初始化库存（保证只执行一次）
  sync.WaitGroup → 等所有用户抢完
  atomic.CAS  → 乐观锁扣库存（比 Mutex 吞吐量高得多）
  atomic.Add  → 并发安全地统计成功数
```

**为什么用 CAS 而非 Mutex？**

Mutex 在高竞争下，每个 goroutine 都要排队等锁 → 大量上下文切换。CAS 失败了立即重试，CPU 一直在干活——100 人抢 10 件库存的场景，CAS 吞吐量远超 Mutex。

---

## sync 包工具速查

| 类型 | 关键方法（含原型） | 用途 | 核心注意事项 |
|------|-------------------|------|-------------|
| `sync.WaitGroup` | `Add(delta int)`, `Done()`, `Wait()` | 等一组 goroutine 完成 | 必须传指针；Add 在 goroutine 外部 |
| `sync.Once` | `Do(f func())` | 只执行一次（延迟初始化） | 不要嵌套调同一个 Once；f 不能有参/返回值 |
| `sync.Mutex` | `Lock()`, `Unlock()`, `TryLock() bool` | 互斥保护临界区 | 不能值拷贝；Lock 后立刻 defer Unlock |
| `sync.RWMutex` | `RLock()`, `RUnlock()`, `Lock()`, `Unlock()`, `TryLock() bool`, `TryRLock() bool` | 读写锁，读多写少 | 读读共享，读写互斥；不能读锁升级为写锁 |
| `sync.Map` | `Store(k,v)`, `Load(k) (v,ok)`, `Delete(k)`, `LoadOrStore(k,v) (v,loaded)`, `LoadAndDelete(k) (v,loaded)`, `Swap(k,v) (prev,loaded)`, `CompareAndSwap(k,old,new) bool`, `CompareAndDelete(k,old) bool`, `Range(f)` | 并发安全的 map | 无 Len()；K/V 是 any；适合读多写少 |
| `sync.Cond` | `NewCond(l Locker) *Cond`, `Wait()`, `Signal()`, `Broadcast()` | 条件变量，等信号再干活 | Wait 必须用 for 不是 if；Wait 前必须持有锁 |
| `sync.Pool` | `Get() any`, `Put(x any)`, `New func() any` | 临时对象复用，减少 GC | GC 会清空池；Get 后初始化字段；Put 后别碰对象 |
| `atomic.Add*` | `AddInt32/64(addr, delta) new`, `LoadInt32/64(addr) val`, `StoreInt32/64(addr, val)`, `SwapInt32/64(addr, new) old`, `CompareAndSwapInt32/64(addr, old, new) bool` | 原子操作整数 | 只适合单变量；无阻塞；CAS 有 ABA 问题 |
| `atomic.Value` | `Load() any`, `Store(val)`, `Swap(new) old`, `CompareAndSwap(old,new) bool` | 原子存取结构体/任意类型 | 类型必须一致（首次 Store 锁定）；Store 不能 nil |

---

## 易错点

1. **Mutex 值拷贝死锁**——传值不传指针，锁的内部状态被复制，副本的 Lock 永远等不到 → 死锁。Mutex 有 `noCopy` 标记，`go vet` 能检测到
2. **Lock/Unlock 不成对**——忘了写 Unlock 或提前 return 跳过，下一个 Lock 永远等不到。**永远 `Lock()` 后立刻 `defer Unlock()`**
3. **原生 map 并发写直接 fatal**——不是数据错，是程序崩。读也要加锁（读和写互斥）
4. **循环等待死锁**——两个 goroutine 按不同顺序加锁（A: mu1→mu2，B: mu2→mu1），形成等待环。**统一加锁顺序**可避免
5. **once.Do 嵌套死锁**——在 once.Do 的回调里再调同一个 once 的 Do，外层持有锁等内层，内层等外层释放 → 死锁
6. **atomic.Value 类型不一致会 panic**——第一次 Store 就固定了类型，后续 Store 不同类型会 panic；Store(nil) 也会 panic
7. **WaitGroup 的 Add 写在 goroutine 里**——可能出现主 goroutine 已经 Wait 了而 Add 还没执行，计数器为 0 直接返回
8. **Cond.Wait 用 if 不用 for**——虚假唤醒后不再检查条件，可能在不满足条件时继续执行
9. **Pool 的对象会被 GC 回收**——不能假设 Put 之后 Get 一定拿到同一个对象；Pool 只能存"丢了也无所谓"的临时对象
10. **Pool.Get 拿到脏对象**——Get 返回的对象可能是别人 Put 回来改过的，状态不确定。Get 后必须自己初始化关键字段
11. **Pool.Put 后继续用对象 → 数据竞争**——Put 后对象所有权归池，可能被其他 goroutine Get 走，继续访问 = 并发读写同一块内存
12. **RWMutex 不能读锁升级为写锁**——先 RLock 再 Lock 会死锁，因为 Lock 要等所有读者释放（包括自己）

---

## 快问快答

### Q1：`sync.Mutex` 为什么不能值拷贝？

Mutex 内部有状态字段（state, sema），值拷贝会把这些状态完整复制给副本。如果原版已经 Lock，副本里的状态也是"已锁定"——但副本的 Lock 会等副本的 Unlock，没人会来 Unlock 副本，于是死锁。永远传指针。

### Q2：读写锁的"读读共享"具体怎么理解？

如果当前只有读锁被持有（没有被写锁），新的 RLock 会立刻成功，不阻塞——因为读操作不修改数据，多个读者各读各的互不干扰。当有写锁请求时，新的读锁请求会被阻塞（防止写者饿死）。

### Q3：什么时候用 sync.Map，什么时候用 Mutex + map？

sync.Map 适合"写入一次、读取多次"的缓存场景，或者不同 key 被不同 goroutine 操作的场景——内部做了读写分离，读操作几乎不用锁。大多数情况下 Mutex + map 更简单且保留了类型安全。非并发场景直接用原生 map，最快。

### Q4：atomic 和 Mutex 的本质区别是什么？

atomic 是 CPU 硬件指令（LOCK 前缀）保证的，执行不被中断，不涉及 OS 调度器，不阻塞——适合单个变量的简单操作。Mutex 是 OS 调度器管理的，等不到锁就让出 CPU 去睡觉，被唤醒后再抢——可以保护任意代码块。一句话：atomic 是轻量快刀，Mutex 是万能重器。

### Q5：`wg.Add()` 为什么要写在 goroutine 外部？

如果把 `wg.Add(1)` 写在 goroutine 里面，可能出现这样的顺序：主 goroutine 创建了子 goroutine 但子 goroutine 还没来得及执行 `wg.Add(1)`，主 goroutine 已经执行到了 `wg.Wait()`——此时计数器是 0，Wait 直接返回，子 goroutine 还没跑。外部先 Add 保证 Wait 之前计数器已经加上去了。

### Q6：Go 的 Mutex 是可重入的吗？

不可重入。同一个 goroutine 如果 Lock 了两次，第二次 Lock 会死锁——因为锁已经被自己持有了，再次 Lock 会等自己 Unlock，而第二次 Lock 之后的代码里才有 Unlock。这是 Go 团队的设计选择：可重入锁虽然方便，但会掩盖设计问题。

### Q7：为什么原生 map 并发写会 fatal 而不是静默出错？

这是 Go 运行时的设计哲学——宁可直接崩掉让你看到，也别带着数据竞争上线。并发写 map 会导致内存损坏、程序在完全无关的地方崩溃，极难排查。运行时检测到后主动 fatal，告诉你是 map 并发写的问题。

### Q8：CAS 的 ABA 问题是什么意思？

线程 A 用 CAS 检查某个值，发现是 X（以为没变），于是换成 Y。但实际上值先被线程 B 改成了 Z，又被线程 C 改回了 X——CAS 只看"是不是同一个值"，不知道"中间有没有变化过"。在指针场景下这很危险——地址没变，但指向的对象可能已被回收。解决：带版本号的 tagged pointer。

### Q9：sync.Pool 和普通缓存（map+Mutex）有什么区别？

本质区别在于生命周期和可靠性。sync.Pool 是临时对象池——GC 一来就清空，对象随时可能消失，只能存"丢了也无所谓"的临时对象，优势是零维护成本、并发安全内置。map+Mutex 是持久缓存——你放进去的东西会一直在，GC 不会回收，适合存需要长期保留的数据。

### Q10：Cond 的 Wait 为什么必须用 for 循环？

防止**虚假唤醒**（spurious wakeup）——操作系统可能在没有明确 Signal/Broadcast 的情况下唤醒等待线程。如果用 `if`，被虚假唤醒后不会重新检查条件，可能在不满足条件时继续执行。用 `for` 循环保证唤醒后重新检查条件，不满足就继续等。

### Q11：sync.Pool 里 Get 拿到的对象为什么要自己初始化字段？

因为 Get 可能拿到别人 Put 回来的"脏"对象——字段值是上一次使用留下的，不是 New 创建的"干净"默认值。如果依赖 Get 拿到默认值，就会读到脏数据。正确做法是 Get 后对需要使用的字段自行赋值。

---

## 一句话总结

`sync` 包是 Go 并发安全的工具箱：WaitGroup 等协程、Mutex 保护临界区、RWMutex 优化读多写少、Once 保证只执行一次、Cond 条件变量等信号、Pool 复用临时对象减轻 GC 压力、atomic 无锁原子操作适合简单场景、sync.Map 在特定场景替代 Mutex+map。用对工具、记住坑点（不拷贝锁、不成对 Lock/Unlock、加锁顺序一致、Wait 用 for 不是 if、Pool 对象会被 GC 回收、atomic.Value 类型必须一致），并发代码才能安全高效。
