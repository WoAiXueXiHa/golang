# Go 接口

## 这一章要记住什么

这一章主要讲六个点：

- 接口定义的是一组方法约束。
- 一个类型只要实现了接口要求的所有方法，就自动实现这个接口。
- 一个类型可以同时实现多个接口。
- 空接口 `interface{}` 可以接收任意类型。
- 类型断言可以把接口里的真实值取出来。
- 接口可以嵌套接口。

---

## 1. 接口基本使用

代码里：

```go
type Phone interface {
	Call()
	SendMessage()
}
```

这个接口要求实现者必须有两个方法：

```text
Phone 接口
+---------------+
| Call()        |
| SendMessage() |
+---------------+
```

`Apple` 和 `Oppo` 都实现了这两个方法：

```go
func (this *Apple) Call() {}
func (this *Apple) SendMessage() {}
```

所以 `*Apple` 可以当作 `Phone` 使用。

## 总结一下

接口不关心你是什么类型，只关心你有没有实现它要求的方法。

---

## 2. 接口和多态

代码里：

```go
var phoneC Phone
phoneC = new(Apple)
phoneC.Call()
phoneC.SendMessage()
```

`phoneC` 是接口变量，里面实际装的是 `*Apple`。

```text
phoneC interface
+----------------+
| 动态类型 *Apple |
| 动态值 Apple地址|
+----------------+
        |
        v
调用 Apple 的方法
```

接口变量调用方法时，会根据里面装的真实类型执行对应方法。

## 总结一下

多态就是同一个接口变量，可以装不同的实现类型，然后用统一的方法调用。

---

## 3. 一个类型实现多个接口

代码里：

```go
type MyWriter interface {
	MyWriter(s string)
}

type MyRead interface {
	MyReader()
}
```

`MyReadWriter` 同时实现了两个方法：

```go
func (this *MyReadWriter) MyWriter(s string) {}
func (this *MyReadWriter) MyReader() {}
```

所以它可以同时满足两个接口。

```text
MyReadWriter
    |
    +-- 实现 MyReader() -> 满足 MyRead
    |
    +-- 实现 MyWriter() -> 满足 MyWriter
```

## 总结一下

Go 的接口是隐式实现的，不需要写 `implements`。

一个类型可以自然地满足多个小接口。

---

## 4. 空接口

代码里：

```go
var any interface{}
any = 10
any = "Vect"
any = map[string]int{"aa": 1}
```

空接口没有要求任何方法，所以任何类型都满足它。

```text
interface{}
+----------------+
| 没有方法要求   |
+----------------+
        |
        v
任意类型都能放进去
```

## 总结一下

`interface{}` 表示“我可以接收任意类型”。

现在新代码里也经常用 `any`，它本质上就是 `interface{}` 的别名。

---

## 5. 类型断言

代码里：

```go
var x interface{}
x = 9
val, ok := x.(int)
```

断言的意思是：我认为接口里装的是 `int`，请帮我取出来。

```text
x interface{}
+------------+
| type: int  |
| value: 9   |
+------------+
      |
      v
x.(int) 成功
```

推荐写法是带 `ok`：

```go
val, ok := x.(int)
```

如果不带 `ok`，断言失败会直接 `panic`。

## 总结一下

从接口变量里取具体类型时，用类型断言。

不确定类型时，一定用 `value, ok := x.(T)`。

---

## 6. 接口作为函数参数

代码里：

```go
type Reader interface {
	Read() int
}

func DoJob(r Reader) {
	fmt.Printf("myReader is %d\n", r.Read())
}
```

`DoJob` 不关心传进来的具体类型，只关心它有没有 `Read()` 方法。

```text
DoJob 需要 Reader
       |
       v
只要有 Read() int 就能传进来
```

## 总结一下

接口作为参数，可以让函数依赖抽象，而不是依赖具体结构体。

---

## 7. 接口嵌套

代码里：

```go
type A interface {
	run1()
}

type B interface {
	run2()
}

type V interface {
	A
	B
	run3()
}
```

要实现 `V`，就必须同时实现：

```text
run1()
run2()
run3()
```

```text
V 接口
+--------+
| A      | -> run1()
| B      | -> run2()
| run3() |
+--------+
```

## 总结一下

接口嵌套就是把多个接口组合成一个更大的接口。

实现外层接口时，必须满足里面所有方法。

---

## 易错点

1. 接口是隐式实现，不需要显式声明。
2. 指针接收者方法会影响到底是值类型实现接口，还是指针类型实现接口。
3. 空接口可以接收任意类型，但取出来时需要类型断言。
4. 不带 `ok` 的类型断言失败会 `panic`。
5. 接口越小越灵活，Go 里很常见小接口设计。

---

## 快问快答

### Q1：Go 接口和 Java 接口最大的区别是什么？

答：

Go 接口是隐式实现的。一个类型只要实现了接口要求的方法，就自动满足这个接口，不需要写 `implements`。

### Q2：空接口为什么能接收任意类型？

答：

因为空接口没有任何方法要求，所有类型都天然满足它。

### Q3：类型断言为什么要带 `ok`？

答：

因为断言失败时，如果不带 `ok` 会直接 `panic`。带 `ok` 可以安全判断是否成功。

### Q4：接口作为参数有什么好处？

答：

函数只依赖行为，不依赖具体类型，代码更容易扩展和测试。

---

## 一句话总结

接口定义行为，类型隐式实现行为；空接口接收任意值，类型断言负责把真实类型取出来。

