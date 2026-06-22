# sync 包：并发安全与同步原语

## 这一章要记住什么

- `sync` 包提供了 Go 并发编程的三板斧——**WaitGroup**（等协程完成）、**Mutex/RWMutex**（保护临界区）、**Once**（保证只执行一次）
- **原生 map 不是并发安全的**——多 goroutine 同时写会直接 fatal；小并发用 Mutex 保护，特定场景用 sync.Map
- `sync.Mutex` **绝不能值拷贝**——拷贝会把锁的内部状态一起复制，导致副本里的 Lock 永远等不到解锁 → 死锁
- **原子操作**（`sync/atomic`）比 Mutex 更快，不涉及 OS 调度，适合简单计数器和状态标志
- 锁的获取顺序不一致会导致**循环等待死锁**——所有 goroutine 必须按相同顺序加锁
- 读写锁把读者和写者分开，**读读共享、读写互斥、写写互斥**，适合读多写少的场景
- **sync.Pool 是临时对象池**——Get/Put 复用高频临时对象，减少 GC 压力；但对象随时可能被 GC 回收，不能当持久缓存

---

## 1. channel vs WaitGroup：两种协程同步方式

你的代码先用 `chan struct{}` 实现同步，再用 `sync.WaitGroup`，实际上展示了两种等协程完成的方案。

### 方案一：chan struct{} 信号通知

```go
ch := make(chan struct{}, 10)  // 有缓冲，避免发送方阻塞
for i := 0; i < 10; i++ {
    go func(i int) {
        fmt.Printf("num:%d\n", i)
        ch <- struct{}{}       // 干完活，往 channel 塞一个信号
    }(i)
}
for i := 0; i < 10; i++ {
    <-ch                        // 收到 10 个信号 = 10 个协程都完成了
}
```

```
协程 0: 干活 → ch <- struct{}{}
协程 1: 干活 → ch <- struct{}{}
   ...
主协程: 收 10 次 <-ch → 收够就继续执行
```

`struct{}` 在这里只充当**信号**——我不关心你传什么数据，我只是通知"我搞定了"。`struct{}` 类型的内存占用是 **0 字节**，是 Go 里表达"无数据含义的信号"的最佳选择。

### 方案二：sync.WaitGroup

```go
var wg sync.WaitGroup

myGo := func() {
    defer wg.Done()   // 完成时计数器 -1
    fmt.Println("myGo!")
}

wg.Add(10)             // 先登记：我要等 10 个任务
for i := 0; i < 10; i++ {
    go myGo()
}
wg.Wait()              // 等计数器归零
```

### 什么时候用哪个？

| 场景 | 推荐 |
|------|------|
| 单纯等 N 个 goroutine 完成 | `WaitGroup`，语义更直接 |
| 需要在协程间**传数据**同时同步 | `channel` |
| 生产者-消费者模型（协程数量动态变化） | `channel` |
| 限流/令牌桶（控制并发数量） | 有缓冲 `channel` |

### 总结一下

`chan struct{}` 用 channel 的发送/接收当信号，忙完就发一个，收够就知道全完了。`WaitGroup` 用计数器：Add 登记、Done 递减、Wait 等归零。单纯等协程完成用 WaitGroup，需要带数据用 channel。

---

## 2. sync.Once：保证只执行一次

很多逻辑只需要执行一次——加载配置、初始化连接池、注册路由。如果放在 `init()` 里，程序启动就执行了，不管你用不用。`sync.Once` 实现了**延迟初始化**：第一次用到的时候才初始化，后续调用直接返回。

```go
var instance *Config
var once sync.Once

func InitConfig() *Config {
    once.Do(func() {
        instance = &Config{}  // 这个函数只会执行一次
    })
    return instance
}
```

```
goroutine A: InitConfig() → once.Do(f) → 执行 f，初始化 instance
goroutine B: InitConfig() → once.Do(f) → f 已经执行过，直接跳过
goroutine C: InitConfig() → once.Do(f) → f 已经执行过，直接跳过
...
```

**底层保证：** 即使 100 个 goroutine 同时并发调用 `InitConfig()`，`once.Do` 里的函数也**只会执行一次**，不会有多个 goroutine 各自创建一个 Config 实例。`Once` 内部用了原子操作 + 互斥锁来保证这个语义。

