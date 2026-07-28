## 0. done channel 取消模型
首先记住：
> goroutine 停不下来，除非它自己愿意停下。`close(done)` 不是杀死 goroutine，而是发一个“该停了”的通知。

```go
package main

import (
	"fmt"
	"time"
)

func worker(done <-chan struct{}) {
	for {
		select {
		case <-done:		// 每一轮都要检查，有没有人通知我停下
			fmt.Println("worker: receive the canceling signal, exit")
			return			// 真正让 goroutine 退出的地方
		default:
			fmt.Println("worker: working")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	done := make(chan struct{})		// 创建一个停止信号通道

	go worker(done)					// 把停止信号交给 worker

	time.Sleep(2 * time.Second)

	fmt.Println("main: inform worker to stop")
	close(done)						// 广播停止信号

	time.Sleep(time.Second)
	fmt.Println("main: done~")
}
```
输出：
```text
worker: working
worker: working
worker: working
worker: working
main: inform worker to stop
worker: receive the canceling signal, exit
main: done~
```
看这条链路：
```text
main 创建 done channel
		|
		V
worker 循环工作，同时 select 监听 done
		|
		V
main close(done)
		|
		V
worker 的 <-done 被触发
		|
		V
worker return，goroutine 结束
```
需要明确一个概念，**取消的本质不是外部杀死 goroutine，而是外部发信号，goroutine 自己收到信号，goroutine 主动 return**

所以现在就能解释：`close(done)` 之后，为什么 `worker` 会退出？
> `close(done)` 会关闭 `channel`，所有监听 `<-done` 的 `goroutine` 都会立刻收到信号。`select` 检测到 `done` 后进入对应分支，`goroutine` 执行 `return`，所以它是自己退出的，不是被外部强杀的。

## 1. 函数调用链中的 ctx

一次请求通常不是一个函数完成的，而是一条调用链：
```text
handler
  -> service
    -> dao
      -> database
```
如果用户取消请求，最外层 `handler` 知道了，但真正耗时的可能是最里面的 `dao/database`。
所以，取消信号是如何从 `handler` 传到 `dao` 的？
答案就是：**ctx 一路作为参数传递下去**

```go
func Handler(ctx context.Context) error {
	return Service(ctx)
}

func Service(ctx context.Context) error {
	return QueryDB(ctx)
}

func QueryDB(ctx context.Context) error {
	return nil
}
```
`ctx` 从 `Handler` 传到 `Service`，再从 `Service` 传到 `QueryDB`

又来一个问题：为什么不能在 `QueryDB` 重新创建 `Background`？
```go
func QueryDB() error {
	ctx := context.Background()
	_ = ctx
	return nil
}
```
这样创建的是一个全新的空 `context`，和外面的请求无法关联，意思是：
```text
用户取消请求
  -> Handler 的 ctx 被取消
  -> 但 QueryDB 自己创建了 Background
  -> QueryDB 完全不知道外面取消了
```
正确的做法是：
```go
func QueryDB(ctx context.Context) error {
	// 用传进来的 ctx
	return nil
}
```

做个总结：
`ctx` 不是某个函数自己的东西，而是**这一次请求的生命周期**，所以它必须从请求入口一路往下传递：
```text
Handler(ctx)
  -> Service(ctx)
    -> QueryDB(ctx)
```
这样，最内层的数据库查询、goroutine 才能知道：
- 外面的请求是不是真的取消了？
- 是不是超时了？
- 有没有 request id / trace id？

官方文档明确强调：
> Do not store Contexts inside a struct type; instead, pass a Context explicitly to each function that needs it.

所以，现在可以回答这个问题：为什么下游操作不要自己写 `context.Background()`，而是要接收外面传进来的 `ctx`？
> `ctx` 代表一次请求的生命周期，必须沿着调用链显式传递。如果下游自己创建 `context.Background()`，就会切断和上游请求的关联。这样上游取消或超时后，最底层的 `DB/RPC/goroutine` 无法感知取消信号，可能继续占用资源。

