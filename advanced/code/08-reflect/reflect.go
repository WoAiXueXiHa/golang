package main

import (
	"fmt"
	"reflect"
)

// =============================================================
// 反射：程序在运行时检测自身类型和值、甚至修改自身的一种能力
// =============================================================
// 两个核心入口：
//   reflect.TypeOf(x)  → 返回 reflect.Type（接口），描述 x 的类型信息
//   reflect.ValueOf(x) → 返回 reflect.Value（结构体），承载 x 的值信息

// ---------------------------
// Demo 复用类型定义
// ---------------------------

// Stu 用于 struct 反射相关的 demo
type Stu struct {
	Name  string
	Age   int
	Score float64
}

// WrapInt 用于 Type vs Kind 的 demo（自定义类型，底层是 int）
type WrapInt int

// Calculator 用于方法反射的 demo
type Calculator struct{}

func (c Calculator) Add(a, b int) int {
	return a + b
}

func (c Calculator) Mul(a, b int) int {
	return a * b
}

// Config 用于结构体标签的 demo
type Config struct {
	Name string `json:"name" yaml:"name"`
	Port int    `json:"port" yaml:"port"`
}

// Add 用于函数反射的 demo
func Add(num1, num2 int) (int, error) {
	return num1 + num2, nil
}

// =============================================================
// 1. reflect.TypeOf() — 获取反射类型对象（reflect.Type 接口）
// =============================================================
// TypeOf 接收任意值，返回它的运行时类型描述。
// 返回的 reflect.Type 是一个接口，背后有具体实现类型承载信息。

func demo1_TypeOf() {
	fmt.Println("=== demo1: reflect.TypeOf() ===")

	// 基础类型
	var num int64 = 100
	t1 := reflect.TypeOf(num)
	fmt.Printf("int64 的 Type: %s\n", t1.String())

	// struct 类型
	st := Stu{Name: "zhangsan", Age: 18, Score: 95.5}
	t2 := reflect.TypeOf(st)
	fmt.Printf("Stu 的 Type: %s\n", t2.String())
}

// =============================================================
// 2. reflect.ValueOf() — 获取反射值对象（reflect.Value 结构体）
// =============================================================
// ValueOf 返回一个 reflect.Value，它不仅包含值，还可以通过它反向拿到 Type。
// TypeOf 只有类型信息；ValueOf 有类型信息 + 具体值。

func demo2_ValueOf() {
	fmt.Println("=== demo2: reflect.ValueOf() ===")

	var num int64 = 100
	v1 := reflect.ValueOf(num)
	fmt.Printf("int64 的 Value: %v\n", v1)

	st := Stu{Name: "zhangsan", Age: 18, Score: 95.5}
	v2 := reflect.ValueOf(st)
	fmt.Printf("Stu 的 Value: %v\n", v2)

	// 从 Value 也能获取 Type
	fmt.Printf("v2 的 Type: %s\n", v2.Type().String())
}

// =============================================================
// 3. Type vs Kind — 类型名 vs 底层分类
// =============================================================
// Type.Name()   → 返回 Go 语言层面的类型名（如 "int"、"Stu"、"WrapInt"）
// Type.Kind()   → 返回 reflect.Kind，表示底层是哪一类（如 int64、struct）
// 自定义类型 Name 是自己的名字，但 Kind 是底层基础类型的 Kind。

func demo3_TypeKind() {
	fmt.Println("=== demo3: Type vs Kind ===")

	var num1 int = 100
	var num2 WrapInt = 10004

	t1 := reflect.TypeOf(num1)
	fmt.Printf("num1 Type.Name()=%s, Kind=%v\n", t1.Name(), t1.Kind())

	t2 := reflect.TypeOf(num2)
	fmt.Printf("num2 Type.Name()=%s, Kind=%v\n", t2.Name(), t2.Kind())
	// num1 和 num2 的 Kind 相同（都是 int），但 Name 不同（int vs WrapInt）
}