**坑：once.Do 里不能再嵌套调同一个 once**

```go
var once sync.Once
once.Do(func() {
    once.Do(func() {  // ❌ 死锁！内部的 Do 等外部的 Do 完成，外部的等内部的
        fmt.Println("nested")
    })
})
```

`Once` 的 Do 方法内部持有锁，递归调用同一个 Once 会造成死锁。

### 总结一下

`sync.Once` 实现线程安全的延迟初始化——不管多少个 goroutine 并发调用，传入的函数只执行一次。适合配置加载、连接池初始化等场景。不要嵌套调用同一个 Once，会死锁。

---

## 3. sync.Mutex：互斥锁

### 不加锁的后果

```go
num := 0
add := func() { num += 1 }  // 读-改-写，不是原子操作

for i := 0; i < 100000; i++ {
    go add()
}
fmt.Println(num == 100000)  // false！num 一定小于 100000
```

为什么？回头翻 `num += 1` 的 CPU 指令分解：

```
goroutine A: LOAD num(100) → ADD(101) → STORE num(101)
goroutine B: LOAD num(100) → ADD(101) → STORE num(101)  ← A 刚写入的 101 被覆盖

两个 goroutine 各自加了一次，最终 num 只增加了 1，而不是 2
```

10 万次并发自增，有大量这种"覆盖丢失"，最终结果会远小于 10 万。

### 加锁后

```go
var guard sync.Mutex

add := func() {
    guard.Lock()        // 拿到锁才能继续，拿不到就在这等着
    num += 1            // 临界区：同一时刻只有一个 goroutine 在执行
    guard.Unlock()      // 释放锁，让给别人
}
// num 一定等于 100000
```

```
goroutine A: Lock() ✅ → num += 1 → Unlock()
goroutine B: Lock() ⏸  (等 A 释放) → Lock() ✅ → num += 1 → Unlock()
goroutine C: Lock() ⏸  (等 B 释放) → Lock() ✅ → num += 1 → Unlock()
...

同一时刻只有一个人在临界区里，num 的每次 +1 都不会被覆盖
```

### 锁的配对是铁律

```go
// ✅ 标准写法：Lock 和 Unlock 靠在一起
mu.Lock()
defer mu.Unlock()
// 临界区代码...

// ❌ 忘了 Unlock → 下一个 Lock 永远等不到 → 死锁
mu.Lock()
num += 1
// 没有 Unlock！

// ❌ Unlock 不在 defer 里，中间 return/panic 导致 Unlock 没执行
mu.Lock()
if something {
    return  // ← 提前返回，Unlock 没执行！
}
mu.Unlock()
```

**永远用 `defer mu.Unlock()`** 紧跟在 Lock 之后，这是 Go 社区的惯用法。

### 总结一下

`sync.Mutex` 用 Lock/Unlock 保护临界区，保证同一时刻只有一个 goroutine 在执行被保护的代码。Lock 和 Unlock 必须成对出现，最佳实践是 `Lock()` 后立刻写 `defer Unlock()`。底层原理见 [并发控制底层原理文档](./03a-concurrency-os-principles.md) 第 3 节。

---

## 4. sync.RWMutex：读写锁

### 为什么需要读写锁

代码里的读写锁 demo 展示了核心思想：

```go
read := func(mr *sync.RWMutex, i int) {
    mr.RLock()          // 读锁：多个 goroutine 可以同时持有
    fmt.Printf("reading count: %d\n", cnt)
    mr.RUnlock()
}

write := func(mr *sync.RWMutex, i int) {
    mr.Lock()           // 写锁：独占，与读锁和写锁都互斥
    cnt++
    mr.Unlock()
}
```

### 实际运行过程