## 2. ctx.Done() 和 ctx.Err()

现在 `ctx` 已经传到了最底层：
```go
func QueryDB(ctx context.Context) error {
	// 这里怎么知道外面取消了？
}
```
答案是两个方法：
```go
ctx.Done()
ctx.Err()
```

可以这样处理：
```go
func QueryDB(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case result := <-queryResult:
		_ = result
		return nil
	}
}
```
意思是：
谁先来就先处理谁：
- `ctx` 先取消 -> 退出，返回 `ctx.Err()`
- 查询结果先回来 -> 正常返回

源码是这样定义的：
```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}
```
- `Done()`：返回一个 `channel`，它会在当前 `context` 被取消时关闭 -> 通知你该停了
- `Err()`：如果 `Done` 还没有关闭，返回 `nil`，如果已经关闭，返回取消原因 -> 告诉你为什么要停

注意：`Done` 和 `Err` 要同时使用，因为需要告诉调用方，到底发生了什么：超时？用户取消？
所以可以总结出一个固定写法：
```go
select {
case <-ctx.Done():
	return ctx.Err()
case v := <-result:
	_ = v
	return nil
}
```
代码表明：**我在等一个操作完成，但如果外部取消或者超时，我要及时退出**

## 3. WithCancel 和 defer cancel()
假设有一个函数，里面启动了一个 worker
```go
func Run() {
	ctx, cancel := context.WithCancel(context.Background())

	go worker(ctx)

	time.Sleep(2 * time.Second)

	cancel()

	time.Sleep(1 * time.Second)
}
```
worker 是这样的：
```go
func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("worker stop:", ctx.Err())
			return
		default:
			fmt.Println("worker working")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
```

走一下完整的流程：
```text
Run 创建 ctx 和 cancel
Run 把 ctx 传给 worker
worker 一边工作，一边监听 ctx.Done()
2s 后调用 cancel()
cancel 关闭 ctx.Done()
worker 监听到了 ctx.Done()
worker return
```

这一句很关键：`go worker(ctx)`，这表示 `worker` 拿到了同一个 `ctx`，所以不是作用域控制 goroutine，而是：
- `Run` 和 `worker` 共享同一个 `ctx`
- `Run` 调用 `cancel()`
- `worker` 监听 `ctx.Done()`

和之前的 done channel 对照一下：
`context.WithCancel` 只是把流程标准化封装了：
| 手写 done | context |
| --- | --- |
| `done := make(chan struct{})` | `ctx, cancel := context.WithCancel(parent)` |
| `close(done)` | `cancel()` |
| `<-done` | `<-ctx.Done()` |
| 没有原因 | `ctx.Err()` 有原因 |

为了防止忘记 `cancel()`，一般这样写：
```go
func Run(parent context.Context) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	go worker(ctx)

	return doSomething(ctx)
}
```
为什么函数正常结束了，还要 `cancel`？
因为创建了一个子 `context`：
```text
parent
  |
  v
child ctx
```
标准库为了实现“父取消，子也取消”，会把 `child` 挂到 `parent` 下面，也就是：**parent.children 里保存 child**
如果不 `cancel`，这个 `child` 可能继续被 `parent` 引用，如果 `child` 下面还有 `worker / timer / 子 context`，也可能继续占用资源

而 `defer cancel()` 的意思是：**这个函数创建了 `child ctx`，函数结束时，我负责把它清理掉**，它有两个作用：
- 通知使用这个 `ctx` 的 `goroutine` 停。
- 释放 `context` 自己的父子引用等资源。


看一下源码：
[context.go (line 240)](/usr/lib/go-1.26/src/context/context.go:240)
```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	c := withCancel(parent)
	return c, func() { c.cancel(true, Canceled, nil) }
}
```
注意这一行：`c := withCancel(parent)`，继续往下看：
[context.go (line 273)](/usr/lib/go-1.26/src/context/context.go:273)
```go
func withCancel(parent Context) *cancelCtx {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := &cancelCtx{}
	c.propagateCancel(parent, c)
	return c
}
```
也就是说：
```go
ctx, cancel := context.WithCancel(parent)
```
表面拿到的是：
```go
ctx context.Context
```
底层真实创建的是：`*cancelCtx`，**`WithCancel` 返回的是 `Context` 接口，但接口里装的真实对象是 `*cancelCtx`**

