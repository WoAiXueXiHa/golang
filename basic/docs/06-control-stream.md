# Go 控制流

## 这一章要记住什么

这一章主要讲四个点：

- Go 的循环只有 `for`，没有 `while`。
- `if` 可以带初始化语句，变量作用域只在 `if` 内。
- `switch` 默认每个 `case` 自动结束，不需要手写 `break`。
- `for range` 很方便，但要注意循环变量、切片长度、map 顺序这些细节。

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

## 4. for range 的循环变量

当前代码里重点放在闭包捕获循环变量：

```go
var funcs []func()

for i := 0; i < 3; i++ {
	funcs = append(funcs, func() {
		fmt.Println(i)
	})
}
```

匿名函数捕获了外层的 `i`。

在老版本 Go 中，循环变量可能复用同一块内存，所以最后执行这些函数时，容易拿到同一个最终值。

```text
循环阶段：
i = 0 -> 保存函数，函数记住 i
i = 1 -> 保存函数，函数还是记住同一个 i
i = 2 -> 保存函数，函数还是记住同一个 i

执行阶段：
循环结束后 i 已经变了
多个函数读到同一个 i
```

常见修复方式：

```go
for i := 0; i < 3; i++ {
	i := i
	funcs = append(funcs, func() {
		fmt.Println(i)
	})
}
```

```text
每轮循环：
外层 i -> 拷贝出一个新的内层 i
匿名函数捕获的是这一轮自己的 i
```

## 总结一下

循环里创建闭包时，要确认闭包捕获的是不是每轮独立的变量。

写 `i := i` 或者通过函数参数传进去，是很常见的保护写法。

---

## 5. for range 的几个细节

代码注释里还提到了几个常见点：

- range 遍历切片时，循环开始前会先确定长度。
- range 得到的值变量是元素副本，改它不会改原切片。
- range 遍历 map 的顺序是随机的。
- range 遍历字符串时，拿到的是 rune，不是单纯的字节。

```text
range slice:
每轮拿到 index 和 value
value 是元素副本

如果要改原切片：
用 slice[index] 修改
```

## 总结一下

`range` 很好用，但不要误以为它拿到的一定是原始元素本身。

需要修改原数据时，优先用索引。

---

## 易错点

1. `if num, ok := dic[key]; ok {}` 里的 `num` 和 `ok` 只在 if 内有效。
2. Go 的 `switch` 默认不会贯穿到下一个 case。
3. `for range` 的 value 通常是副本。
4. map 遍历顺序不固定。
5. 循环里创建闭包时，要注意捕获变量的问题。

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

不会。value 通常是元素副本。要改原切片，应该用索引修改 `slice[i]`。

### Q4：map 遍历顺序稳定吗？

答：

不稳定。不能依赖 map 的遍历顺序，需要稳定顺序时，先取 key 排序再遍历。

---

## 一句话总结

Go 控制流语法不多，但 `if 初始化`、`switch 默认 break`、`for range 副本和闭包` 这些细节很容易考。