```
同时启动 3 个 writer goroutine 和 3 个 reader goroutine：

Write 1: Lock() → 拿到写锁 ✅，cnt++，Unlock()
Write 2: Lock() → 拿不到（1 拿着），阻塞 ⏸
Write 3: Lock() → 拿不到（1 拿着），阻塞 ⏸
Read 1: RLock() → 拿不到（1 拿着写锁，读写互斥），阻塞 ⏸
Read 2: RLock() → 同上
Read 3: RLock() → 同上

Write 1: Unlock() → 释放写锁
    → Write 2 被唤醒：Lock() ✅ → cnt++ → Unlock()
    → 注意：writer 可能优先于 reader（取决于 Go 的实现，偏向写者防止写饥饿）
    
... 交错的顺序取决于调度，但核心规则不变
```

### 多读者并发读的示例

```
当前状态：有 3 个 reader 持有读锁

Reader 1: RLock() ✅（第一个读锁）
Reader 2: RLock() ✅（读读不互斥！）
Reader 3: RLock() ✅（读读不互斥！）

Writer 1: Lock()  ⏸ 阻塞！（读写互斥，等到所有读者释放）

Reader 1: RUnlock() → 还剩 2 个读者
Reader 2: RUnlock() → 还剩 1 个
Reader 3: RUnlock() → 全部释放 → 唤醒 Writer 1
```

### 读写锁 vs 互斥锁的选择

| 场景 | 推荐 | 理由 |
|------|------|------|
| 读多写少（如配置、缓存） | RWMutex | 读者可以并发，性能好很多 |
| 读写差不多 | Mutex | RWMutex 本身有额外开销 |
| 写多读少 | Mutex | 读锁的开销浪费了 |

### 总结一下

RWMutex 把锁分成读锁和写锁——读读共享，读写互斥，写写互斥。读多写少的场景性能远超 Mutex。但读写锁本身比 Mutex 重，不适合读写均衡或写多的场景。

---

## 5. 死锁实战分析

你的代码里演示了两种死锁，正好对应底层原理文档里死锁两个最重要的场景。

### 5.1 Lock/Unlock 不成对 + 值拷贝陷阱

```go
func demo4_no_couple_lock() {
    copyMutex := func(mu sync.Mutex) {  // ⚠️ 参数是值传递！
        mu.Lock()
        defer mu.Unlock()
        fmt.Println("ok")
    }

    var mu sync.Mutex
    mu.Lock()          // 给原版 mu 加锁
    defer mu.Unlock()
    copyMutex(mu)      // 把 mu 拷贝给函数 → 锁状态也被复制了
}
```

**发生了什么：**

```
原版 mu（栈上）：                  副本 mu（copyMutex 的参数）：
┌──────────────┐                 ┌──────────────┐
│ state: 1     │  Lock() 之后    │ state: 1     │  ← 值拷贝！锁状态被复制
│ （已加锁）   │ ── 值拷贝 ──→  │ （显示已加锁）│
└──────────────┘                 └──────────────┘

copyMutex 里执行 mu.Lock()：
  检查副本的 state → 已经锁了 → 等待原持有者 Unlock
  
但是！原持有者是"原版 mu"，而副本里的 Lock 等的是"副本的 Unlock"
→ 原版 mu.Unlock() 不会释放副本的锁
→ 副本里的 Lock 永远等不到 → 死锁
```

**关键教训：`sync.Mutex` 绝对不能值拷贝。** Go 的 `go vet` 可以检测到这种问题——`Mutex` 结构体内部有 `noCopy` 标记。如果你需要把锁传到函数里，**永远传指针**：

```go
copyMutex := func(mu *sync.Mutex) {  // ✅ 传指针
    mu.Lock()
    defer mu.Unlock()
    fmt.Println("ok")
}
```

### 5.2 循环等待

你的代码演示的就是上一份原理文档里分析的经典循环等待死锁：

```go
// goroutine A
mu1.Lock()       // 拿到 mu1
time.Sleep(1s)
mu2.Lock()       // 想拿 mu2，但 B 持有 mu2

// goroutine B
mu2.Lock()       // 拿到 mu2
time.Sleep(1s)
mu1.Lock()       // 想拿 mu1，但 A 持有 mu1
```

```
A 持有 mu1，等 mu2
B 持有 mu2，等 mu1
     ↓
  循环等待环 → 死锁
```

**解决方法：统一加锁顺序**

```go
// ✅ 两个 goroutine 都先 Lock mu1 再 Lock mu2
go func() {
    mu1.Lock()
    mu2.Lock()
    // ...
}()

go func() {
    mu1.Lock()  // 先等 mu1
    mu2.Lock()  // 再等 mu2
    // ...
}()
```