可以这样理解链路：
```text
context.WithCancel(parent)
    |
    v
withCancel(parent)
    |
    v
c := &cancelCtx{}
    |
    v
return c 作为 Context 接口
```

`cancelCtx` 的结构在这里：[context.go (line 431)](/usr/lib/go-1.26/src/context/context.go:431)
```go
type cancelCtx struct {
	Context        						// 我的父 context

	mu       sync.Mutex
	done     atomic.Value				// 我的取消信号
	children map[canceler]struct{}		// 我的子 context
	err      atomic.Value			
	cause    error
}
```

回到 `WithCancel` 返回的第二个值：
```go
return c, func() { c.cancel(true, Canceled, nil) }
```
说明 `cancel()` 本质上是调用：
```go
c.cancel(...)
```
而 `c` 就是刚刚创建的：
```go
&cancelCtx{}
```
重新梳理一下完整链路：
```text
ctx, cancel := WithCancel(parent)

内部真实发生：

c := &cancelCtx{}
把 c 挂到 parent 下面
返回 c 给你，当作 ctx
返回一个函数 cancel

你调用 cancel()
    -> 执行 c.cancel(...)
    -> 关闭 c.done
    -> worker 的 <-ctx.Done() 被触发
```

## 4. cancel() 怎么关闭 Done()？
先看核心源码：[context.go (line 549)](/usr/lib/go-1.26/src/context/context.go:549)
```go
func (c *cancelCtx) cancel(removeFromParent bool, err, cause error) {
	if cause == nil {
		cause = err
	}

	c.mu.Lock()
	if c.err.Load() != nil {
		c.mu.Unlock()
		return
	}

	c.err.Store(err)
	c.cause = cause

	d, _ := c.done.Load().(chan struct{})
	if d == nil {
		c.done.Store(closedchan)
	} else {
		close(d)
	}

	for child := range c.children {
		child.cancel(false, err, cause)
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.Context, c)
	}
}
```
总共可以拆解为五步：
1. 如果已经取消过，直接返回
2. 设置 `err/cause`
3. 关闭 `done`
4. 遍历 `children`，把子 `context` 也取消
5. 从父 `context` 里移除自己

关键是：
```go
d, _ := c.done.Load().(chan struct{})
if d == nil {
	c.done.Store(closedchan)
} else {
	close(d)
}
```
就是上文手写的 `close(done)`，只是这里的 `done` 是 `cancelCtx` 内部保存的

所以，`cancel()` 最终效果是：**关闭这个 `ctx` 的 `Done` 信号**

那么 `worker` 就会被触发：`case <-ctx.Done(): return`

串联起整个 `WithCancel`：
```text
WithCancel(parent)
  -> 创建 *cancelCtx
  -> 把 parent 存进 cancelCtx.Context
  -> 返回 ctx 和 cancel 函数

worker(ctx)
  -> 监听 ctx.Done()

cancel()
  -> 调用 cancelCtx.cancel
  -> 设置 err
  -> close(done)
  -> worker 收到信号
  -> worker return
```
为什么还需要 `children`？因为 `context` 可以继续派生：

```go
parent, cancelParent := context.WithCancel(context.Background())
child, cancelChild := context.WithCancel(parent)
```
它们形成了：
```text
parent cancelCtx
  children:
    child cancelCtx
```

当调用了 `cancelParent()`，父会遍历 `children`：
```go
for child := range c.children {
	child.cancel(false, err, cause)
}
```
所以 `child` 也会取消


