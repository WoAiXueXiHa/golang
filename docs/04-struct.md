# 04-struct：结构体是 Go 后端的数据骨架

这个目录的重点不是“struct 怎么写”，而是：**Go 用 struct 组织数据，用字段名表达业务含义。** 后端里的请求体、响应体、数据库模型、配置对象，最后都会落到结构体上。

## 关键代码怎么看

`testStruct()` 里定义了一个局部结构体：

```go
type Student struct {
	ID    int
	Name  string
	Age   int
	Score int
}
```

教学代码里定义在函数内部没问题，但真实项目里，结构体一般放在包级别，比如 `model.User`、`dto.CreateUserRequest`、`config.Config`。因为它们通常要被多个函数、多个文件复用。

第一种初始化方式：

```go
st := Student{
	ID:    100,
	Name:  "kun",
	Age:   29,
	Score: 100,
}
```

这是你以后最应该常用的写法。字段名清楚，顺序不敏感，漏填字段也能看出来。

第二种初始化方式：

```go
stu := &Student{
	101,
	"li",
	20,
	99,
}
```

这叫值列表初始化。它依赖字段顺序，字段一变，代码就容易出问题。真实后端里，除了 `Point{X, Y}` 这种极小结构体，基本不推荐这么写。

代码最后有一个小 bug：

```go
fmt.Printf("学生stu的名字: %s\n", st.Name)
```

文字说打印 `stu`，实际访问的是 `st.Name`。应该是 `stu.Name`。这个问题很适合提醒自己：**变量名、输出文案、实际对象必须对齐。后端日志里这种错会非常误导排查。**

## 必须掌握的点

**1. struct 是值类型。**  
把结构体传给函数时，默认复制一份。如果结构体很大，或者函数需要修改它，通常传指针。

**2. 字段名大写才跨包可见。**  
如果以后写 `type User struct { name string }`，其他包访问不到 `name`，JSON 默认也不会导出它。

**3. 优先使用键值对初始化。**  
`Student{Name: "li", Score: 99}` 比 `Student{101, "li", 20, 99}` 更稳。

**4. `&Student{...}` 返回指针。**  
Go 后端里很多构造函数会返回指针：

```go
func NewStudent(name string) *Student {
	return &Student{Name: name}
}
```

## 用一个形象例子理解

值列表初始化像“按座位号点名”：第一个位置是 ID，第二个位置是 Name。只要中间多加一个座位，后面的人全乱。

键值对初始化像“按身份证点名”：你明确写 `Name: "li"`，无论字段顺序怎么调整，它都知道你要填的是姓名。

## 和 Go 后端开发的关系

真实项目里 struct 常见形态：

```go
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
```

你要逐渐建立这个习惯：**不同层的数据结构承担不同职责**。请求结构体不一定等于数据库结构体，响应结构体也不一定直接暴露数据库字段。

## 更像工程代码的写法

- `Student` 提到包级别，方便复用。
- 初始化 `stu` 时也用键值对：

```go
stu := &Student{
	ID:    101,
	Name:  "li",
	Age:   20,
	Score: 99,
}
```

- 打印结构体调试时用 `%+v`，能看到字段名。
- 修正 `st.Name` 为 `stu.Name`。
- 后续学习 JSON、数据库时，重点看 struct tag。

## 复习时问自己

1. **为什么不推荐 `Student{101, "li", 20, 99}`？**  
   因为它依赖字段顺序，字段变化时容易错位。

2. **`Student{}` 的零值是什么？**  
   每个字段取自己的零值，数字是 0，字符串是 `""`。

3. **`stu := &Student{...}` 和 `st := Student{...}` 最大区别是什么？**  
   前者是 `*Student` 指针，后者是 `Student` 值。

4. **真实后端为什么常写 struct tag？**  
   用来控制 JSON 字段名、数据库列映射、参数校验等。