只要所有 goroutine 按**相同顺序**加锁，就不会形成环。

### 总结一下

Mutex 值拷贝会把锁状态一起复制，导致副本的 Lock 永远等不到——这是 Go 里最隐蔽的死锁场景之一。循环等待死锁是两个 goroutine 各自持有一把锁、互相等对方的锁；解决办法是统一所有 goroutine 的加锁顺序。

---

## 6. 并发安全的 Map

### 原生 map 不是并发安全的

Go 的内置 `map` 在多 goroutine 同时写时，**不是"结果不正确"的问题，而是直接 fatal**：

```go
var m = make(map[string]int)

for i := 0; i < 10; i++ {
    go func(num int) {
        m[strconv.Itoa(num)] = num  // ❌ 并发写 → fatal
    }(i)
}
// fatal error: concurrent map writes
```

Go 运行时会在检测到并发 map 写时主动抛出 fatal，而不是静默产生错误数据。这是**检测型保护**——宁可直接崩掉让你看到，也不让你带着数据竞争上线。

### 方案一：Mutex + map

```go
var m = make(map[string]int)
var mu sync.Mutex

mu.Lock()
m[key] = value    // 写操作加锁
mu.Unlock()

mu.Lock()
val := m[key]     // 读操作也要加锁！读和写互斥
mu.Unlock()
```

**读也要加锁**——因为写入过程中如果有人在读，可能读到写了一半的不完整数据。

你的代码里 `demo5_lock_map` 有一个**风格上的小问题**：

```go
defer func() {
    wg.Done()
    mu.Unlock()   // Unlock 在 defer 里
}()
// ...
mu.Lock()          // Lock 在外面
```

Lock 和 Unlock 隔着 defer 层的距离，不太容易一眼看出配对关系。更推荐的写法是 Lock 后紧接 defer Unlock：

```go
mu.Lock()
defer mu.Unlock()
// 操作 map...
```

### 方案二：sync.Map（Go 1.9+）

```go
var m sync.Map

m.Store("name", "zhangsan")   // 写入
m.Store("age", 18)

age, ok := m.Load("age")      // 读取，ok 表示 key 是否存在

m.Range(func(key, value interface{}) bool {  // 遍历
    fmt.Printf("key: %v, val: %v\n", key, value)
    return true  // 返回 false 会提前终止遍历
})

m.Delete("age")               // 删除

actual, loaded := m.LoadOrStore("name", "lisi")  // 有就返回已有值，没有就存储
```

### sync.Map 的设计考虑

**为什么不是直接替换原生 map？** sync.Map 不是万能的，它有特定适用场景：

| 场景 | 推荐 | 理由 |
|------|------|------|
| 写入后读很多次（缓存） | sync.Map | 读操作几乎不用锁，很快 |
| 多个 goroutine 读写不同的 key | sync.Map | 减少锁竞争 |
| 频繁读写同一个 key | Mutex + map | sync.Map 的额外开销反而更大 |
| 非并发场景 | 原生 map | 没有同步开销，最快 |

**sync.Map 的坑：**

1. **没有 Len() 方法**——要用 `Range` 自己数
2. **Range 不保证强一致性**——遍历过程中有并发写入，可能看到一部分新数据
3. **Key 和 Value 是 `interface{}`**——丢失了类型安全，需要自己做类型断言
4. **Range 的回调返回 false 会提前终止**——跟 `for range` 的 break 类似

### 总结一下

原生 map 并发写会直接 fatal（Go 运行时检测）。一般场景用 Mutex + map 就行，读写都要加锁。sync.Map 适合"一次写入多次读取"的缓存场景或不同 key 被不同 goroutine 操作的场景——但它没有 Len()，key/value 是 interface{}，不是通用替代品。

---

## 7. sync/atomic：原子操作

### 原子操作 vs Mutex

你的代码里有两段 atomic demo，演示了原子操作和 Mutex 的不同使用方式：