// =============================================================
// 4. 获取 struct 反射值 — 遍历字段、读取值
// =============================================================
// Value.NumField()   → 字段数量
// Value.Field(i)     → 第 i 个字段的 reflect.Value
// 取值方法：.String()、.Int()、.Float()、.Interface() 等
// 注意：用错取值方法会 panic（如对 string 字段调 .Int()），
// 安全的做法是用 .Interface() 获取 interface{} 再 %v 打印。

func demo4_StructValue() {
	fmt.Println("=== demo4: struct 反射值 ===")

	st := Stu{Name: "zhangsan", Age: 18, Score: 95.5}
	v := reflect.ValueOf(st)

	fmt.Printf("字段数量: %d\n", v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		// 用 Interface() 安全获取值，避免类型不匹配 panic
		fmt.Printf("  字段%d: 类型=%v, 值=%v\n", i+1, field.Type().Name(), field.Interface())
	}
}

// =============================================================
// 5. 获取 map 反射值 — 遍历 key/value
// =============================================================
// Value.MapKeys()     → 返回所有 key 的 []reflect.Value
// Value.MapIndex(key) → 根据 key 获取对应的 value 的 reflect.Value

func demo5_MapValue() {
	fmt.Println("=== demo5: map 反射值 ===")

	m := map[int]uint32{1: 100, 2: 200, 3: 300}

	v := reflect.ValueOf(m)
	for _, k := range v.MapKeys() {
		val := v.MapIndex(k)
		fmt.Printf("  key=%d (类型=%s), value=%d (类型=%s)\n",
			k.Int(), k.Type().Name(), val.Uint(), val.Type().Name())
	}
}

// =============================================================
// 6. 获取 slice / array 反射值 — 下标遍历
// =============================================================
// Value.Len()    → 长度
// Value.Index(i) → 第 i 个元素的 reflect.Value
// 注意：slice 和 array 都用同一个 Index/Len 方法

func demo6_SliceValue() {
	fmt.Println("=== demo6: slice / array 反射值 ===")

	// slice
	slice := []int{1, 2, 3}
	v1 := reflect.ValueOf(slice)
	fmt.Print("slice: ")
	for i := 0; i < v1.Len(); i++ {
		fmt.Printf("%v ", v1.Index(i).Interface())
	}
	fmt.Println()

	// array
	arr := [3]int{4, 5, 6}
	v2 := reflect.ValueOf(arr)
	fmt.Print("array: ")
	for i := 0; i < v2.Len(); i++ {
		fmt.Printf("%v ", v2.Index(i).Interface())
	}
	fmt.Println()
}

// =============================================================
// 7. struct 反射类型 — 获取类型级别的信息（不涉及具体值）
// =============================================================
// Type.Name()       → 类型名
// Type.Kind()       → 底层分类
// Type.NumField()   → 字段数
// Type.Field(i)     → 第 i 个字段的 reflect.StructField（含 Name、Type 等）

func demo7_StructType() {
	fmt.Println("=== demo7: struct 反射类型 ===")

	st := Stu{Name: "zhangsan", Age: 18, Score: 95.5}
	t := reflect.TypeOf(st)

	fmt.Printf("类型名: %s\n", t.Name())
	fmt.Printf("Kind:   %v\n", t.Kind())
	fmt.Printf("字段数: %d\n", t.NumField())

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("  字段%d: name=%s, type=%s\n", i+1, f.Name, f.Type.String())
	}
}

// =============================================================
// 8. 指针反射类型 — 必须 Elem() 才能拿到指针指向的类型
// =============================================================
// 对指针调用 TypeOf，Kind 是 ptr，不能直接拿字段信息。
// 必须通过 .Elem() 解引用，拿到指向的具体类型后才能访问字段。