## 5. WithTimeout
假设调用了一个接口，最多等待两秒：
```go
// 一次请求
func Handler(ctx context.Context) error {
	return CallAPI(ctx)
}

// API 调用最多 2s
func CallAPI(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, 2*time.Second)
	defer cancel()

	return doRequest(ctx)
}

// 一边等 API，一边监听 ctx.Done()
func doRequest(ctx context.Context) error {
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("api success")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
```
调用链路：
```text
Handler 把请求 ctx 传给 CallAPI
CallAPI 创建一个 2 秒超时的子 ctx
doRequest 拿到这个 2 秒 ctx
doRequest 模拟 API 要 3 秒才返回
2 秒到了，ctx.Done() 关闭
doRequest 返回 ctx.Err()
```

各个函数的分工：
| 函数 | 负责什么 |
| --- | --- |
| `Handler` | 持有整次请求的 ctx |
| `CallAPI` | 给某个外部调用加超时 |
| `doRequest` | 真正执行操作，并监听 ctx.Done() |
| `defer cancel()` | CallAPI 结束时释放超时 context |
继续往下看：`WithTimeout` 不是单独的新机制，它其实是基于 `WithDeadline` 实现的。

源码在这里：[context.go (line 703)](/usr/lib/go-1.26/src/context/context.go:703)
```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}
```
翻译一下：
```text
WithTimeout(parent, 2s)
  = WithDeadline(parent, 当前时间 + 2s)
```

所以，`WithTimeout` 的本质是：**给传进来的 parent ctx 派生一个带截止时间的子 ctx**。

这就解释了为什么不能在 `CallAPI` 里自己写 `context.Background()`：
```go
func CallAPI(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, 2*time.Second)
	defer cancel()

	return doRequest(ctx)
}
```
这里必须基于 `parent` 派生，而不是自己创建新的根 `context`。因为新的 `ctx` 要同时受两个条件控制：
- `parent` 被取消，新的 `ctx` 也要取消。
- 2 秒到了，新的 `ctx` 自己也要取消。

底层真实创建的是 `timerCtx`：
[context.go (line 662)](/usr/lib/go-1.26/src/context/context.go:662)
```go
type timerCtx struct {
	cancelCtx
	timer    *time.Timer
	deadline time.Time
}
```
只看三个点：
```text
cancelCtx -> 负责取消能力
Timer     -> 负责到点触发取消
deadline  -> 记录截止时间
```
注意：`timerCtx` 里面嵌入了 `cancelCtx`，所以它也有 `done / children / err / cancel` 这些能力。它不是重新做了一套取消机制，而是在 `cancelCtx` 的基础上加了一个定时器。

关键源码在这里：
[context.go (line 652)](/usr/lib/go-1.26/src/context/context.go:652)
```go
c.timer = time.AfterFunc(dur, func() {
	c.cancel(true, DeadlineExceeded, cause)
})
```
翻译一下：
```text
启动一个 timer
时间到了
自动调用 c.cancel(...)
取消原因是 DeadlineExceeded
```

所以 `WithTimeout` 的完整链路是：
```text
WithTimeout(parent, 2s)
  -> WithDeadline(parent, now + 2s)
  -> 创建 timerCtx
  -> timerCtx 里面嵌入 cancelCtx
  -> 启动 timer
  -> 2 秒到了
  -> timer 自动调用 cancel
  -> 关闭 ctx.Done()
  -> doRequest 返回 ctx.Err()
  -> ctx.Err() = context deadline exceeded
```

还要注意：即使设置了超时，也要写 `defer cancel()`。

因为任务可能提前完成，比如：
```text
超时设置 2 秒
API 100ms 就返回了
```
如果不调用 `cancel()`，这个子 `ctx` 和它的 `timer` 可能还要等到 2 秒后才清理。`defer cancel()` 的意义是：函数结束时提前释放资源，停止 timer，并从父 `context` 的 `children` 里移除自己。

所以固定模板是：
```go
func CallAPI(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return doRequest(ctx)
}
```

现在可以总结 `WithCancel` 和 `WithTimeout` 的区别：

