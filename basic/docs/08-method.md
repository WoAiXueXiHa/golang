# Go 方法

## 这一章要记住什么

这份代码主要讲三个点：

- 方法本质上还是函数，只是多了一个接收者。
- Go 没有传统继承，通常用组合来复用字段和方法。
- 方法集决定一个类型能不能实现某个接口，尤其要注意值接收者和指针接收者的区别。

---

## 1. 方法就是带接收者的函数

普通函数是这样：

```go
func Add()
```

方法是这样：

```go
func (c *Counter) AddPtr() {
	c.Count++
}
```

`(c *Counter)` 叫接收者。它表示这个方法属于 `Counter` 这个类型。

可以简单理解成：

```go
func AddPtr(c *Counter) {
	c.Count++
}
```

```text
普通函数：
Add(x, y)
   |
   v
函数独立存在

方法：
c.AddPtr()
|
v
AddPtr 绑定到了 Counter 类型上
```

## 总结一下

方法 = 函数 + 接收者。

接收者决定了这个方法属于哪个类型，也决定了方法里面操作的是值副本，还是原对象地址。

---

## 2. 值接收者和指针接收者

代码里有两个方法：

```go
func (c *Counter) AddPtr() {
	c.Count++
}

func (c Counter) AddValue() {
	c.Count++
}
```

区别是：

- `AddPtr` 用的是指针接收者，改的是原对象。
- `AddValue` 用的是值接收者，改的是副本。

```text
调用 c.AddPtr()

指针接收者拿到地址
+----------+        +----------+
| 地址     | -----> | Count: 1 |
+----------+        +----------+
                         |
                         v
                     Count: 2

调用 c.AddValue()

值接收者拿到副本
+----------+        +----------+
| 原对象   |        | 副本     |
| Count: 2 |        | Count: 2 |
+----------+        +----------+
                         |
                         v
                     副本变 3

原对象还是 Count: 2
```

## 总结一下

值接收者改副本，指针接收者改原对象。

如果方法要修改结构体内部状态，优先用指针接收者。

---

## 3. 编译器的自动转换

代码里这句可以正常运行：

```go
c.AddPtr()
```

虽然 `c` 是值类型，`AddPtr` 要的是 `*Counter`，但 Go 编译器会自动帮你转成：

```go
(&c).AddPtr()
```

```text
值变量调用指针方法：
c.AddPtr()
   |
   v
(&c).AddPtr()

指针变量调用值方法：
pc.AddValue()
   |
   v
(*pc).AddValue()
```

## 总结一下

平时调用方法时，Go 会帮你自动取地址或解引用。

但是判断一个类型是否实现接口时，要看方法集，不能只看“能不能调用”。

---

## 4. 组合代替继承

代码里：

```go
type Engine struct {
	Power int
}

type Car struct {
	Name string
	Engine
}
```

`Car` 直接嵌入了 `Engine`，这叫匿名嵌入。

组合之后，外层可以直接访问内层字段：

```go
fmt.Println(myCar.Power)
```

```text
Car
+----------------------+
| Name: "Audi"         |
| Engine               |
| +------------------+ |
| | Power: 500       | |
| +------------------+ |
+----------------------+

myCar.Power
    |
    v
myCar.Engine.Power
```

`Car` 和 `Engine` 都有 `Start` 方法时：

```text
myCar.Start()
    |
    v
优先找 Car 自己的 Start
    |
    v
调用 Car.Start()

myCar.Engine.Start()
    |
    v
明确调用 Engine.Start()
```

## 总结一下

Go 不靠继承复用代码，而是靠组合。

匿名嵌入后，内层字段和方法会被提升到外层；如果外层有同名方法，优先调用外层自己的方法。

---

## 5. 方法集

代码里有一个接口：

```go
type Payer interface {
	Pay(amount int)
}
```

`WeChatPay` 的 `Pay` 方法是指针接收者：

```go
func (w *WeChatPay) Pay(amount int) {
	w.Balance -= amount
}
```

所以只有 `*WeChatPay` 实现了 `Payer`，`WeChatPay` 本身没有实现。

```text
WeChatPay
+----------------+
| 值接收者方法   |
+----------------+

*WeChatPay
+----------------+
| 值接收者方法   |
| 指针接收者方法 |
+----------------+

Pay 是指针接收者方法
所以只有 *WeChatPay 有 Pay
```

```go
var p2 Payer = &wxWallet
```

可以。

```go
var p1 Payer = wxWallet
```

会报错。

## 总结一下

如果一个接口要求的方法是通过指针接收者实现的，那么只有这个类型的指针实现了接口，值类型本身没有实现。

---

## 易错点

1. `c.AddPtr()` 能调用，不代表 `Counter` 值类型实现了需要指针方法的接口。
2. 值接收者方法里修改字段，不会影响原对象。
3. 指针接收者常用于需要修改对象状态，或者结构体比较大、不想复制的场景。
4. 匿名嵌入不是继承，Go 没有父类子类那套关系。
5. 判断是否实现接口时，看方法集，不看语法糖。

---

## 快问快答

### Q1：Go 方法和函数有什么区别？

答：

方法本质上也是函数，只是多了一个接收者。接收者把函数绑定到某个类型上，所以可以用 `对象.方法()` 的形式调用。

### Q2：值接收者和指针接收者怎么选？

答：

需要修改原对象时用指针接收者。结构体比较大时，也常用指针接收者避免复制。只读且对象很小时，可以用值接收者。

### Q3：为什么 `var p Payer = wxWallet` 会报错？

答：

因为 `Pay` 是定义在 `*WeChatPay` 上的指针接收者方法，所以实现 `Payer` 的是 `*WeChatPay`，不是 `WeChatPay`。

---

## 一句话总结

Go 方法就是带接收者的函数；组合负责能力复用；方法集决定一个类型到底有没有实现接口。

