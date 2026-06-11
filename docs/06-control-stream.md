# 06-control-stream：把 if、for、switch 写成后端里的控制节奏

这个目录要抓住 Go 控制流的三个习惯：**循环只有 for，错误和存在性判断常写在 if 初始化里，switch 默认不穿透。** 这些不是语法细节，而是你以后写 handler、service、批处理时的日常节奏。

## 关键代码怎么看

`for/for.go` 里有三类循环：

```go
for i := 0; i < 5; i++ {
	fmt.Printf("cur i is %d\n", i)
}
```

这是固定次数循环。

```go
j := 0
for {
	if j == 5 {
		break
	}
	fmt.Printf("cur j is %d\n", j)
	j++
}
```

这是死循环加 `break`。能用，但如果只是模拟 while，更自然的是：

```go
for j < 5 {
	fmt.Printf("cur j is %d\n", j)
	j++
}
```

`for range` 是以后最常用的：

```go
for i, s := range sliceArr {
	fmt.Printf("slice index is %d, str is %s\n", i, s)
}

for k, v := range dic {
	fmt.Printf("key is %s, value is %d\n", k, v)
}
```

注意 map 遍历顺序不稳定，不要依赖输出顺序。需要稳定顺序时，先把 key 收集出来排序。

`if_switch/if_switch.go` 里最值得记的是 map 判断：

```go
if num, ok := dic["apple"]; ok {
	fmt.Printf("appele num %d\n", num)
}
```

这就是 Go 的经典句式：**声明一个临时变量，只在 if 块里用。** 真实项目里你会高频写：

```go
if err := svc.Do(); err != nil {
	return err
}
```

`switch`：

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

Go 的 `switch` 不需要写 `break`。每个 case 执行完自动结束。只有你明确写 `fallthrough`，才会继续往下走。

## 必须掌握的点

**1. Go 没有 while，只有 for。**  
`for condition {}` 就是 while 的写法。

**2. range 拿到的是副本。**  
`for _, v := range slice` 中的 `v` 是元素副本，改 `v` 不会改原 slice。要改原数据，用下标：

```go
for i := range users {
	users[i].Name = "new"
}
```

**3. if 初始化语句控制作用域。**  
`if v, ok := m[k]; ok {}` 让 `v` 和 `ok` 只活在 if 里，代码更干净。

**4. switch 默认自动 break。**  
这减少了 C 语言里忘写 break 的坑。

## 用一个形象例子理解

`if err := do(); err != nil {}` 像在门口做安检：检查完没问题就让流程继续，检查用的临时票据不会带进办公室。变量作用域越小，代码越不容易乱。

map 的 range 像从一个无序抽屉里抓东西，每次抓到的顺序都可能不同。你要按顺序处理，就先把标签拿出来排序，再一个个取。

## 和 Go 后端开发的关系

- 遍历请求参数：`for _, item := range req.Items`
- 批量处理数据库结果：`for rows.Next()`
- 错误处理：`if err != nil { return err }`
- 状态分支：`switch order.Status`
- map 缓存判断：`if v, ok := cache[key]; ok { ... }`

控制流写得清楚，后端业务逻辑才不会变成一坨嵌套。

## 更像工程代码的写法

- `for { if j == 5 { break } }` 可改成 `for j < 5 {}`。
- `appele` 拼错，应该是 `apple`。
- 如果 switch 处理业务状态，`default` 里最好返回错误或记录日志，不要静默吞掉未知状态。
- 遍历 map 时不要假设顺序。

## 复习时问自己

1. **Go 里 while 怎么写？**  
   `for condition { ... }`。

2. **range map 的顺序可靠吗？**  
   不可靠，需要有序时自己排序 key。

3. **`if num, ok := dic["apple"]; ok {}` 中 num 的作用域在哪里？**  
   只在 if/else 语句块内。

4. **switch 需要 break 吗？**  
   不需要，Go 默认每个 case 执行后自动结束。