| 对比项 | WithCancel | WithTimeout |
| --- | --- | --- |
| 取消触发方式 | 手动调用 `cancel()` | 时间到了自动取消，也可以手动 `cancel()` |
| 底层结构 | `cancelCtx` | `timerCtx`，里面嵌入 `cancelCtx` |
| `Err()` 常见结果 | `context.Canceled` | 超时后是 `context.DeadlineExceeded` |
| 是否要 `defer cancel()` | 要 | 也要 |

## 6. WithValue

最后看 `WithValue`。它解决的问题不是取消，也不是超时，而是传递**请求级数据**。

比如一次请求里有一个 `request id`：
```go
type requestIDKey struct{}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func RequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(requestIDKey{}).(string)
	return v, ok
}
```

这样下游日志就可以拿到这次请求的 `request id`：
```go
func Log(ctx context.Context, msg string) {
	id, _ := RequestID(ctx)
	fmt.Println(id, msg)
}
```

但是要注意，`WithValue` 不是用来传普通业务参数的。它适合放：
- `request id`
- `trace id`
- 用户身份
- 认证信息

不适合放：
- `page`
- `size`
- `keyword`
- 业务配置
- 本来就应该作为函数参数传递的东西

官方文档也明确强调：
> Use context Values only for request-scoped data that transits processes and APIs, not for passing optional parameters to functions.

看源码：[context.go (line 727)](/usr/lib/go-1.26/src/context/context.go:727)
```go
func WithValue(parent Context, key, val any) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	if !reflectlite.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parent, key, val}
}
```
它最后返回的是：
```go
return &valueCtx{parent, key, val}
```
所以表面上返回的是 `Context` 接口，底层真实对象是 `*valueCtx`。

`valueCtx` 的结构在这里：[context.go (line 742)](/usr/lib/go-1.26/src/context/context.go:742)
```go
type valueCtx struct {
	Context
	key, val any
}
```
它只保存一对 `key/value`，不是一个大 map。

如果连续调用多次：
```go
ctx := context.Background()
ctx = context.WithValue(ctx, requestIDKey{}, "req-001")
ctx = context.WithValue(ctx, traceIDKey{}, "trace-abc")
```
底层更像这样：
```text
valueCtx(traceIDKey = trace-abc)
  -> valueCtx(requestIDKey = req-001)
    -> backgroundCtx
```

查找的时候，从当前层往父层找。源码在这里：[context.go (line 768)](/usr/lib/go-1.26/src/context/context.go:768)
```go
func (c *valueCtx) Value(key any) any {
	if c.key == key {
		return c.val
	}
	return value(c.Context, key)
}
```
这段代码很直白：
```text
如果当前这一层 key 命中，就返回当前 val
如果没命中，就继续去父 context 里找
```

所以 `WithValue` 的设计思想是：**每次只包一层，查询时沿着父链一层层往上找**。

这也解释了为什么不要用字符串当 key。比如两个包都写：
```go
context.WithValue(ctx, "user", user)
```
就可能互相冲突。更推荐定义自己的私有 key 类型：
```go
type userKey struct{}
```

最后总结 `WithValue`：
```text
WithValue(parent, key, val)
  -> 创建 *valueCtx
  -> valueCtx 保存一对 key/value
  -> 查 Value 时从当前层往父链找
  -> 只适合请求级数据，不适合普通业务参数
```

## 7. context 的使用场景

学完前面的内容后，`context` 的使用场景就比较清楚了。

### 场景 1：请求链路取消

用户关闭页面、客户端断开、上游请求取消，下游操作应该停止：
```go
func QueryDB(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case result := <-queryResult:
		_ = result
		return nil
	}
}
```

核心：下游必须接收上游传进来的 `ctx`，不能自己创建 `context.Background()`。

### 场景 2：控制某个操作的最长耗时

调用外部接口、RPC、数据库时，经常要设置超时：
```go
func CallAPI(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return doRequest(ctx)
}
```

核心：谁要限制这一步最多执行多久，谁就在这一层派生 `WithTimeout`。

### 场景 3：防止 goroutine 泄漏

启动 goroutine 时，把 `ctx` 传进去，让它能主动退出：
```go
func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// do work
		}
	}
}
```