```go
// 用 atomic：一行搞定，不需要 Lock/Unlock
var sum int32 = 0
atomic.AddInt32(&sum, 1)  // 原子加 1，不会被其他 goroutine 干扰

// 用 Mutex：
mu.Lock()
num += 1
mu.Unlock()
```

| 维度 | atomic | Mutex |
|------|--------|-------|
| 操作范围 | 单个变量的简单操作 | 任意代码块 |
| 阻塞 | 不阻塞，立即完成 | 等不到锁时阻塞 |
| OS 调度器 | 不涉及 | 涉及（睡眠/唤醒） |
| 性能 | 极快（几个 CPU 周期） | 较重（几百个 CPU 周期起） |
| 复杂度 | 简单计数/标志 | 可以保护复杂逻辑 |

### atomic 提供的方法

```go
// 针对 int32/int64/uint32/uint64/uintptr
atomic.AddInt32(&sum, 1)      // 原子加
atomic.LoadInt32(&sum)        // 原子读
atomic.StoreInt32(&sum, 100)  // 原子写
atomic.SwapInt32(&sum, 200)   // 原子交换，返回旧值

// CAS：如果当前值等于 old，就替换为 new，返回是否成功
swapped := atomic.CompareAndSwapInt32(&sum, 100, 200)
```

### atomic.Value：原子地存取整个值

当要原子操作的不止是整数，而是结构体时：

```go
var v atomic.Value

v.Store(Student{Name: "zhangsan", Age: 19})  // 原子存储
stu := v.Load().(Student)                     // 原子加载（需要类型断言）

old := v.Swap(Student{Name: "lisi", Age: 20}) // 原子交换

// CAS：如果 v 当前存的是 st2，就换成 st3
swapped := v.CompareAndSwap(st2, st3)
```

### 原子操作的坑

**1. 类型必须一致**

```go
var v atomic.Value
v.Store(1)
v.Store("hello")  // ❌ panic！之前存的是 int，不能换成 string
```

`atomic.Value` 一旦存了一种类型，之后只能存同类型。它内部用接口类型做了检查。

**2. 不要用 atomic 保护复杂逻辑**

```go
// ❌ atomic 不够用——两个账户的转账需要同时原子修改
atomic.AddInt64(&accountA, -amount)
atomic.AddInt64(&accountB, +amount)
// 这两行之间如果有 goroutine 读到余额，会看到 A 扣了 B 还没加的状态

// ✅ 这种情况必须用 Mutex
mu.Lock()
accountA -= amount
accountB += amount  // 两个操作打包成一个原子逻辑
mu.Unlock()
```

**3. CAS 的 ABA 问题**

如果两次读到的是同一个值，不能保证中间没被人改过——可能被改成别的、又改回来了。这在指针场景下可能很危险。解决方案是加版本号或者用 tagged pointer（你代码里暂时没涉及，先知道有这个问题就好）。

### 总结一下

atomic 是 CPU 指令级的原子操作，不阻塞、极快，适合简单计数器和状态标志位。复杂逻辑、多变量事务性修改必须用 Mutex。atomic.Value 可以原子存取结构体，但类型必须一致，第一次 Store 就锁定了类型。

---

## 8. sync.Pool：临时对象池与复用

### 为什么需要对象池

有些对象创建频繁、用完就扔——比如网络读写时的缓冲区 `[]byte`、字符串拼接的 `bytes.Buffer`。每次都 `new` 一个，GC 扫描和回收的开销不断累积。`sync.Pool` 提供了一种**临时对象复用**机制：用完了放回去，下次需要时直接拿出来用，不用重新分配内存。

```go
type Student struct {
    Name string
    Age  int
}

pool := sync.Pool{
    New: func() interface{} {     // 池空时，调用 New 创建新对象
        return &Student{
            Name: "zhangsan",
            Age:   18,
        }
    },
}

st := pool.Get().(*Student)       // 从池里取；池空则调用 New
fmt.Printf("addr is %p\n", st)   // 比如 0xc00001c030

// 修改对象
st.Name = "hghhhhhh"
st.Age = 34

pool.Put(st)                       // 用完放回去

st1 := pool.Get().(*Student)       // 这次拿到了刚才 Put 回去的同一个对象！
println(st1.Name, st1.Age)         // hghhhhhh 34
fmt.Printf("addr1 is %p\n", st1)  // 0xc00001c030 ← 地址一样！
```

