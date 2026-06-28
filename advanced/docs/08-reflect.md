# reflect 包：运行时反射

## 这一章要记住什么

- **反射就是程序在运行时检测自身类型和值的能力** — 两个入口函数：`reflect.TypeOf()` 返回类型信息、`reflect.ValueOf()` 返回值信息
- **Type 是接口，Value 是结构体** — TypeOf 只回答"是什么类型"；ValueOf 除了类型还承载具体值，而可以修改值也需要从 Value 入手
- **Type.Name() 和 Type.Kind() 不同** — Name 是 Go 层面的类型名（自定义类型有自己的名字），Kind 是底层基础分类（int、struct、ptr 等），写通用框架时用 Kind 做分支
- **指针反射必须用 Elem() 解引用** — `TypeOf(&st).Kind() == ptr`，必须 `.Elem()` 后才能拿到 struct 的字段和方法
- **反射修改值三个条件：传指针 → Elem() → CanSet 为 true** — 直接 `ValueOf(x)` 拿到的是副本，不可改；只有 `ValueOf(&x).Elem()` 才能修改原值
- **反射有代价** — 性能比直接访问慢 1~2 个数量级，编译期类型安全丢失，容易写出运行时 panic 的代码

---

## 0. 什么是反射

Go 的变量在编译时类型是确定的。但有一个特殊类型 `interface{}`（或 `any`），它可以接收任意类型的值。当一个具体类型赋值给 `interface{}` 后，这个接口变量内部其实有两层信息：

```
interface{} 变量（如 var x any = Stu{...}）

+-----------------------------+
| 静态类型 (编译期确定)         |  → 永远是 interface{} / any
+-----------------------------+
| 动态类型 (运行时确定)         |  → 被赋值的具体类型，这里是 Stu
+-----------------------------+
| 动态值   (运行时确定)         |  → 具体的 {Name, Age, Score}
+-----------------------------+
```

**反射就是在运行时读取（甚至修改）这个动态类型和动态值。** Go 的 `reflect` 包提供了两个入口函数：

- `reflect.TypeOf(x any) reflect.Type` — 拿到动态类型，返回 `reflect.Type` 接口
- `reflect.ValueOf(x any) reflect.Value` — 拿到动态值（同时也包含类型），返回 `reflect.Value` 结构体

几乎所有反射操作都以这两个函数为起点。

### 总结一下

反射让程序在运行时知道自己是什么类型、有什么字段、能调什么方法。TypeOf 回答"是什么"，ValueOf 回答"有什么值"。底层依赖的就是 `interface{}` 存储动态类型/动态值的机制。

---

## 1. reflect.TypeOf() 与 reflect.ValueOf()

两个函数的签名：

```go
func TypeOf(i any) Type
func ValueOf(i any) Value
```

都接收 `any`（空接口），这意味着可以把任何值传进去。区别在于返回：

| 函数 | 返回类型 | 本质 | 能做什么 |
|------|---------|------|---------|
| `TypeOf` | `reflect.Type`（接口） | 类型描述符 | 查字段名、方法签名、Kind |
| `ValueOf` | `reflect.Value`（结构体） | 类型 + 值的容器 | 上述所有 + 读取值、修改值 |

关键代码片段：

```go
var num int64 = 100

// TypeOf：只拿到类型描述
t := reflect.TypeOf(num)
fmt.Println(t.String()) // "int64"

// ValueOf：拿到值（也包含类型）
v := reflect.ValueOf(num)
fmt.Println(v.Int())        // 100
fmt.Println(v.Type())       // int64 —— Value 也能拿到 Type
```

简单说：**TypeOf 只关心类型，ValueOf 关心类型 + 值。** 如果需要读/写具体值，必须用 ValueOf。而且 Value 身上也有 `.Type()` 方法，所以拿到 Value 就同时拥有了类型和值。

### 总结一下

`TypeOf` 和 `ValueOf` 是 reflect 的两个入口。TypeOf 返回类型描述（纯元数据），ValueOf 返回类型 + 值容器。大部分反射操作都需要从 Value 出发，因为它同时包含类型和值的信息。

---

## 2. Type 与 Kind：名字和底层分类

`reflect.Type` 有两个容易混淆的方法：

- **Name()** — 返回 Go 语言层面的类型名。内置类型（如 `int`）有名字；自定义类型（如 `type WrapInt int`）也有自己的名字。
- **Kind()** — 返回底层基础分类（`reflect.Kind` 枚举：`Int`、`String`、`Struct`、`Ptr`、`Slice` 等）。

