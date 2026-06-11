# 08-method：方法是行为归属，接收者决定语义

这一章的本质：**函数解决动作，方法说明动作属于谁。**  
Go 没有 class，但可以给类型绑定方法；没有继承，但可以用组合复用能力。

## 这份代码最该看什么

```go
type People struct {
	Name string
	Age  int
}

func (this *People) GetName() string {
	return this.Name
}
```

`GetName` 是方法。它和普通函数的区别就是多了接收者：

```go
func (p *People) GetName() string
```

接收者就是这个方法“归属”的对象。  
工程里不要用 `this`，Go 惯用短名：`p *People`、`s *Student`、`h *Handler`。

## 值接收者 vs 指针接收者

代码里：

```go
func (st Student) GetScore() float64 {
	return st.Score
}

func (this *Student) SetScore(score float64) {
	this.Score = score
}
```

值接收者拿到副本。适合小对象、只读、不关心复制成本的场景。  
指针接收者拿到地址。适合修改对象、避免复制、保持状态一致。

选择原则很简单：

- 要修改接收者：用指针
- 结构体较大：用指针
- 包含 `sync.Mutex`：必须用指针
- 同一类型已有指针方法：尽量统一用指针
- 小且不可变语义：可以用值

面试会问：**值变量能调用指针接收者方法吗？**

```go
st.SetScore(99.5)
```

可以。因为 `st` 是可寻址变量，编译器自动转成：

```go
(&st).SetScore(99.5)
```

但不是所有值都可寻址。比如 map 元素不行：

```go
students["a"].SetScore(90) // 如果 SetScore 是指针接收者，这通常不行
```

因为 map 元素取出来不是稳定地址。

## 嵌入不是继承

```go
type Student struct {
	ID    int
	Score float64
	People
}
```

`Student` 里匿名嵌入了 `People`。Go 会把 `People` 的字段和方法提升上来：

```go
st.Name
st.GetName()
```

但本质仍然是：

```go
st.People.Name
st.People.GetName()
```

这叫组合，不是继承。  
`Student` 不是 `People`，它只是包含了 `People`。

通俗理解：继承像“血缘关系”，组合像“装备关系”。`Student` 装备了 `People` 的字段和方法，但它没有变成 `People`。

## 方法覆盖和提升

代码里 `People` 和 `Student` 都有 `GetName`：

```go
func (this *People) GetName() string
func (this *Student) GetName() string
```

调用：

```go
st.GetName()
```

优先调用外层 `Student.GetName`。如果要调用内部的：

```go
st.People.GetName()
```

工程里不要为了“像继承”而乱覆盖方法。覆盖应该表达明确差异，否则只是增加理解成本。

## 方法集：接口拷打重点

方法接收者会影响类型是否实现接口。

```go
type Namer interface {
	GetName() string
}
```

如果方法是：

```go
func (p *People) GetName() string
```

那么通常是 `*People` 实现了 `Namer`，不是 `People` 值本身。

这点在接口赋值时经常被问：

```go
var n Namer
var p People
n = p  // 可能不行
n = &p // 可以
```

简化记忆：

- 值接收者方法：值和指针都能用
- 指针接收者方法：通常只有指针满足接口

## 后端实践怎么用

**Handler：**

```go
type Handler struct {
	svc *UserService
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// h.svc.FindUser(...)
}
```

方法让 handler 携带依赖，而不是到处用全局变量。

**Service：**

```go
type UserService struct {
	repo UserRepo
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserReq) error {
	return s.repo.Create(ctx, req)
}
```

方法把业务行为收拢到对象上，依赖关系也更清楚。

**Model 嵌入基础字段：**

```go
type BaseModel struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	BaseModel
	Name string
}
```

这就是组合的实际用途。但不要搞多层嵌入迷宫，清楚比炫技重要。

## 本目录代码可改进点

- `this` 改为 `p`、`s`。
- `Student.GetName` 和 `People.GetName` 内容一样，没必要覆盖，除非要展示差异。
- `Student` 的方法建议统一指针接收者。
- 可以补一个接口示例，展示接收者类型如何影响接口实现。

## 面试拷打

1. **方法和函数区别？**  
   方法有接收者，函数没有。方法本质上也是函数，只是语法上绑定到类型。

2. **值接收者和指针接收者怎么选？**  
   修改、大对象、含锁、保持一致性用指针；小对象只读可以用值。

3. **Go 的嵌入是继承吗？**  
   不是，是组合加字段/方法提升。

4. **指针接收者方法，值能不能调用？**  
   可寻址的值可以，编译器自动取地址。

5. **指针接收者对接口实现有什么影响？**  
   如果方法只定义在 `*T` 上，通常只有 `*T` 实现接口，`T` 不实现。

6. **为什么含 `sync.Mutex` 的 struct 方法必须用指针接收者？**  
   因为复制 Mutex 会破坏锁语义，可能导致并发问题。
