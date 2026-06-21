# Go 控制流

## 这一章要记住什么

这一章主要讲四个点：

- Go 的循环只有 `for`，没有 `while`。
- `if` 可以带初始化语句，变量作用域只在 `if` 内。
- `switch` 默认每个 `case` 自动结束，不需要手写 `break`。
- `for range` 很方便，但循环变量、切片长度快照、值拷贝、map 顺序、字符串遍历这些细节踩坑率很高。

---

## 1. if 和 if 初始化语句

代码里：

```go
if localStr == "case3" {
    fmt.Println("into true logic")
} else {
    fmt.Println("into false logic")
}
```

这是最普通的条件判断。

更常见的是带初始化语句：

```go
if num, ok := dic["apple"]; ok {
    fmt.Printf("apple num %d\n", num)
}
```

`num, ok := dic["apple"]` 先执行。

`ok` 表示这个 key 是否存在。

```text
dic["apple"]
    |
    v
返回两个值：num, ok
    |
    v
ok == true  -> key 存在
ok == false -> key 不存在
```

## 总结一下

`if 初始化; 条件 {}` 很适合处理 map 查询、函数返回值检查这类场景。

初始化出来的变量只在这个 `if` 里有效。

---

## 2. switch

代码里：

```go
switch localStr {
case "case1":
    fmt.Println("case1")
case "case2":
    fmt.Println("case2")
case "case3":
    fmt.Println("case3")
default:
    fmt.Println("default")
}
```

Go 的 `switch` 匹配到一个 `case` 后，默认不会继续往下执行。

```text
localStr = "case3"
        |
        v
case1? no
case2? no
case3? yes
        |
        v
执行 case3，然后结束 switch
```

## 总结一下

Go 的 `switch` 默认自带 `break` 效果，不需要每个分支都写 `break`。

---

## 3. for 的几种写法

代码注释里总结了三种：

```go
for i := 0; i < 5; i++ {}

for condition {}

for {}
```

可以对应理解成：

```text
传统 for：
初始化 -> 条件判断 -> 循环体 -> 后置语句

类似 while：
条件判断 -> 循环体

死循环：
一直执行，通常配合 break
```

## 总结一下

Go 只有 `for` 一个循环关键字，但它可以覆盖 C/C++ 里的 `for`、`while` 和死循环。

---

## 4. for range 的循环变量与闭包陷阱

`for range` 的循环变量有两个经典的"地址/闭包"问题，根源相同，但表现形式不同。

### 4.1 取地址：大家都指向同一个变量？

```go
arr := [2]int{1, 2}
res := []*int{}
for _, v := range arr {
    res = append(res, &v)
}
for _, p := range res {
    fmt.Println(*p)
}
```

这段代码看似要把 `1` 和 `2` 的指针都存下来，实际行为取决于 Go 版本：

**1.22 之前**：整个循环从头到尾只有一个 `v` 变量，一块内存地址。每轮迭代把新值塞进这块内存，`res` 里存的全是同一个 `&v`。循环结束后 `v` 里是最后一个值 `2`，所以解引用打印出来全是 `2`。

**1.22 之后**：每轮迭代创建一个全新的 `v` 变量，各自有独立的内存地址。`res` 里存的是不同地址，解引用打印出 `1 2`。

```text
1.22 之前：
v 的地址: 0x100 ──┐
                 ├── 第0轮：v=1 → res = [&v]
                 ├── 第1轮：v=2 → res = [&v, &v]
                 └── 全部指向 0x100，最终值都是 2

1.22 之后：
v0 地址: 0x100 ── 第0轮：v0=1 → res = [&v0]
v1 地址: 0x108 ── 第1轮：v1=2 → res = [&v0, &v1]
                 各自独立，打印 1 2
```

### 4.2 闭包捕获：所有函数都记住了同一个 i？

```go
var funcs []func()
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() {
        fmt.Println(i)
    })
}
for _, f := range funcs {
    f()
}
```

匿名函数捕获了外层的 `i`。

**1.22 之前**：`i` 只有一块内存，三个闭包捕获的是同一个 `i` 的地址。循环结束后 `i` 被加到了 `3`（第四轮 `i=3`，条件不成立退出），执行所有函数都打印 `3`。

**1.22 之后**：每轮迭代的 `i` 地址相互隔离，三个闭包各抓各的，打印 `0 1 2`。

