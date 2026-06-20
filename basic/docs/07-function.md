# Go 函数

## 这一章要记住什么

这一章主要讲三个点：

- Go 函数支持多返回值，常用来返回结果和错误。
- Go 只有值传递，传指针也是复制一份地址值。
- 结构体值传递是浅拷贝，字段里的指针仍然指向同一份数据。

---

## 1. 多返回值和 error

代码注释里写过这个函数：

```go
func RegisterUser(username string, password string) (string, error) {
	if len(username) < 3 {
		return "", errors.New("用户名长度不能小于3位")
	}
	return username + ", 欢迎加入!", nil
}
```

Go 很常见的函数返回形式是：

```text
结果值, error
```

调用方通常这样处理：

```go
msg, err := RegisterUser("Golang", "123456")
if err != nil {
	return
}
fmt.Println(msg)
```

```text
调用函数
   |
   v
返回 result, err
   |
   +-- err != nil：处理错误，提前返回
   |
   +-- err == nil：继续使用 result
```

## 总结一下

Go 不把异常作为主要错误处理方式，而是把错误当成普通返回值交给调用方显式处理。

---

## 2. Go 只有值传递

当前代码里真正运行的是这个例子：

```go
type UserData struct {
	Age      int
	PtrScore *int
}

func PassByValue(u UserData) {
	u.Age = 30
	*u.PtrScore = 99
}
```

调用时：

```go
initialScore := 100
user := UserData{
	Age:      18,
	PtrScore: &initialScore,
}

PassByValue(user)
```

`PassByValue(user)` 会复制一份 `UserData`。

```text
main 里的 user
+----------------+
| Age: 18        |
| PtrScore: ---- | ----+
+----------------+     |
                       v
                  initialScore = 100

传入函数后，复制一份 u
+----------------+
| Age: 18        |
| PtrScore: ---- | ----+
+----------------+     |
                       v
                  initialScore = 100
```

函数里：

```go
u.Age = 30
```

改的是副本里的 `Age`。

```go
*u.PtrScore = 99
```

通过指针字段改的是原来的 `initialScore`。

```text
u.Age = 30
只改副本：
main.user.Age 还是 18

*u.PtrScore = 99
顺着同一个地址修改：
initialScore 变成 99
```

## 总结一下

Go 只有值传递。

但是如果被复制的值里面有指针字段，那么指针地址也会被复制，两个结构体副本仍然可能指向同一份底层数据。

---

## 3. 值传递和浅拷贝

这个例子最关键的地方是：结构体被复制了，但结构体里的指针指向的数据没有被复制。

```text
结构体浅拷贝：

原结构体             副本
+--------+          +--------+
| Age 18 |          | Age 18 |
| ptr ---|---+  +---| ptr    |
+--------+   |  |   +--------+
             v  v
          同一个 score
```

所以输出结果会表现为：

```text
Age 没变
Score 被改了
```

## 总结一下

浅拷贝只复制结构体字段本身。

如果字段是指针、slice、map 这种内部带引用关系的类型，就要小心共享底层数据。

---

## 4. 函数作为值

代码注释里也写过函数类型：

```go
type LogFilter func(msg string) bool
```

这说明函数可以像普通值一样传递。

```go
func ProcessLogs(logs []string, f LogFilter) {
	for _, log := range logs {
		if f(log) {
			fmt.Println(log)
		}
	}
}
```

```text
ProcessLogs
    |
    v
接收一组日志 + 一个过滤函数
    |
    v
用传进来的函数决定哪些日志要保留
```

## 总结一下

函数在 Go 里是一等公民，可以赋值给变量，也可以作为参数和返回值。

---

## 5. 闭包

代码注释里有这个例子：

```go
factor := 2
multiplier := func(val int) int {
	factor++
	return val * factor
}
```

匿名函数引用了外层变量 `factor`，这就是闭包。

```text
外层变量 factor
      ^
      |
匿名函数内部继续使用它
```

只要这个函数还在，`factor` 就会跟着这个函数一起被保存下来。

## 总结一下

闭包就是函数捕获了外层变量。

它很方便，但也要注意变量生命周期和循环变量捕获的问题。

---

## 易错点

1. Go 只有值传递，没有引用传递。
2. 传指针也是值传递，只是复制的是地址值。
3. 结构体拷贝是浅拷贝，指针字段可能共享底层数据。
4. 多返回值里 `error` 通常放最后。
5. 闭包捕获的是变量，不是简单地保存当时的值。

---

## 快问快答

### Q1：Go 是值传递还是引用传递？

答：

Go 只有值传递。传指针时复制的是指针这个地址值，所以函数内部可以通过地址改到原对象。

### Q2：为什么 `Age` 没变，但 `Score` 变了？

答：

因为结构体被复制了一份，`Age` 改的是副本字段。但 `PtrScore` 里面保存的是同一个地址，解引用后改到了原来的分数。

### Q3：什么是闭包？

答：

闭包就是函数引用了自己外层作用域里的变量，并让这个变量跟着函数继续存在。

---

## 一句话总结

Go 函数的核心是多返回值、显式错误处理和值传递；结构体里有指针字段时，要特别注意浅拷贝带来的共享数据。

