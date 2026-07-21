# Golang 学习笔记

这个仓库用于系统复习 Go 基础语法、标准库常用能力和并发编程基础。当前内容以可运行的小示例为主，配套整理了一篇入门复习博客。

`theory/` 目录预留给 Go 底层原理剖析，目前还没开始写，后续可以放调度器、内存分配、GC、channel、map、interface 等专题笔记。

## 目录结构

```text
.
├── go.mod
├── docs/
│   └── 一文带你快速上手Go.md
├── quick_start/
│   ├── channels/
│   ├── context/
│   ├── defer/
│   ├── enums/
│   ├── errors/
│   ├── function/
│   ├── goroutines/
│   ├── interfaces/
│   ├── map/
│   ├── methods/
│   ├── panic/
│   ├── range/
│   ├── rune/
│   ├── select/
│   ├── slice/
│   ├── struct/
│   ├── sync/
│   ├── timer_ticker/
│   └── worker_pool/
└── theory/
```

## 当前内容

### 1. Go 基础语法

这部分放在 `quick_start/` 下，每个目录都是一个独立的小主题：

| 主题 | 目录 | 主要内容 |
| --- | --- | --- |
| 程序基本形态 | `function/` | 普通函数、返回值、递归、闭包 |
| 数据组织 | `struct/`、`methods/`、`interfaces/` | 结构体、方法、接口、`any`、类型断言 |
| 内置容器 | `map/`、`slice/`、`range/` | map 增删查、slice 追加复制、range 遍历和常见坑 |
| 字符处理 | `rune/` | 字节、rune、UTF-8 字符串遍历 |
| 枚举风格 | `enums/` | `iota`、状态流转、`String()` 方法 |

### 2. 错误处理

| 主题 | 目录 | 主要内容 |
| --- | --- | --- |
| error | `errors/` | 哨兵错误、自定义错误、`errors.Is`、`errors.As`、`errors.Join` |
| defer | `defer/` | 延迟执行、执行顺序、参数求值、资源释放 |
| panic / recover | `panic/` | panic 传播、recover 捕获、goroutine 内部恢复 |

### 3. 并发基础

| 主题 | 目录 | 主要内容 |
| --- | --- | --- |
| goroutine | `goroutines/` | `go` 关键字、`WaitGroup`、循环变量捕获、`atomic` 计数 |
| channel | `channels/` | 无缓冲、有缓冲、阻塞、关闭、`for range` |
| select | `select/` | 同时等待多个 channel、超时、取消、非阻塞分支 |
| worker pool | `worker_pool/` | 固定 worker、任务队列、结果 channel |

### 4. 并发控制

| 主题 | 目录 | 主要内容 |
| --- | --- | --- |
| Mutex | `sync/sync.Mutex/` | 保护共享变量 |
| RWMutex | `sync/sync.RWMutex/` | 读多写少缓存 |
| Once | `sync/sync.Once/` | 懒加载、只初始化一次 |
| atomic | `sync/aotmic/` | 原子计数、原子读写 |
| sync.Map | `sync/sync.Map/` | 并发安全 map 基础用法 |

### 5. 超时与定时任务

| 主题 | 目录 | 主要内容 |
| --- | --- | --- |
| context | `context/` | 取消信号、超时控制、请求级值 |
| timer / ticker | `timer_ticker/` | `Timer`、`Ticker`、`time.After`、`time.AfterFunc` |

## 推荐学习路线

1. 先看 `function/`、`struct/`、`methods/`，掌握 Go 程序、函数和类型的基本形态。
2. 再看 `slice/`、`map/`、`range/`、`rune/`，熟悉日常数据处理。
3. 接着看 `interfaces/` 和 `errors/`，理解 Go 如何用接口表达能力、如何把错误作为返回值处理。
4. 然后看 `defer/`、`panic/`，区分普通错误处理和异常退出。
5. 最后进入并发：`goroutines/`、`channels/`、`select/`、`sync/`、`context/`、`timer_ticker/`。
6. 用 `worker_pool/` 把 goroutine、channel、select、WaitGroup 串起来复习。

## 如何运行示例

仓库模块名：

```text
module learn/golang
```

Go 版本：

```text
go 1.25.0
```

运行某个示例：

```bash
go run ./quick_start/channels
go run ./quick_start/context
go run ./quick_start/worker_pool
```

运行 `sync` 子目录里的示例：

```bash
go run ./quick_start/sync/sync.Mutex
go run ./quick_start/sync/sync.RWMutex
go run ./quick_start/sync/sync.Once
go run ./quick_start/sync/sync.Map
go run ./quick_start/sync/aotmic
```

注意：`sync/aotmic` 目录名目前是 `aotmic`，内容对应的是 `sync/atomic` 示例。

## 配套文章

- [一文带你快速上手 Go](docs/一文带你快速上手Go.md)

这篇文章基于 `quick_start/` 目录整理，适合按“问题 -> 最小代码 -> 关键语法 -> 常见坑”的节奏复习 Go 基础。

## theory 规划

`theory/` 计划用于整理 Go 底层原理。当前还没有内容，后续可以按下面顺序补充：

| 主题 | 计划内容 |
| --- | --- |
| goroutine 调度 | G、M、P 的基本关系，调度发生的常见时机 |
| channel 原理 | 阻塞、唤醒、关闭、无缓冲和有缓冲差异 |
| map 原理 | 哈希表、扩容、遍历顺序不稳定 |
| interface 原理 | 动态类型、动态值、类型断言 |
| defer / panic / recover | 延迟调用、panic 传播和恢复边界 |
| GC 与内存 | 栈、堆、逃逸分析、垃圾回收基础 |
| sync 原语 | Mutex、RWMutex、Once、atomic 的使用边界 |


## 官方资料

- [A Tour of Go](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
- [The Go Blog](https://go.dev/blog/)
- [Go Packages](https://pkg.go.dev/)
- [The Go Programming Language Specification](https://go.dev/ref/spec)