核心：`context` 不会强杀 goroutine，它只是发取消信号，goroutine 要自己监听并退出。

### 场景 4：传递请求级数据

比如 `request id / trace id / 用户身份`：
```go
type requestIDKey struct{}
ctx = context.WithValue(ctx, requestIDKey{}, "req-001")
```

核心：只传请求级元数据，不传普通业务参数。

## 8. context 有哪几种数据结构实现

面试时经常问：`context` 底层有哪些实现？不用背所有边角类型，先掌握这几个核心结构。

| 数据结构 | 来自哪个 API | 作用 |
| --- | --- | --- |
| `emptyCtx` | `Background` / `TODO` 的基础 | 空 context，不会取消、没有 deadline、没有 value |
| `backgroundCtx` | `context.Background()` | 正常根 context |
| `todoCtx` | `context.TODO()` | 临时占位 context |
| `cancelCtx` | `WithCancel` / `WithCancelCause` | 负责取消，内部有 `done`、`children`、`err/cause` |
| `timerCtx` | `WithTimeout` / `WithDeadline` | 负责超时，嵌入 `cancelCtx`，额外有 `timer` 和 `deadline` |
| `valueCtx` | `WithValue` | 保存一对 `key/value`，查找时沿父链向上找 |
| `withoutCancelCtx` | `WithoutCancel` | 保留 value，但切断取消和 deadline |
| `afterFuncCtx` | `AfterFunc` | context 取消后触发回调 |

如果只面试核心，可以重点说前三类派生结构：
```text
cancelCtx -> 取消
timerCtx  -> 超时
valueCtx  -> 请求级值
```

其中最重要的是 `cancelCtx`：
```text
cancelCtx
  -> Context: 父 context
  -> done: 取消信号
  -> children: 子 context
  -> err/cause: 取消原因
```

`timerCtx` 是在 `cancelCtx` 上加 timer：
```text
timerCtx
  -> cancelCtx
  -> timer
  -> deadline
```

`valueCtx` 是链式包装：
```text
valueCtx(key2, val2)
  -> valueCtx(key1, val1)
    -> backgroundCtx
```

## 9. 最后总结

`context` 的本质是：**在一条调用链中传递取消信号、超时时间和请求级数据**。

先理解这条主线：
```text
handler(ctx)
  -> service(ctx)
    -> dao(ctx)
      -> db/rpc(ctx)
```
`ctx` 代表这次请求的生命周期。它必须一路传下去，不能在下游重新创建 `context.Background()`，否则取消链路就断了。

取消的核心是：
```text
cancel()
  -> 关闭 ctx.Done()
  -> 监听 ctx.Done() 的 goroutine 收到信号
  -> goroutine 自己 return
```

源码上，`WithCancel` 返回的 `ctx` 底层是 `*cancelCtx`。它里面有父 `Context`、`done`、`children`、`err/cause`。父 context 取消时，会遍历 `children`，把子 context 一起取消。

超时的核心是：
```text
WithTimeout
  -> WithDeadline
  -> timerCtx
  -> timer 到期自动 cancel
  -> Err() 返回 context deadline exceeded
```

传值的核心是：
```text
WithValue
  -> valueCtx
  -> 每次只保存一对 key/value
  -> Value 查询时沿父链查找
```

最后记住实战规则：
1. `ctx` 作为函数第一个参数一路传递。
2. 不要把 `context` 存进结构体。
3. 不要传 `nil`，不确定时用 `context.TODO()`。
4. 派生了 `WithCancel / WithTimeout / WithDeadline`，就要 `defer cancel()`。
5. `WithValue` 只放请求级元数据，不放普通业务参数。

面试可以这样一句话收尾：
> `context` 用来控制请求生命周期。`cancelCtx` 负责取消传播，`timerCtx` 负责超时控制，`valueCtx` 负责请求级值传递。实战中要把 `ctx` 显式往下传，监听 `ctx.Done()` 及时退出，派生子 context 后及时 `defer cancel()` 释放资源。