```
第一次 Get：
  检查池 → 空的 → 调用 New() → 创建 &Student{Name:"zhangsan", Age:18} → 返回 0xc00001c030

Put(st)：
  把 0xc00001c030 放回池中 → 池：[obj1]

第二次 Get：
  检查池 → 有 obj1 → 直接取出 0xc00001c030（就是刚才那个，Name 还是 "hghhhhhh"）
  不用再 New！节省了一次内存分配
```

### sync.Pool 的生命周期：GC 会清空池

**这是 sync.Pool 最关键的特性：Pool 里的对象随时可能被 GC 回收。** sync.Pool 只在两次 GC 之间保留对象——一旦发生 GC，池里的所有对象都会被清空。

```
GC 之前：
  ┌──────────────────┐
  │    sync.Pool     │
  │  obj1    obj2    │  ← 池里有 2 个复用对象
  └──────────────────┘

GC 发生 → 池被清空

GC 之后：
  ┌──────────────────┐
  │    sync.Pool     │
  │    （空的）       │  ← 全被回收了
  └──────────────────┘

下次 Get → 池空 → 调用 New() → 重新创建
```

这意味着：

- **不能假设 Put 之后 Get 一定拿到**——两次 Get 之间如果发生了 GC，拿到的是 New 创建的新对象
- **不能把 Pool 当持久缓存**——持久缓存用 `map` + `Mutex`，Pool 只能存"丢了也无所谓"的对象
- **Pool 的设计哲学：减少 GC 压力，不是消灭 GC**——高频创建临时对象时，池里的对象被反复复用，GC 扫描次数减少

### Get 的三条铁律

**① Get 拿到的对象状态不确定——必须自己初始化字段**

```go
st := pool.Get().(*Student)
// st 可能是：
//   - New 刚创建的（Name="zhangsan", Age=18）
//   - 别人 Put 回来的（Name 被改过、Age 可能是个脏值）
//
// ✅ 正确做法：Get 之后，需要什么字段就自己设
st.Name = "new_name"
st.Age = 0
// 不要依赖 Get 拿到的"默认值"
```

**② Put 之后不能再碰这个对象**

```go
st := pool.Get().(*Student)
pool.Put(st)
// ❌ Put 之后不要再访问 st！它可能已经被另一个 goroutine Get 走了
fmt.Println(st.Name)  // 数据竞争！
```

Put 意味着你放弃了对这个对象的所有权——它现在是池的财产，随时可能被别人拿走。

**③ Pool 本身是并发安全的**

多个 goroutine 可以同时 `Get` 和 `Put`，Pool 内部处理了竞态——不需要你额外加锁。

### sync.Pool 的典型应用场景

| 场景 | 例子 | 为什么适合 Pool |
|------|------|----------------|
| 网络/IO 缓冲区 | `[]byte` 频繁分配和释放 | 每次读写都分配新 buffer，GC 压力巨大 |
| 字符串拼接 | `bytes.Buffer` 的复用 | `fmt` 包内部就用 Pool 管理 Buffer |
| 序列化中间对象 | JSON/Protobuf 编解码的临时结构体 | 每个请求都要 new→用完→GC，复用后 GC 压力骤降 |
| 请求上下文包装 | 每个 HTTP 请求的 RequestContext | 高频创建、用完即弃 |

### 标准库中的 Pool 实例

`fmt` 包内部就用了 `sync.Pool` 来复用 printer：

```go
// fmt/print.go 中的实际代码（简化版）
var ppFree = sync.Pool{
    New: func() interface{} { return new(pp) },
}

func newPrinter() *pp {
    p := ppFree.Get().(*pp)
    // 拿到后重置内部状态
    return p
}

func (p *pp) free() {
    ppFree.Put(p)  // 用完放回去
}
```

每次 `fmt.Printf` 调用都会 `Get` 一个 printer、用完 `Put` 回去——这就是为什么 fmt 在高并发下性能还不错的原因之一。

### 总结一下