func demo8_PtrType() {
	fmt.Println("=== demo8: 指针反射类型 ===")

	st := &Stu{Name: "zhangsan", Age: 18, Score: 95.5}
	t := reflect.TypeOf(st)

	fmt.Printf("指针的 Kind: %v\n", t.Kind()) // ptr
	// 接下来必须用 Elem() 才能真正拿到 Stu 的类型信息
	fmt.Printf("解引用后类型名: %s\n", t.Elem().Name())
	fmt.Printf("解引用后 Kind:   %v\n", t.Elem().Kind())
	fmt.Printf("解引用后字段数: %d\n", t.Elem().NumField())

	for i := 0; i < t.Elem().NumField(); i++ {
		f := t.Elem().Field(i)
		fmt.Printf("  字段%d: name=%s, type=%s\n", i+1, f.Name, f.Type.String())
	}
}

// =============================================================
// 9. 函数反射类型 — 查看参数和返回值签名
// =============================================================
// Type.NumIn()   → 入参个数
// Type.In(i)     → 第 i 个入参的 Type
// Type.NumOut()  → 返回值个数
// Type.Out(i)    → 第 i 个返回值的 Type

func demo9_FuncType() {
	fmt.Println("=== demo9: 函数反射类型 ===")

	t := reflect.TypeOf(Add)

	fmt.Printf("入参个数: %d\n", t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		fmt.Printf("  入参%d: %s\n", i+1, t.In(i).Name())
	}

	fmt.Printf("返回值个数: %d\n", t.NumOut())
	for i := 0; i < t.NumOut(); i++ {
		fmt.Printf("  返回值%d: %s\n", i+1, t.Out(i).Name())
	}
}

// =============================================================
// 10. 获取 struct 的方法 — 反射查看类型上有哪些方法
// =============================================================
// Type.NumMethod()  → 方法数量（注意：只包含导出的方法）
// Type.Method(i)    → 第 i 个方法的 reflect.Method（含 Name、Type 等）

func demo10_StructMethods() {
	fmt.Println("=== demo10: 获取 struct 方法 ===")

	c := Calculator{}
	t := reflect.TypeOf(c)

	fmt.Printf("方法数: %d\n", t.NumMethod())

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fmt.Printf("  方法%d: name=%s, type=%s\n", i+1, m.Name, m.Type)
		// m.Type 是 func 类型，包含接收者作为第一个参数
	}
}

// =============================================================
// 11. 通过反射调用方法 — 动态选择方法并传参执行
// =============================================================
// Value.MethodByName(name)  → 获取指定名称的方法的 reflect.Value
// Value.Method(i)           → 按索引获取方法
// Value.Call(args)          → 调用方法，args 是 []reflect.Value，返回 []reflect.Value
// 注意：参数和返回值都要包装成 reflect.Value

func demo11_CallMethod() {
	fmt.Println("=== demo11: 反射调用方法 ===")

	c := Calculator{}
	v := reflect.ValueOf(c)

	// 方式1：按方法名调用
	method := v.MethodByName("Add")
	// 参数必须包装成 []reflect.Value
	args := []reflect.Value{
		reflect.ValueOf(3),
		reflect.ValueOf(4),
	}
	result := method.Call(args)
	// 返回值也是 []reflect.Value，用 Interface() 提取
	fmt.Printf("Add(3, 4) = %v\n", result[0].Interface())

	// 方式2：按方法名调用另一个方法
	method2 := v.MethodByName("Mul")
	result2 := method2.Call([]reflect.Value{
		reflect.ValueOf(5),
		reflect.ValueOf(6),
	})
	fmt.Printf("Mul(5, 6) = %v\n", result2[0].Interface())
}

// =============================================================
// 12. 通过反射修改值 — CanSet + Elem + SetXxx
// =============================================================
// 反射修改值有三个条件：
//   1. 必须传入指针（否则拿到的是副本，改不了原值）
//   2. 必须 .Elem() 解引用（拿到指针指向的可修改对象）
//   3. 调用对应的 Set 方法（SetInt、SetString、Set 等）
// CanSet() 返回 false = 不可修改；返回 true = 可以修改