以自定义类型为例：

```go
type WrapInt int

var a int = 100
var b WrapInt = 200

t1 := reflect.TypeOf(a)
fmt.Println(t1.Name()) // "int"
fmt.Println(t1.Kind()) // int

t2 := reflect.TypeOf(b)
fmt.Println(t2.Name()) // "WrapInt"  ← 自己的名字
fmt.Println(t2.Kind()) // int        ← 底层和 int 一样
```

Name 不同，Kind 相同。写通用反射代码（比如"遍历所有 int 类型的字段"）时，应该用 **Kind** 做条件判断，而不是 Name——因为你不知道使用代码的人会不会定义一个 `type MyInt int`。

```go
switch t.Kind() {
case reflect.Int, reflect.Int64:
    // 处理所有整数类字段
case reflect.String:
    // 处理字符串字段
case reflect.Struct:
    // 递归处理嵌套结构体
}
```

### 总结一下

Name 是类型在 Go 代码里被称呼的名字（包括自定义类型名），Kind 是底层是"哪一类"（int、struct、ptr 等）。写通用框架时用 Kind 判断，因为自定义类型的 Kind 和底层类型一致。

---

## 3. 通过反射遍历结构体字段值

拿到 struct 的 `reflect.Value` 后，可以遍历每个字段：

```go
st := Stu{Name: "zhangsan", Age: 18, Score: 95.5}
v := reflect.ValueOf(st)

for i := 0; i < v.NumField(); i++ {
    field := v.Field(i)
    fmt.Printf("字段%d: 类型=%v, 值=%v\n", i+1, field.Type().Name(), field.Interface())
}
```

关键方法：

| 方法 | 说明 |
|------|------|
| `Value.NumField()` | 字段个数 |
| `Value.Field(i)` | 第 i 个字段的 `reflect.Value` |
| `Value.Interface()` | 把字段值还原为 `interface{}` |

> **推荐用 `.Interface()` 取值**，而不是 `.Int()`、`.String()` 等类型化方法。因为如果字段类型不匹配，类型化方法会直接 panic。`.Interface()` 配合 `%v` 打印最安全。

### 总结一下

反射遍历 struct 字段用 `NumField + Field(i)`，取值用 `Interface()` 最安全。这也就是 `json.Marshal` 内部做的事情——它用反射遍历所有字段、读取值、拼成 JSON。

---

## 4. 通过反射操作 Map 和 Slice

**Map** 的反射遍历：

```go
m := map[int]uint32{1: 100, 2: 200}
v := reflect.ValueOf(m)

for _, k := range v.MapKeys() {
    val := v.MapIndex(k)
    fmt.Printf("key=%d, value=%d\n", k.Int(), val.Uint())
}
```

- `Value.MapKeys()` — 返回所有 key 的 `[]reflect.Value`
- `Value.MapIndex(key)` — 根据 key 取对应的 value

**Slice/Array** 的反射遍历（两者用同一套方法）：

```go
s := []int{1, 2, 3}
v := reflect.ValueOf(s)

for i := 0; i < v.Len(); i++ {
    fmt.Printf("%v ", v.Index(i).Interface())
}
```

- `Value.Len()` — 长度
- `Value.Index(i)` — 第 i 个元素

slice 和 array 在反射层面用同一套 `Len` + `Index`，不需要区分。

### 总结一下

Map 用 `MapKeys` + `MapIndex` 遍历，Slice/Array 用 `Len` + `Index` 遍历。反射抹平了 slice 和 array 的差异，统一用下标操作。

---

## 5. 指针反射与 Elem()

对指针直接调 `TypeOf` 或 `ValueOf`，返回的 Kind 是 `ptr`，不能直接访问字段/方法。必须用 **`Elem()`** 解引用：

```
reflect.TypeOf(&st)

+--------+       .Elem()        +---------+
|  *Stu  | ──────────────────→ |  Stu    |
|  Kind= |  ptr 解引用          |  Kind=  | struct
+--------+                     +---------+
                                   │
                                   ├─ .NumField() → 3
                                   ├─ .Field(0).Name → "Name"
                                   └─ .Name() → "Stu"
```

代码示例：