### 4.3 兼容写法：不用猜版本也能写对

**方案一：同名变量遮蔽（最常用）**

```go
for _, v := range arr {
    v := v                // 在当前作用域内强行分配一个全新的局部变量
    res = append(res, &v) // 这个 &v 指向的是新变量，互不干扰
}

for i := 0; i < 3; i++ {
    i := i                // 同理，闭包捕获的是这个新 i
    funcs = append(funcs, func() {
        fmt.Println(i)
    })
}
```

```text
每轮循环：
外层 v 拷贝 ──→ 内层 v（新地址，独立内存）
                 ↑
            闭包/指针绑定的是这个新的内层 v
```

**方案二：函数参数传值**

```go
for i := 0; i < 3; i++ {
    funcs = append(funcs, func(val int) {
        fmt.Println(val)
    }(i)) // i 作为参数传入，发生值拷贝
}
```

**方案三：用索引访问原始数据（适用于切片/数组）**

```go
for i := range arr {
    res = append(res, &arr[i]) // 直接取原数组元素的地址
}
```

## 总结一下

> 1.22 之前，循环变量是一块"公用内存"，每轮覆盖旧值；1.22 之后，每轮分配独立变量，生而隔离。

写 `v := v` 或通过函数参数传递，是兼容所有版本的稳妥写法。核心原则：**确认闭包或指针绑定的变量，到底是每轮共享的那一个，还是每轮独立的那一个。**

---

## 5. for range 的几个细节

### 5.1 切片遍历：循环次数在开始前就已确定

```go
v := []int{1, 2, 3}
for i := range v {
    v = append(v, i)
}
// 这会死循环吗？—— 不会。
```

`for range` 在循环开始前，会把切片的长度"拍照"记录下来。循环体内对原切片做任何 `append`，都是在修改 `v` 这个变量，但控制循环次数的那个长度值早就被编译期锁死在了临时变量里。

```text
编译期展开的等价伪代码：

rangeSlice := v
length := len(rangeSlice) // ← 循环开始前就固定了

for i := 0; i < length; i++ {
    v = append(v, i)      // 修改的是 v，length 纹丝不动
}
```

所以上面那段代码只会执行 3 次循环，不会死循环。

### 5.2 切片遍历：value 是元素副本

```go
slice := []int{1, 2, 3}
for _, v := range slice {
    v *= 10 // 改的是副本，原切片毫发无伤
}
fmt.Println(slice) // [1 2 3]
```

`for _, v := range slice` 执行时，底层发生了一次值拷贝：`v = slice[index]`。`v` 拥有自己独立的栈内存，你怎么改它，都影响不了原切片的数据。

```text
slice[0] = 1 ──拷贝──→ v = 1（独立内存）
                         │
                      v *= 10 → v = 10
                         │
                      slice[0] 还是 1
```

要修改原切片，直接用索引：

```go
for i := range slice {
    slice[i] *= 10 // 这才是改原切片
}
```

### 5.3 map 遍历：顺序不固定

```go
dic := map[string]int{"a": 1, "b": 2, "c": 3}
for k, v := range dic {
    fmt.Printf("%s:%d ", k, v)
}
// 每次运行，输出顺序可能都不一样
```

这不是 bug，是 Go 故意设计的。`for range map` 底层调用 `mapiterinit` 初始化迭代器时，会用一个随机数种子来选择遍历的起始桶和桶内偏移量。

```text
每次 range map：
  → mapiterinit()
  → 随机选起始桶
  → 从那个桶开始遍历
  → 顺序每次都不同
```

为什么要随机？**防止开发者依赖某种特定的遍历顺序，导致代码在底层哈希表扩容时出现诡异 bug。** 实际上数据在哈希桶里的物理存放是有序的，只是每次遍历起点被随机"洗了一次牌"。

需要稳定顺序时，先取 key 排序再遍历：

```go
keys := make([]string, 0, len(dic))
for k := range dic {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Printf("%s:%d ", k, dic[k])
}
```

### 5.4 字符串遍历：rune 视角 vs 字节视角

```go
str := "hello 世界"

// for range 遍历：rune 视角 — 按 UTF-8 字符解码
for i, r := range str {
    fmt.Printf("index: %d, rune: %c\n", i, r)
}
// 输出：
// index: 0, rune: h
// index: 1, rune: e
// ...
// index: 6, rune: 世
// index: 9, rune: 界
//              ↑ 注意：index 从 6 直接跳到了 9！

// 传统 for 循环：字节视角 — 逐字节裸读
for i := 0; i < len(str); i++ {
    fmt.Printf("index: %d, byte: %x\n", i, str[i])
}
// 中文的 3 个字节被拆开单独打印，根本拼不出完整的"世界"
```