`sync.Pool` 用 Get/Put 实现临时对象复用，降低 GC 压力。核心坑：①GC 会清空 Pool，对象随时可能消失，不能当持久缓存；②Get 拿到的对象状态不确定，必须自己初始化关键字段；③Put 之后不能再访问该对象（所有权已转移）。适合高频创建、生命周期短、丢了也无所谓的临时对象——缓冲区、字符串 builder、编解码中间结构体。

---

## 易错点

1. **Mutex 值拷贝死锁**——`copyMutex(mu)` 传值不传指针，锁的内部状态被复制，副本的 Lock 永远等不到 → 死锁。Mutex 有 `noCopy` 标记，`go vet` 能检测到
2. **Lock/Unlock 不成对**——忘了写 Unlock，或者提前 return 跳过了 Unlock，下一个 Lock 永远等不到。**永远在 Lock 后立刻 defer Unlock**
3. **原生 map 并发写直接 fatal**——不是数据错，是程序崩。读也要加锁（读和写互斥），不只是写
4. **循环等待死锁**——两个 goroutine 按不同顺序加锁（A: mu1→mu2，B: mu2→mu1），形成等待环。**统一所有 goroutine 的加锁顺序**可以避免
5. **once.Do 嵌套死锁**——在 once.Do 的回调里再调同一个 once 的 Do，内部死锁
6. **atomic.Value 类型不一致会 panic**——第一次 Store 就固定了类型，后续 Store 不同类型会 panic
7. **代码拼写问题**——注释里 "synv.WaitGroup" 应为 "sync.WaitGroup"；函数名 `IintConfig` 应为 `InitConfig`；`dmeo2_lock` 应为 `demo2_lock`；注释 "通知可以有读个" 应为 "同时可以有多个"
8. **demo5_lock_map 里 Lock 没紧跟 defer Unlock**——虽然这个 demo 里逻辑是对的（Unlock 在 defer 里，Lock 在之后），但可读性不如标准写法 `Lock() + defer Unlock()`
9. **Pool 的对象会被 GC 回收**——不能假设 Put 之后下次 Get 一定拿到同一个对象；两次 Get 之间若发生 GC，拿到的是 New 创建的新对象。Pool 只能存"丢了也无所谓"的临时对象
10. **Pool.Get 拿到脏对象**——Get 返回的对象可能是别人 Put 回来改过的，状态不确定。Get 后必须自己初始化关键字段，不要依赖"默认值"
11. **Pool.Put 后继续用对象 → 数据竞争**——Put 后对象所有权归池，可能被其他 goroutine Get 走。继续访问 = 并发读写同一块内存

---

## 快问快答

### Q1：`sync.Mutex` 为什么不能值拷贝？

答：Mutex 内部有状态字段（当前是否锁住、等待队列等），值拷贝会把这些状态完整复制给副本。如果原版已经 Lock，副本里的状态也是"已锁定"——但副本的 Lock 会等副本的 Unlock，没人会来 Unlock 副本，于是就死锁了。永远传指针。

### Q2：读写锁的"读读共享"具体怎么理解？

答：如果当前只有读锁被持有（没有被写锁），新的 RLock 会立刻成功，不阻塞——因为读操作不修改数据，多个读者各读各的互不干扰。只有写锁的 Lock 请求出现时，新的读锁请求才会被阻塞（防止写者饿死）。

### Q3：什么时候用 sync.Map，什么时候用 Mutex + map？

答：sync.Map 适合"写入一次、读取多次"的缓存场景，或者不同 key 被不同 goroutine 操作的场景——它内部做了优化，读操作几乎不用锁。但大多数情况下 Mutex + map 更简单，且保留了类型安全。非并发场景直接用原生 map，最快。

### Q4：atomic 和 Mutex 的本质区别是什么？

答：atomic 是 CPU 硬件指令保证的，执行不被中断，不涉及 OS 调度器，不阻塞——适合单个变量的简单操作。Mutex 是 OS 调度器管理的，等不到锁就让出 CPU 去睡觉，被唤醒后再抢——可以保护任意代码块。一句话：atomic 是轻量快刀，Mutex 是万能重器。

### Q5：`wg.Add()` 为什么要写在 goroutine 外面？