```go
st := &Stu{Name: "zhangsan", Age: 18}
t := reflect.TypeOf(st)

fmt.Println(t.Kind())        // ptr —— 不能直接在 ptr 上调 NumField
fmt.Println(t.Elem().Kind()) // struct —— 解引用后才能访问字段
fmt.Println(t.Elem().NumField()) // 3
```

`Elem()` 在类型反射和值反射中都有：
- `Type.Elem()` — 返回指针指向的类型
- `Value.Elem()` — 返回指针指向的值（同时是**可修改**的，见后面第 8 节）

### 总结一下

指针对应的反射 Kind 是 `ptr`，拿不到字段和方法。`.Elem()` 相当于解引用，之后就能看字段、调方法、甚至修改值。忘了 `Elem()` 是最常见的反射错误之一。

---

## 6. 函数反射

反射可以查看函数的签名（参数和返回值类型）：

```go
func Add(a, b int) (int, error) { return a + b, nil }

t := reflect.TypeOf(Add)
fmt.Println(t.NumIn())  // 入参个数: 2
fmt.Println(t.In(0))    // 第1个入参类型: int
fmt.Println(t.NumOut()) // 返回值个数: 2
fmt.Println(t.Out(1))   // 第2个返回值类型: error
```

常用方法：

| 方法 | 说明 |
|------|------|
| `Type.NumIn()` | 入参个数 |
| `Type.In(i)` | 第 i 个入参的 `reflect.Type` |
| `Type.NumOut()` | 返回值个数 |
| `Type.Out(i)` | 第 i 个返回值的 `reflect.Type` |

函数反射主要用于框架层面——比如依赖注入容器需要知道构造函数的参数类型才能自动创建依赖。

### 总结一下

反射可以拿到函数的完整签名（几个参数、什么类型、几个返回值）。这为动态调用函数提供了基础——知道签名后，就可以用 `Value.Call()` 传参调用了。

---

## 7. 获取和调用结构体方法

**获取方法列表：**

```go
c := Calculator{}         // 有 Add 和 Mul 两个方法
t := reflect.TypeOf(c)

for i := 0; i < t.NumMethod(); i++ {
    m := t.Method(i)
    fmt.Println(m.Name) // 方法名
    fmt.Println(m.Type) // 完整签名（含接收者）
}
```

注意 `m.Type` 包含接收者作为第一个参数，比如 `func(main.Calculator, int, int) int`。

**调用方法：**

```go
v := reflect.ValueOf(c)

// 按名称获取方法
method := v.MethodByName("Add")

// 参数必须包装成 []reflect.Value
args := []reflect.Value{
    reflect.ValueOf(3),
    reflect.ValueOf(4),
}

// 调用，返回 []reflect.Value
results := method.Call(args)
fmt.Println(results[0].Interface()) // 7
```

几个要点：
- `MethodByName(name)` 返回绑定到该接收者的方法值
- `Call(args)` 的参数和返回值都是 `[]reflect.Value`
- 参数和返回值需要手动包装/解包（`ValueOf` 包装、`Interface()` 解包）

### 总结一下

反射可以列出 struct 的所有导出方法（`NumMethod` + `Method`），也可以按名字动态调用（`MethodByName` + `Call`）。这也就是 RPC 框架的核心机制——收到方法名和参数，反射找到方法，传入参数，拿回结果。

---

## 8. 通过反射修改值

这是反射最容易出错的地方。修改值必须满足三个条件：**传指针、调用 Elem()、CanSet 为 true**。

```
直接传值：                        传指针 + Elem()：
                                 
reflect.ValueOf(x)              reflect.ValueOf(&x).Elem()
       │                              │
       ▼                              ▼
  ┌─────────┐                   ┌─────────┐
  │ 副本    │                    │ 原对象  │ ← 可修改
  │ CanSet: │ false              │ CanSet: │ true
  └─────────┘                    └─────────┘
```

代码演示：

```go
x := 100

// ❌ 错误：直接传值，CanSet = false
v1 := reflect.ValueOf(x)
fmt.Println(v1.CanSet()) // false
// v1.SetInt(200) ← panic!

// ❌ 错误：传了指针但没用 Elem，CanSet 还是 false
v2 := reflect.ValueOf(&x)
fmt.Println(v2.CanSet()) // false（指针本身不可改）

// ✅ 正确：传指针 + Elem()
v3 := reflect.ValueOf(&x).Elem()
fmt.Println(v3.CanSet()) // true
v3.SetInt(200)           // 成功！x 现在是 200
```