func demo12_SetValue() {
	fmt.Println("=== demo12: 反射修改值 ===")

	// ---------- 错误示范：直接传值，CanSet = false ----------
	x := 100
	v1 := reflect.ValueOf(x)
	fmt.Printf("ValueOf(x) CanSet: %v （值不可改，因为是副本）\n", v1.CanSet())
	// v1.SetInt(200) ← 会 panic！

	// ---------- 错误示范：只传指针但不用 Elem ----------
	v2 := reflect.ValueOf(&x)
	fmt.Printf("ValueOf(&x) CanSet: %v （指针本身不可改）\n", v2.CanSet())
	// v2.SetInt(200) ← 也会 panic！

	// ---------- 正确做法：传指针 + Elem() ----------
	v3 := reflect.ValueOf(&x).Elem()
	fmt.Printf("ValueOf(&x).Elem() CanSet: %v （可以改了）\n", v3.CanSet())

	v3.SetInt(200)
	fmt.Printf("x 被修改为: %d\n", x)

	// ---------- 修改 struct 字段 ----------
	st := Stu{Name: "zhangsan", Age: 18, Score: 95.5}
	fmt.Printf("修改前: %+v\n", st)

	v4 := reflect.ValueOf(&st).Elem()
	// 通过 Field 定位 + SetString
	v4.Field(0).SetString("lisi")
	v4.Field(1).SetInt(20)

	fmt.Printf("修改后: %+v\n", st)
}

// =============================================================
// 13. 结构体标签 (StructTag) — 编译期关联到字段上的 key:value 元数据
// =============================================================
// Type.Field(i).Tag             → 返回 reflect.StructTag（本质是 string）
// Tag.Get(key)                  → 获取 key 对应的 value，找不到返回 ""
// Tag.Lookup(key)               → 获取 key 对应的 value + 是否存在（bool）
// 标签格式：`key1:"value1" key2:"value2"`
// 标准库 json、yaml、xml 等都用 struct tag 做字段映射

func demo13_StructTag() {
	fmt.Println("=== demo13: 结构体标签 ===")

	cfg := Config{Name: "server", Port: 8080}
	t := reflect.TypeOf(cfg)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// 打印完整 tag 字符串
		fmt.Printf("字段 %s 的完整 tag: %s\n", f.Name, string(f.Tag))

		// Get：拿不到返回空字符串
		jsonVal := f.Tag.Get("json")
		fmt.Printf("  json: %s\n", jsonVal)

		// Lookup：返回 (value, ok)
		yamlVal, ok := f.Tag.Lookup("yaml")
		if ok {
			fmt.Printf("  yaml: %s\n", yamlVal)
		}

		// 查询不存在的 key
		xmlVal, ok := f.Tag.Lookup("xml")
		fmt.Printf("  xml: %q (存在=%v)\n", xmlVal, ok)

		fmt.Println()
	}
}

// =============================================================
// 反射的优缺点总结
// =============================================================
// 优点：
//   1. 动态分派：可以根据字符串名字调用方法（如 json.Marshal 内部用反射取字段）
//   2. 通用工具：ORM、序列化库、依赖注入框架等都依赖反射
//   3. 减少重复代码：一套反射逻辑处理任意 struct
//
// 缺点：
//   1. 性能开销：反射操作涉及堆逃逸、运行时类型检查，比直接访问慢 1~2 个数量级
//   2. 编译期类型安全丢失：用错字段名或类型不会在编译时报错，只会在运行时 panic
//   3. 代码可读性下降：反射代码比直接访问更冗长、更难维护
//   4. 容易写出 panic 代码：CanSet 没检查就 Set、类型不匹配的取值方法都会 panic
//
// 原则：能用泛型/接口/代码生成解决的，尽量不用反射。

// =============================================================
// main：按需取消注释运行对应 demo
// =============================================================

func main() {
	// demo1_TypeOf()
	// demo2_ValueOf()
	// demo3_TypeKind()
	// demo4_StructValue()
	// demo5_MapValue()
	// demo6_SliceValue()
	// demo7_StructType()
	// demo8_PtrType()
	// demo9_FuncType()
	// demo10_StructMethods()
	// demo11_CallMethod()
	// demo12_SetValue()
	// demo13_StructTag()
}
