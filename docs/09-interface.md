# 09-interface：接口不是继承，是工程解耦

这一章是 Go 后端核心。  
一句话：**接口描述能力，调用方依赖能力，不依赖具体实现。**

## 这份代码最该看什么

```go
type Phone interface {
	Call()
	SendMessage()
}
```

接口只关心方法，不关心字段。  
`Apple` 和 `Oppo` 只要拥有这两个方法，就自动实现 `Phone`。

```go
func (this *Apple) Call() {}
func (this *Apple) SendMessage() {}
```

没有 `implements`。这是 Go 接口的灵魂：**隐式实现**。

```go
var phoneC Phone
phoneC = new(Apple)
phoneC.Call()
```

接口变量里装的是具体实现。调用 `phoneC.Call()` 时，运行时会找到实际类型 `*Apple` 的方法。

如果要访问 `PhoneName`：

```go
phoneC.(*Apple).PhoneName = "APPLE"
```

必须类型断言。因为 `Phone` 接口没有承诺自己有 `PhoneName` 字段。

更清楚的写法是：

```go
var phone Phone = &Apple{PhoneName: "APPLE"}
```

少一次断言，意图更直接。

## 接口的本质

接口不是“父类”。接口是“能力清单”。

通俗案例：  
招聘要求写“会开车、能送货”。不管你开面包车还是电动车，只要能完成这两个动作，就符合岗位。Go 不要求你提前报名“我实现了送货员接口”，你真有这些方法，编译器就认。

这和 Java 最大区别：

- Java：类型声明自己实现了某接口
- Go：接口看方法集合，方法对上就实现

## 小接口才是好接口

工程里不要一上来定义巨型接口：

```go
type UserRepo interface {
	Create(...)
	Update(...)
	Delete(...)
	Find(...)
	List(...)
	Count(...)
}
```

如果某个服务只需要查询，就定义小接口：

```go
type UserFinder interface {
	FindByID(ctx context.Context, id int64) (*User, error)
}
```

小接口的好处：

- 更容易 mock
- 更容易复用
- 调用方依赖更少
- 代码边界更清楚

Go 标准库就是这么做的：`io.Reader` 只有一个方法。

## 接口应该由谁定义

Go 里常见原则：**接口由使用方定义。**

Service 需要什么能力，就在 service 侧定义：

```go
type UserRepo interface {
	FindByID(ctx context.Context, id int64) (*User, error)
}

type UserService struct {
	repo UserRepo
}
```

MySQL 实现不需要知道这个接口存在：

```go
type MySQLUserRepo struct {}

func (r *MySQLUserRepo) FindByID(ctx context.Context, id int64) (*User, error) {
	// query db
}
```

只要方法匹配，它就自动满足 `UserRepo`。

这就是依赖倒置：业务层不依赖 MySQL 细节，只依赖自己需要的能力。

## 空接口 any：能不用就别用

代码里：

```go
var any interface{}
any = 10
any = "Vect"
any = map[string]int{"aa": 1}
```

`interface{}` 可以接收任何类型。Go 1.18 后 `any` 是它的别名。

问题是：类型信息丢给运行时了。你要拿回来，就得断言：

```go
v, ok := any.(string)
```

工程建议：

- 能用具体类型，就用具体类型
- 能用小接口，就用小接口
- 能用泛型，就别滥用 any
- any 适合日志、JSON 任意结构、通用容器边界

## 类型断言和 type switch

安全断言：

```go
v, ok := x.(int)
if !ok {
	return
}
```

不安全断言：

```go
v := x.(int) // 类型不对直接 panic
```

多类型判断用 type switch：

```go
switch v := x.(type) {
case int:
	fmt.Println(v)
case string:
	fmt.Println(v)
default:
	fmt.Println("unknown")
}
```

面试官问断言，重点不是语法，而是你知道 panic 风险。

## nil 接口：经典拷打点

接口值内部可以粗略理解成两部分：

- 动态类型
- 动态值

只有两者都为空，接口才是 nil。

经典坑：

```go
type MyErr struct{}

func (e *MyErr) Error() string { return "bad" }

func f() error {
	var e *MyErr = nil
	return e
}

err := f()
fmt.Println(err == nil) // false
```

为什么？  
返回的接口里有动态类型 `*MyErr`，动态值是 nil。接口本身不是 nil。

工程结论：返回 error 时，不要把 nil 的具体指针塞进 error 接口。没错误就直接 `return nil`。

## 接口和方法接收者

如果接口要求：

```go
type Reader interface {
	Read() int
}
```

实现方法是：

```go
func (r *MyReader1) Read() int {
	return r.a + r.b
}
```

那么通常是 `*MyReader1` 实现了 `Reader`，不是 `MyReader1` 值。

```go
var r Reader
r = &MyReader1{2, 10} // ok
```

这和方法集有关。简化记：

- 值接收者方法：`T` 和 `*T` 都能满足接口
- 指针接收者方法：通常只有 `*T` 满足接口

## 接口嵌套

代码里：

```go
type V interface {
	A
	B
	run3()
}
```

意思是实现 `V` 必须同时拥有 `A`、`B` 和 `run3` 的方法。

标准库例子：

```go
type ReadWriter interface {
	Reader
	Writer
}
```

接口嵌套不是继承，是组合。把小能力拼成大能力。

## 后端实践怎么用

**Repository 抽象：**

```go
type UserRepo interface {
	FindByID(ctx context.Context, id int64) (*User, error)
}
```

**Service 依赖接口：**

```go
type UserService struct {
	repo UserRepo
}
```

**测试替身：**

```go
type fakeUserRepo struct{}

func (f fakeUserRepo) FindByID(ctx context.Context, id int64) (*User, error) {
	return &User{ID: id}, nil
}
```

这样单元测试不用连真实数据库。

## 本目录代码可改进点

- `this` 改成 Go 惯用接收者名：`a *Apple`、`o *Oppo`。
- `MyWriter` 接口的方法也叫 `MyWriter`，工程里更自然是 `Write`。
- `DoJob1(val interface{})` 可写成 `DoJob1(val any)`。
- `runer` 应为 `runner`。
- `phoneC = new(Apple)` 后再断言赋值不够自然，直接 `&Apple{PhoneName: "APPLE"}`。

## 面试拷打

1. **Go 接口是显式实现还是隐式实现？**  
   隐式。方法集合匹配就实现。

2. **接口的核心价值是什么？**  
   解耦。调用方依赖能力，不依赖具体类型。

3. **为什么推荐小接口？**  
   更容易复用、mock、替换实现，依赖面更小。

4. **接口一般定义在实现方还是使用方？**  
   Go 更推荐使用方定义接口。

5. **`interface{}` 和 `any` 有区别吗？**  
   没有，`any` 是别名。

6. **类型断言不带 ok 会怎样？**  
   类型不匹配会 panic。

7. **nil 指针放进接口后，接口一定等于 nil 吗？**  
   不一定。接口有动态类型时，哪怕动态值是 nil，接口也不是 nil。

8. **指针接收者方法对接口实现有什么影响？**  
   通常只有指针类型 `*T` 实现接口，值类型 `T` 不实现。