修改 struct 字段：

```go
st := Stu{Name: "zhangsan", Age: 18}
v := reflect.ValueOf(&st).Elem()

v.Field(0).SetString("lisi") // 改 Name
v.Field(1).SetInt(20)        // 改 Age
// st 现在是 {lisi 20}
```

常用 Set 方法：`SetInt`、`SetString`、`SetFloat`、`SetBool`、`Set`（设置任意 reflect.Value）。

### 总结一下

反射修改值的核心公式：`ValueOf(&x).Elem().SetXxx(v)`。没有指针就没有 CanSet，没有 Elem 也拿不到 CanSet。这是反射里最容易踩的坑，但记住三步骤就不会出错。

---

## 9. 结构体标签（StructTag）

StructTag 是编译期写在 struct 字段上的 `key:"value"` 字符串。标准库 `encoding/json`、`gopkg.in/yaml.v3` 等都靠它来映射字段名：

```go
type Config struct {
    Name string `json:"name" yaml:"name"`
    Port int    `json:"port" yaml:"port"`
}
```

反射读取标签：

```go
cfg := Config{Name: "server", Port: 8080}
t := reflect.TypeOf(cfg)

for i := 0; i < t.NumField(); i++ {
    f := t.Field(i) // reflect.StructField

    // 方式1：Get —— 找不到返回空字符串
    jsonName := f.Tag.Get("json")  // "name"

    // 方式2：Lookup —— 返回 (value, ok)，能区分"空值"和"不存在"
    yamlName, ok := f.Tag.Lookup("yaml") // ("name", true)
    _, notFound := f.Tag.Lookup("xml")   // ("", false)
}
```

`StructTag` 本质是 `string` 类型，有两个方法：
- `Get(key)` — 查单个 key，找不到返回 `""`
- `Lookup(key)` — 返回 `(value, bool)`，能区分"值是空字符串"和"key 不存在"

### 总结一下

StructTag 是编译期关联到字段的 key:value 元数据，通过反射的 `Field(i).Tag.Get()` / `Tag.Lookup()` 在运行时读取。几乎所有序列化库（json、yaml、xml、protobuf）都依赖这个机制。

---

## 10. 反射的优缺点

### 优点

1. **动态分派**：可以根据字符串名字找到方法并调用，RPC 框架、HTTP 路由分发都靠这个
2. **减少重复代码**：一套反射逻辑处理任意 struct，比如 `json.Marshal` 不需要为每种类型写一套序列化代码
3. **通用工具的基础**：ORM、依赖注入、配置绑定、测试框架都依赖反射

### 缺点

1. **性能开销**：反射涉及堆逃逸（值要复制到 `interface{}` 中）、运行时类型检查，比直接访问慢 1~2 个数量级
2. **编译期类型安全丢失**：字段名写错、类型不匹配都不会编译报错，只会在运行时 panic
3. **代码可读性下降**：反射代码比直接访问冗长、意图不直观
4. **容易 panic**：CanSet 没检查就 Set、对 string 字段调 Int()、类型不匹配的 Set 都会 panic

### 什么时候用、什么时候不用

| 场景 | 建议 |
|------|------|
| 序列化/反序列化库 | ✅ 用反射，标准库 json 就这么做的 |
| ORM 映射 struct → 数据库 | ✅ 用反射 |
| RPC/路由分发 | ✅ 用反射 |
| 普通业务代码 | ❌ 别用，直接访问即可 |
| 性能热点路径 | ❌ 别用，考虑代码生成 |
| 能用泛型/接口解决 | ❌ 优先用泛型/接口 |

原则：**能用泛型、接口、代码生成解决的，不用反射。**

---

## 易错点