答：如果把 `wg.Add(1)` 写在 goroutine 里面，可能出现这样的顺序：主 goroutine 创建了子 goroutine 但子 goroutine 还没来得及执行 `wg.Add(1)`，主 goroutine 已经执行到了 `wg.Wait()`——此时计数器是 0，Wait 直接返回，子 goroutine 还没跑。写在 goroutine 外面保证 Wait 之前 Add 已经生效。

### Q6：Go 的 Mutex 是可重入的吗？

答：不可重入。同一个 goroutine 如果 Lock 了两次，第二次 Lock 会死锁——因为锁已经被自己持有了，再次 Lock 会等自己 Unlock，而第二次 Lock 之后的代码里才有 Unlock。这是 Go 团队的设计选择：可重入锁虽然方便，但会掩盖设计问题。如果你需要可重入的话，用 `sync.Mutex` 做不到，得自己实现或者重新设计代码结构。

### Q7：为什么原生 map 并发写会 fatal 而不是静默出错？

答：这是 Go 运行时的设计哲学——**宁可直接崩掉让你看见，也别带着数据竞争上线**。并发写 map 会导致内存损坏、程序在完全无关的地方崩溃，极难排查。与其让你找 bug 找到怀疑人生，不如运行时直接告诉你是 map 并发写的问题。

### Q8：CAS 的 ABA 问题是什么意思？

答：线程 A 用 CAS 检查某个值，发现是 X（以为没变），于是换成 Y。但实际情况可能是：值先被线程 B 改成了 Z，又被线程 C 改回了 X——值确实又变成了 X，但中间经历了变化。CAS 只看"是不是同一个值"，不知道"中间有没有变化过"。在指针场景下这很危险——地址没变，但指向的对象可能已经被回收了。解决方案通常是加版本号（tagged pointer）。

### Q9：sync.Pool 和普通缓存（map+Mutex）有什么区别？

答：本质区别在于**生命周期和可靠性**。sync.Pool 是临时对象池——GC 一来就清空，对象随时可能消失，只能存"丢了也无所谓"的临时对象，优势是零维护成本、并发安全内置。map+Mutex 是持久缓存——你放进去的东西会一直在，GC 不会回收，适合存需要长期保留的数据。实际使用中它们经常组合出现：Pool 管临时对象的复用（如 `bytes.Buffer`），map+Mutex 管需要持久化的数据（如配置、缓存结果）。

---

## sync 包工具速查

| 类型 | 关键方法 | 用途 | 注意事项 |
|------|---------|------|---------|
| `sync.WaitGroup` | `Add / Done / Wait` | 等一组 goroutine 完成 | 必须传指针；Add 在 goroutine 外部 |
| `sync.Once` | `Do(func())` | 只执行一次（单例/初始化） | 不要嵌套调同一个 Once |
| `sync.Mutex` | `Lock / Unlock` | 互斥保护临界区 | 不能值拷贝；Lock 后紧接 defer Unlock |
| `sync.RWMutex` | `RLock / RUnlock / Lock / Unlock` | 读写锁，读多写少 | 读读共享，读写互斥 |
| `sync.Map` | `Store / Load / Delete / Range / LoadOrStore` | 并发安全的 map | 无 Len()；K/V 是 interface{}；适合读多写少 |
| `atomic.Add*` | `Add / Load / Store / Swap / CAS` | 原子操作简单变量 | 只适合单变量；不阻塞 |
| `atomic.Value` | `Store / Load / Swap / CAS` | 原子存取结构体 | 类型必须一致；适合读多写少 |
| `sync.Pool` | `Get / Put` + `New` 构造函数 | 临时对象复用，减少 GC | GC 会清空池；Get 后初始化字段；Put 后别碰对象 |

---

## 一句话总结

`sync` 包是 Go 并发安全的工具箱：WaitGroup 等协程、Mutex 保护临界区、RWMutex 优化读多写少、Once 保证只执行一次；map 并发要用 Mutex 保护或用 sync.Map；原子操作是轻量级的无阻塞方案适合简单场景；Pool 复用临时对象减轻 GC 压力——用对工具、记住坑点（不拷贝锁、不成对 Lock/Unlock、不加锁顺序不一致、Pool 对象会被 GC 回收），并发代码才能又安全又高效。
