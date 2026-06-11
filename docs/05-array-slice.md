# 05-array-slice：slice 和 map 是 Go 后端的日常武器

这一章只抓本质：**数组是固定容器，slice 是窗口，map 是哈希表。**  
后端开发里，数组了解即可；真正天天用的是 slice 和 map。

## 这份代码最该看什么

```go
var strArr = [10]string{"aa", "bb", "cc"}
arr := [3]int{1, 2, 3}
```

数组长度是类型的一部分。`[10]string` 和 `[3]string` 是两个类型。数组是值类型，赋值、传参都会复制整块数据。真实后端业务里很少直接用数组，除非是固定长度数据，比如 `[16]byte`、`[32]byte`。

```go
var sliceArr = make([]string, 0)
sliceArr = strArr[1:3]
```

这里第一行 `make` 是多余的，因为下一行立刻把 `sliceArr` 指向了 `strArr[1:3]`。更直接：

```go
sliceArr := strArr[1:3]
```

`strArr[1:3]` 不复制 `"bb"` 和 `"cc"`，它只是创建一个 slice，指向数组 `strArr` 的一段。改 `sliceArr[0]`，就是改 `strArr[1]`。

```go
var dic = map[string]int{
	"apple":  1,
	"banana": 2,
}
```

map 是哈希表。读取不存在的 key 返回零值，不会报错。必须用 `ok` 判断是否真的存在：

```go
v, ok := dic["orange"]
```

## slice 的本质

slice 不是数组。它是一个小结构：

```go
type sliceHeader struct {
	ptr *T
	len int
	cap int
}
```

可以把 slice 想成“窗口”：

- `ptr`：窗口从哪里开始看
- `len`：当前能看到几个元素
- `cap`：不换底层数组的情况下，最多能往后扩多远

例子：

```go
arr := [5]int{10, 20, 30, 40, 50}
s := arr[1:3] // [20, 30]
```

此时 `s[0]` 对应的是 `arr[1]`。所以：

```go
s[0] = 99
fmt.Println(arr) // [10 99 30 40 50]
```

面试会追问：**切片传参是值传递还是引用传递？**  
答案：还是值传递。复制的是 slice header，但 header 里的 ptr 指向同一个底层数组，所以修改元素会影响外部。

```go
func change(s []int) {
	s[0] = 100
}
```

这个会影响外部元素。  
但如果函数里重新让 `s` 指向新 slice，不会改变调用方的 slice 变量本身。

## append 的核心坑

`append` 可能复用原数组，也可能扩容换新数组。

```go
s = append(s, x)
```

必须接住返回值，因为 append 返回的 slice 可能已经指向新的底层数组。

最容易被拷打的例子：

```go
a := []int{1, 2, 3}
b := a[:2]
b = append(b, 99)
fmt.Println(a) // 可能是 [1 2 99]
```

为什么？因为 `b` 还有容量，append 直接写回原底层数组。

再看扩容：

```go
a := []int{1, 2, 3}
b := a[:2:2]      // 第三个索引限制 cap
b = append(b, 99) // cap 不够，扩容
fmt.Println(a)    // [1 2 3]
fmt.Println(b)    // [1 2 99]
```

三索引切片 `a[low:high:max]` 很少日常写，但它能控制 cap，是理解 append 副作用的好工具。

## nil slice 和空 slice

这点后端接口很常见。

```go
var a []int        // nil slice
b := []int{}       // empty slice
c := make([]int, 0)
```

它们 `len` 都是 0，但：

- `a == nil` 为 true
- `b == nil` 为 false
- JSON 里 nil slice 常见输出 `null`
- 空 slice 常见输出 `[]`

接口返回列表时，通常希望返回 `[]`，而不是 `null`。

```go
users := make([]User, 0)
```

这是很常见的后端习惯。

## map 的工程重点

map 读不存在 key：

```go
v := dic["orange"] // v 是 0
```

问题是：你不知道它是“不存在”，还是“存在但值就是 0”。所以工程里要写：

```go
v, ok := dic["orange"]
if !ok {
	// key 不存在
}
```

nil map：

```go
var m map[string]int
m["a"] = 1 // panic
```

必须先初始化：

```go
m := make(map[string]int)
```

map 并发：

普通 map 不是并发安全的。并发读写可能直接崩：

```text
fatal error: concurrent map read and map write
```

工程选择：

- 简单共享 map：`sync.RWMutex`
- 读多写少、key 独立：`sync.Map`
- 单 goroutine 管理状态：channel 串行化

## 后端实践怎么用

**列表响应：**

```go
users := make([]User, 0, len(ids))
for _, id := range ids {
	users = append(users, loadUser(id))
}
```

提前给 cap，减少扩容。

**去重：**

```go
seen := make(map[int64]struct{})
for _, id := range ids {
	seen[id] = struct{}{}
}
```

`struct{}{}` 不占额外空间，比 `map[int64]bool` 更像集合。

**分组：**

```go
groups := make(map[string][]Order)
for _, order := range orders {
	groups[order.Status] = append(groups[order.Status], order)
}
```

这就是后端里非常高频的数据整理方式。

## 本目录代码可改进点

- `make([]string, 0)` 后立刻覆盖，删掉。
- `fmt.Printf("strArr没有+v: %v\n", arr)` 文案和变量不一致。
- 建议补一个 `append` 示例，否则 slice 最关键的地方没有练到。
- 建议补一个 `v, ok := dic[key]` 示例，这是 map 工程用法。

## 面试拷打

1. **数组和 slice 的区别？**  
   数组长度固定且长度属于类型；slice 是动态视图，底层指向数组。

2. **slice 传参会不会复制？**  
   会复制 slice header，但不会复制底层数组。

3. **为什么 append 后必须接返回值？**  
   append 可能扩容，返回的新 slice 可能指向新数组。

4. **nil slice 和空 slice 区别？**  
   len 都是 0，但 nil 判断和 JSON 输出可能不同。

5. **map 并发安全吗？**  
   不安全。并发读写需要锁、`sync.Map` 或其他同步设计。

6. **map 取不到 key 会怎样？**  
   返回 value 类型零值，所以必须用 `v, ok` 区分。
