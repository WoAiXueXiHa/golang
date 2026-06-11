# 03-const-val：常量重点看 iota 和“编译期确定”

这个目录不需要把常量讲成大部头。你只要抓住两点：**const 是编译期确定的值**，**iota 是 Go 用来写枚举和位标志的计数器**。真实后端里，状态码、错误码、权限位、固定协议值都会用到它。

## 关键代码怎么看

包级常量：

```go
const (
	Unkonw  = 0
	Success = 1
	Fail    = 2
)
```

这里想表达的是“状态枚举”。Go 没有 enum，所以常用 `const (...)` 表示一组相关常量。代码能跑，但 `Unkonw` 拼错了，应该是 `Unknown`。如果这是工程代码，拼写错误会很影响专业感。

这段也很适合改成 iota：

```go
const (
	Unknown = iota
	Success
	Fail
)
```

`testCal()` 里这个点值得记：

```go
const (
	a = "123"
	b = 10
	c = unsafe.Sizeof(a)
)
```

`unsafe.Sizeof(a)` 能用于常量，是因为它能在编译期确定。这里 `a` 是 string，64 位机器上 string header 通常是 16 字节：一个数据指针加一个长度字段。注意它测的是 string 这个“描述头”的大小，不是 `"123"` 内容的长度。内容长度要用 `len(a)`。

`testIota()` 是本目录重点：

```go
const (
	a = iota
	b
	c
)
```

输出是 `0 1 2`。记住：**iota 在每个 const 块里从 0 开始，按行递增。空行不算，常量声明行才算。**

复杂一点的这组：

```go
const (
	val1, val2 = iota + 1, iota + 2
	val3, val4
	val5, val6 = iota + 10, iota * 10
	val7, val8
)
```

关键规则是：没写表达式的行会继承上一行表达式，但 iota 用当前行的值。所以结果是：

```text
val1=1 val2=2
val3=2 val4=3
val5=12 val6=20
val7=13 val8=30
```

## 必须掌握的点

**1. const 不是变量。**  
常量没有地址，不能取 `&`，也不能运行时修改。它更像编译器在代码里替你放好的固定值。

**2. iota 只在当前 const 块内有效。**  
新开一个 `const (...)`，iota 又从 0 开始。

**3. iota 按行递增，不按变量个数递增。**  
同一行多个常量共享同一个 iota。

**4. 后端里不要滥用魔法数字。**  
`if status == 2` 不如 `if status == StatusFailed`。常量让代码表达业务含义。

## 用一个形象例子理解

iota 像排队取号机。每写一行常量声明，就叫一次号：第一行 0，第二行 1，第三行 2。你不写新表达式时，就沿用上一张表格模板，只是把号码换成当前号码。

比如 `iota + 10` 这张模板，在 iota 是 2 时得到 12，在 iota 是 3 时得到 13。

## 和 Go 后端开发的关系

常量最常见的后端场景：

- 业务状态：`OrderPending`、`OrderPaid`、`OrderFailed`
- 错误码：`ErrUserNotFound`、`ErrTokenExpired`
- 权限位：`PermRead = 1 << iota`
- 固定配置默认值：`DefaultPageSize = 20`

但要区分：**编译期不会变的用 const，运行时可能变的用配置**。数据库地址、端口、超时时间通常不应该写死成常量。

## 更像工程代码的写法

- 修正 `Unkonw` 为 `Unknown`。
- 枚举型常量优先用 iota，维护更轻。
- 常量名字要表达业务，不要只叫 `a`、`b`，除非是教学代码。
- 位权限可以练习 `1 << iota`，这是 Go 后端里很实用的写法。

## 复习时问自己

1. **iota 在哪里重置？**  
   每个新的 `const (...)` 块中重置为 0。

2. **`unsafe.Sizeof("abc")` 为什么不是 3？**  
   因为它测的是 string header 的大小，不是内容长度。

3. **状态码为什么适合用 const？**  
   因为它们是编译期固定的业务含义，不应该运行时随便变化。

4. **什么时候不用 const，而用配置？**  
   当值需要不同环境调整，比如端口、数据库 DSN、请求超时。