为什么 `for range` 的 index 会跳？Go 默认使用 **UTF-8** 编码，英文占 1 字节，中文占 3 字节。`for range` 遍历字符串时，底层会调用 `decoderune` 函数，遇到多字节的中文编码时，自动向后"嗅探"并合并 3 个字节，解码成一个完整的 `rune`。下一次迭代的起始索引自然就跨越了 3 个字节。

```text
"hello 世界" 的内存布局（字节序列）：
 h  e  l  l  o  _   [世：3字节]   [界：3字节]
 0  1  2  3  4  5   6   7   8     9  10  11

for range（rune 视角）：
  index: 0,1,2,3,4,5   →   6(世)   →   9(界)
  每次读一个完整字符      自动合并3字节  自动合并3字节

传统 for（byte 视角）：
  index: 0,1,2,3,4,5,6,7,8,9,10,11
  逐字节裸读，中文被肢解
```

## 总结一下

- **循环次数**：`range` 切片时，长度在开始前"拍照"锁定，循环内 `append` 不会导致死循环。
- **值拷贝**：`range` 拿到的 `v` 是元素副本，改它不改原数据，要改原数据请用索引。
- **map 随机**：遍历顺序是故意随机化的，需要稳定顺序就先把 key 排好序。
- **字符串**：`for range` 是"智能嗅探"的 rune 视角，传统 `for` 是"逐字节裸读"的 byte 视角。处理中文等非 ASCII 字符时优先用 `for range`。

---

## 易错点

1. `if num, ok := dic[key]; ok {}` 里的 `num` 和 `ok` 只在 if 内有效。
2. Go 的 `switch` 默认不会贯穿到下一个 case。
3. `for range` 里取 `&v`，1.22 之前拿到的全是同一个地址。兼容写法：`v := v`。
4. 循环里创建闭包时，1.22 之前所有闭包捕获同一个变量地址。兼容写法：`i := i` 或函数参数传值。
5. `for range` 的 value 是元素副本，修改它不影响原切片。要改原数据用索引。
6. `for range` 切片时长度在开始前就已确定，循环内 `append` 不会影响循环次数。
7. map 遍历顺序不固定，不能依赖。需要顺序时先对 key 排序。
8. `for range` 遍历字符串拿到的是 rune，index 遇到多字节字符会跳跃；传统 `for` 拿的是字节。

---

## 快问快答

### Q1：Go 有 `while` 吗？

答：

没有。Go 只有 `for`，可以用 `for condition {}` 表达 while 的效果。

### Q2：map 查询为什么常写 `num, ok := dic[key]`？

答：

因为 value 的零值可能也是合法值，必须用 `ok` 判断 key 是否真的存在。

### Q3：`range` 遍历切片时，修改 value 会影响原切片吗？

答：

不会。value 是元素副本。要改原切片，应该用索引修改 `slice[i]`。

### Q4：map 遍历顺序稳定吗？

答：

不稳定。底层在初始化迭代器时用了随机种子来选择起始桶，每次运行顺序都可能不同。需要稳定顺序时，先取 key 排序再遍历。

### Q5：循环里往切片 `append`，会导致死循环吗？

答：

不会。`for range` 在循环开始前就确定了切片长度，循环体内的 `append` 影响的是原变量，但循环次数已经被"拍照"锁死了。

### Q6：`for range` 字符串和 `for i := 0; i < len(str); i++` 有什么区别？

答：

`for range` 按 UTF-8 字符（rune）解码遍历，遇到中文等多字节字符会自动合并，index 会跳跃。传统 `for` 是逐字节裸读，多字节字符会被拆散。处理非 ASCII 文本时优先用 `for range`。

---

## 一句话总结

Go 控制流语法不多，但 `if 初始化`、`switch 默认 break`、`for range 的循环变量/长度快照/值拷贝/map 随机/字符串 rune` 这些细节，面试和实际写代码都很容易踩坑。记住一个核心原则：**`for range` 本质是编译器在循环开始时"拍照"快照的语法糖——你操作的是快照还是原始数据，心里要有数。**