1. **TypeOf 和 ValueOf 混淆** — `TypeOf` 返回类型信息不能取值，`ValueOf` 除了类型信息还持有值。需要读/写具体值的时候一定是 `ValueOf`
2. **Name 和 Kind 区分不清** — 自定义类型 `type MyInt int` 的 Name 是 `"MyInt"`，Kind 是 `int`。写 switch 做类型分支时用 Kind，不要用 Name 做字符串比较
3. **指针反射忘了 Elem()** — `TypeOf(&st).NumField()` 拿不到字段，必须先 `.Elem()`。这是最常见的反射错误
4. **直接 ValueOf(x) 就想 Set 值** — `ValueOf(x).CanSet()` 永远是 false，因为传值拿到的是副本。必须 `ValueOf(&x).Elem()` 才能 Set
5. **用错取值方法** — 对 string 字段调 `.Int()`、对 int 调 `.String()` 都会 panic。不确定类型时用 `.Interface()` 最安全
6. **CanSet 没检查就直接 Set** — 反射修改值前务必检查 `CanSet()`，否则 production 里 panic 很难排查
7. **NumMethod 数量与预期不符** — `NumMethod()` 只返回导出方法；如果方法定义在指针接收者上，值类型可能找不到该方法
8. **Call 参数没包装成 reflect.Value** — `Call()` 要求 `[]reflect.Value` 类型的参数，不能直接传原始值
9. **StructTag.Get 区分不了空值和不存在** — 当 tag 值本身可能是空字符串时，用 `Lookup` 代替 `Get`，通过第二个返回值判断 key 是否存在
10. **`reflect.TupeOf` 拼写错误** — 是 `TypeOf` 不是 `TupeOf`，少一个字母编译不过

---

## 快问快答

### Q1：`reflect.TypeOf` 和 `reflect.ValueOf` 有什么区别？

答：TypeOf 返回 `reflect.Type` 接口，只包含类型信息（名字、字段、方法签名）。ValueOf 返回 `reflect.Value` 结构体，包含类型信息 + 具体值，还能修改值（如果可寻址的话）。大多数反射操作从 ValueOf 出发，因为 Value 上也有 `.Type()` 方法。

### Q2：Type 和 Kind 有什么区别？什么时候用 Kind？

答：Type.Name() 是 Go 代码层面的类型名（如 `"int"`、`"Stu"`、`"WrapInt"`）；Kind() 是底层分类（Int、Struct、Ptr 等）。自定义类型 Name 是自己的名字但 Kind 和底层类型一致。写通用框架做 switch 分支时用 Kind，因为自定义类型的 Kind 和基础类型一样。

### Q3：为什么直接 `ValueOf(x)` 不能修改值？怎么才能改？

答：`ValueOf(x)` 传入的是 x 的副本，反射拿的是副本的 `reflect.Value`，改副本没意义，所以 Go 直接设置 CanSet=false 禁止修改。要改原值，必须 `ValueOf(&x).Elem()`——传指针拿到的是指针的反射值，再 Elem() 解引用得到指向原对象的可写反射值。

### Q4：`Elem()` 在反射里有什么用？

答：Elem() 就是反射层面的"解引用"。对 Kind 为 Ptr 的 Type/Value 调用 Elem()，返回指针指向的类型/值。没有 Elem() 就看不到 struct 的字段、拿不到可修改的值。**指针对应的反射对象里，所有有用信息都在 Elem() 之后。**

### Q5：StructTag 的工作原理是什么？

答：StructTag 是定义 struct 时写在字段后面的反引号字符串（如 `` `json:"name"` ``），格式是 `key:"value"` 空格分隔的多组键值对。编译时这些字符串被嵌入到类型信息中，运行时通过 `Type.Field(i).Tag.Get(key)` 读取。标准库 json、yaml 等包就是用这个机制把 struct 字段名映射到序列化后的字段名。

### Q6：反射比直接调用慢多少？

答：慢 1~2 个数量级（10~100 倍）。性能开销主要来自两方面：一是值要逃逸到堆上（因为 `interface{}` 参数）；二是每次反射操作都要做运行时类型检查而不是编译时确定。在热点路径上应避免反射，可用代码生成代替。

### Q7：反射的适用场景有哪些？什么时候不该用？

答：适用的场景：序列化库（json/yaml）、ORM、RPC 框架、依赖注入容器、配置绑定——这些工具需要处理任意类型，反射是唯一的办法。不该用的场景：普通业务逻辑、循环中的频繁操作、能用泛型、接口或代码生成替代的情况。

### Q8：`MethodByName` 找不到方法会怎样？

答：返回一个零值 `reflect.Value`，其 `IsValid()` 返回 false。之后调 `Call()` 会 panic。所以调用前最好检查 `method.IsValid()`。

---

## 一句话总结

反射让 Go 程序在运行时拥有了审视和修改自身的能力，以 `TypeOf`/`ValueOf` 为入口，通过 Type 看类型、通过 Value 读值改值、通过 Tag 读元数据；但它的代价是性能下降、编译期安全丢失和更复杂的代码——只在框架层面使用，业务代码优先选择泛型和接口。
