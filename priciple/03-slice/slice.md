# 什么是切片

切片是建立在数组之上的抽象类型

## 这一章要记住什么

- 切片 header 是值拷贝，底层数组是共享的——传到函数里改元素外部看得见，append 扩容后就跟外部没关系了
- `append` 返回的不一定和原切片共享底层数组，取决于 cap 有没有余量，所以永远要接收返回值
- 截取切片 `s[low:high]` 不复制数据，只新建一个 header，cap 会跟着缩（`newCap = oldCap - low`）
- 删除元素惯用 `append(s[:i], s[i+1:]...)`，但会污染共享底层数组的原切片
- 想完全独立的副本，只能是 `make` + `copy`，截取做不到
- **三索引切片 `s[low:high:max]` 可以限制 cap**，防止子切片 append 时意外覆盖父切片数据

---

## 数组

Go 语言中数组是一个值，数组变量表示了整个数组，和 C/C++ 不同（指向数组首元素的指针）

利用代码看一下：

```go
package main

import "fmt"

// 将数组传递到函数中，数组的地址不一样
func test(arr [3]int) {
    fmt.Printf("arr 内: %p\n", &arr)
}
func f1() {
    arr := [3]int{1, 2, 3}
    test(arr)
    fmt.Printf("arr 外: %p\n", &arr)
}

// 拷贝数组，修改旧数组，对新数组无影响
func f2() {
    arr1 := [3]int{1, 2, 3}
    arr2 := arr1

    arr1[0] = 100
    fmt.Println(arr1)
    fmt.Println(arr2)
}

func main() {
    f1()
    fmt.Printf("---------------------------\n")
    f2()
}

```

输出：

```bash
[vect@ubuntu-dev ~/golang/priciple/03-slice/demo1]$ go run demo1.go 
arr 内: 0xc0000a0018
arr 外: 0xc0000a0000
---------------------------
[100 2 3]
[1 2 3]
```

可以发现：

- 数组传值到函数中，数组的地址不一样
- 拷贝数组，修改旧的数组，对新的数组无影响

Go 的数组类似 C++ 的 array，定长数组，长度是固定的

而 slice 就类似 C++ 的 vector，变长数组，动态扩容



# 切片底层原理剖析

## 切片结构

切片底层就是一个结构体：

```golang
type slice struct {
    // 指向一块连续内存空间的起始位置
    array unsafe.Pointer
    len int
    cap int
}
```

## 切片扩容机制

### 1.计算目标容量（预估阶段）

当 slice 触发 `append` 导致超出当前容量时，Go 会通过以下两个阶段来决定最终的容量。首先是计算预估容量：

* **case1**：新切片长度 > 旧切片容量的两倍，则预估容量直接定为**新切片长度**。
* **case2**：若不满足 case1，则根据旧切片容量进行平滑过渡：
  1. **旧切片容量 < 256**：新切片的预估容量直接翻倍，即 `newcap = 2 * oldcap`。
  2. **旧切片容量 $\ge$ 256**：每次扩大为原来的 1.25 倍，并且每次为了平滑过渡，还会固定加上 $\frac{3}{4} \times 256$（即 `+192`），直到预估容量 $\ge$ 新切片长度。

```golang
// newcap = newcap * 1.25 + 192 的底层高效位移写法
newcap += (newcap + 3*threshold) / 4
```

### 2. 内存对齐（最终容量确定）

> **关键结论**：计算出预估容量后，**最终容量并不一定等于预估容量，而是由底层内存分配器决定。**

Go 运行时会调用 `roundupsize` 函数，将预估容量占用的内存大小（预估容量 $\times$ 元素大小）向上对齐到与其最接近的底层内存规格（Size Class）

* **例如**：若预估容量计算出需要 300 字节，而 Go 内存分配器现有的固定分配规格中没有 300 字节，只有 320 字节，则系统会直接分配 320 字节。此时反向推导出的最终 `cap` 就会比预估值稍大一些。
* **不严谨之处在于**：这种机制虽然会导致容量多出预期的几个元素，但在宏观上极大减少了堆内存碎片，并加速了内存分配效率。



## 和 C++ vector 进行对比

先说结论：

**扩容策略差异：** vector 追求 **确定性的几何增长（1.5倍或2倍）**，slice 计算出预估容量后，**强行接入内存对齐**，最终容量由底层内存分配器决定

用个表格总结：

| **维度**           | **Go slice**                                               | **C++ std::vector**                   |
| ------------------ | ---------------------------------------------------------- | ------------------------------------- |
| **扩容系数**       | $<256$ 元素时 2 倍；$\ge 256$ 时过渡到 1.25 倍 + 192       | GCC/Clang 固定 2 倍；MSVC 固定 1.5 倍 |
| **最终容量确定性** | **不确定**。受限于底层内存分配器的 Size Class 规格向上取整 | **确定**                              |



## 内存对齐规则对比

内存对齐本质是**用空间换时间，确保CPU能通过一次总线周期高效读取数据**

### C++ 内存对齐

前提：编译器都有默认的**对齐数**，64位默认为8

规则一：**成员自身对齐（决定字段偏移量）**

$min(自身类型大小，默认对齐数)$的整数倍

规则二：**结构体整体对齐（决定结构体最终大小）**

结构体大小=$min(内部最大基础成员的大小，默认对齐数)$的整数倍

例如：64位系统

```cpp
struct A {
    char c;
    int b;
    double c;
}
c _ _ _ b b b b c c c c c c c c
0 1 2 3 4 5 6 7 8 .......    15
```

最终大小为16字节



### Go 内存对齐

和 C++ 完全一致，但多了针对垃圾回收机制的特殊**尾部边界处理**：

**零大小尾部阻隔：**

若一个结构体的**最后一个字段**的大小是0（例如空结构体`struct{}`），且该结构体还会被其他对象引用或者作为数组元素，Go 编译器会在尾部**强制填充1字节**并进行对齐

设计原因：

若不填充，指向该空结构体的指针就会直接指向结构体外部的下一个对象，误认为下一个对象还在被引用，导致内存泄漏

还是64位系统：

```golang
type BadLayout struct {
    a int32      // 4 字节
    b struct{}   // 0 字节，但处于尾部。为了 GC 安全，强行填充并对齐至 4 字节
} // 总大小 = 4 + 4 = 8 字节

type GoodLayout struct {
    b struct{}   // 0 字节，处于头部
    a int32      // 4 字节
} // 总大小 = 0 + 4 = 4 字节 (无需尾部填充)
```



# 切片行为分析

## 1. 切片传参的本质：值传递 + 共享底层数组

```golang
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func PrintSlice(s *[]int) {
	// Go 是强类型，正常情况 *[]int 绝对不能转成 *reflect.SliceHeader
	// 而 unsafe.Pointer 类似 void*，可以接收任意类型指针
	// 对于 reflect.SliceHeader
	// type SliceHeader struct {
	// 	Data uintptr  // 对应底层数组地址，这个不是指针，就是存了地址数字的类型而已
	// 	Len  int      // 对应长度
	// 	Cap  int      // 对应容量
	// }

	ss := (*reflect.SliceHeader)(unsafe.Pointer(s))

	fmt.Printf("slice struct: %+v, slice is %v\n", ss, s)
}

func test(s []int) {
	PrintSlice((&s))
}

func demo2_slice_func() {
	s := make([]int, 5, 10)
	PrintSlice(&s)
	test(s)
}

func main() {
	demo2_slice_func()
}

```

输出：

```bash
slice struct: &{Data:824633884752 Len:5 Cap:10}, slice is &[0 0 0 0 0]
slice struct: &{Data:824633884752 Len:5 Cap:10}, slice is &[0 0 0 0 0]
```

先看代码做了什么：用 `make([]int, 5, 10)` 创建一个 len=5、cap=10 的切片，在 main 里打印一次 header，传到 `test` 函数里再打印一次。

两次输出的 **Data 地址完全一样**。

这说明什么？Go 所有函数参数都是值传递，切片也不例外——传进去的是 slice header 的一份**拷贝**。但 header 里的 Data 字段是个指针（准确说是 uintptr，存的是地址值），拷贝后仍然指向**同一块底层数组**。

```text
main 里的 s:                     test 里的 s (拷贝):
+------------------+           +------------------+
| Data: 0x...4752  |---┐       | Data: 0x...4752  |---┐
| Len:  5          |   |       | Len:  5          |   |
| Cap:  10         |   |       | Cap:  10         |   |
+------------------+   |       +------------------+   |
                       |                              |
                       v  同一块底层数组                v
              +-----------------------------------+
              | [0] [1] [2] [3] [4] (预留 5 个空位) |
              +-----------------------------------+
               len = 5                cap = 10
```

两个 header 是独立的（Len/Cap 互不影响），但它们看到的底层数组是同一片内存。

总结一下：

> 切片传到函数里，切片 header 是复制品，底层数组是共享的

## 2. 修改切片：下标修改 vs append

```golang
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func PrintSlice(s *[]int) {
	ss := (*reflect.SliceHeader)(unsafe.Pointer(s))

	fmt.Printf("slice struct: %+v, slice is %v\n", ss, s)
}

func test(s []int) {
	PrintSlice((&s))
}

// 底层数组不变
func demo3_case1(s []int) {
	s[1] = 1000
	PrintSlice(&s)
}

// 底层数组变化
func demo3_case2(s []int) {
	s = append(s, 1000)
	s[1] = 1000
	PrintSlice(&s)
}
func demo3_infunc_modify() {
	s := make([]int, 5)
	demo3_case1(s)
	demo3_case2(s)
	PrintSlice(&s)
}

func main() {
	demo3_infunc_modify()
}

```

输出：

```bash
slice struct: &{Data:824633811472 Len:5 Cap:5}, slice is &[0 1000 0 0 0]
slice struct: &{Data:824633884832 Len:6 Cap:10}, slice is &[0 1000 0 0 0 1000]
slice struct: &{Data:824633811472 Len:5 Cap:5}, slice is &[0 1000 0 0 0]
```

对于通过索引修改：

调用前 `s := make([]int, 5)`，底层数组全是 0。`s[1] = 1000` 直接修改了底层数组的第 1 号位置。外部切片和它共享同一个底层数组，所以外部看到的也是 `[0 1000 0 0 0]`。

对于先 append 再通过索引修改

```go
func demo3_case2(s []int) {
    s = append(s, 1000)   // len==cap==5，append 触发扩容！
    s[1] = 1000           // 这次改的是新数组
}
```

关键在这里：`make([]int, 5)` 创建的切片 len=cap=5，**没有预留空间**。`append(s, 1000)` 发现 len(6) > cap(5)，必须扩容——于是分配一块**新**的底层数组，把旧元素拷过去，再追加 1000。

此时函数里的 `s` 已经指向新数组了，后续 `s[1] = 1000` 改的是新数组，**跟外部切片已经没关系了**。

```text
append 前（case1 执行后）:
外部 s:                         case2 内 s:
+------------------+            +------------------+
| Data: 0x...1472  |--+        | Data: 0x...1472  |--+
| Len:  5          |  |        | Len:  5          |  |
| Cap:  5          |  |        | Cap:  5          |  |
+------------------+  |        +------------------+  |
                      +-------> [0,1000,0,0,0]  <----+
                                (底层数组，cap=5，已满)

append(s, 1000) 之后:
外部 s:                         case2 内 s (重新赋值后):
+------------------+            +------------------+
| Data: 0x...1472  |--+        | Data: 0x...4752  |-----+
| Len:  5          |  |        | Len:  6          |     |
| Cap:  5          |  |        | Cap:  10         |     |
+------------------+  |        +------------------+     |
                      |                                 |
                      v                                 v
              [0,1000,0,0,0]              [0,1000,0,0,0,1000]
              (旧数组，外部还指着它)         (新数组，函数里的 s 指着它)
                                          s[1]=1000 改这里
```

输出验证了这一点：

- case1 内 Data 和外部一样（`0x...1472`）
- case2 内 Data 变了（`0x...4752`），是新数组的地址
- 外部 Data 仍然是旧地址 `0x...1472`，且值还是 case1 留下的 `[0 1000 0 0 0]`

总结一下：

> 通过下标改元素，影响的是共享的底层数组，外部可见。通过 append 追加导致扩容时，函数内部的 s 指向新数组，后续操作跟外部完全脱钩。**能不能影响外部，取决于 append 是否触发扩容。**

## 3. 截取切片：新建视图，不复制数据

```golang
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func PrintSlice(s *[]int) {
	ss := (*reflect.SliceHeader)(unsafe.Pointer(s))

	fmt.Printf("slice struct: %+v, slice is %v\n", ss, s)
}

func case1(s []int) {
	s = s[1:]
	PrintSlice(&s)
}

func case2(s []int) {
	s = s[1:3]
	PrintSlice(&s)
}

func case3(s []int) {
	s = s[len(s)-1:]
	PrintSlice(&s)
}

func case4(s []int) {
	s1 := s[2:]
	PrintSlice(&s1)
}

func main() {
	s := make([]int, 5)

	case1(s)
	case2(s)
	case3(s)
	case4(s)

	PrintSlice(&s)
}

```

输出：

```bash
slice struct: &{Data:824633811480 Len:4 Cap:4}, slice is &[0 0 0 0]
slice struct: &{Data:824633811480 Len:2 Cap:4}, slice is &[0 0]
slice struct: &{Data:824633811504 Len:1 Cap:1}, slice is &[0]
slice struct: &{Data:824633811488 Len:3 Cap:3}, slice is &[0 0 0]
slice struct: &{Data:824633811472 Len:5 Cap:5}, slice is &[0 0 0 0 0]
```

原始切片 `s := make([]int, 5)` 生成 5 个零值 int，Data 从 `0x...1472` 开始，每个 int 占 8 字节。

```text
底层数组 (每个格子 8 字节):
地址:   0x1472  0x147A  0x1482  0x148A  0x1492  (十六进制，差 8)
       +-------+-------+-------+-------+-------+
       |  [0]  |  [1]  |  [2]  |  [3]  |  [4]  |
       +-------+-------+-------+-------+-------+
         ^                                       ^
         |                                       |
      原始 s.Data                          原始 s 能看到的最远位置
      (0x...1472)                          (0x...1472 + 5*8)
```

四个 case 分别做了不同截取，看输出数据来推理规律：

| 操作     | Data 地址   | 相比原 Data 偏移   | Len  | Cap  |
| -------- | ----------- | ------------------ | ---- | ---- |
| 原始 `s` | `0x...1472` | 0                  | 5    | 5    |
| `s[1:]`  | `0x...1480` | +8（跳 1 个 int）  | 4    | 4    |
| `s[1:3]` | `0x...1480` | +8（跳 1 个 int）  | 2    | 4    |
| `s[4:]`  | `0x...1504` | +32（跳 4 个 int） | 1    | 1    |
| `s[2:]`  | `0x...1488` | +16（跳 2 个 int） | 3    | 3    |

规律非常清楚：**`s[i:j]` 不会分配新内存，只是构造了一个新的 slice header。**

```text
s[1:]  = s[1:5] → Data = 原Data + 1*8, Len = 4, Cap = 原Cap - 1
s[1:3]          → Data = 原Data + 1*8, Len = 2, Cap = 原Cap - 1
s[4:]  = s[4:5] → Data = 原Data + 4*8, Len = 1, Cap = 原Cap - 1
s[2:]  = s[2:5] → Data = 原Data + 2*8, Len = 3, Cap = 原Cap - 2
```

通用公式（`s[low:high]`）：

- `newData = oldData + low * sizeof(element)`
- `newLen = high - low`
- `newCap = oldCap - low`

注意 `s[1:3]` 虽然 Len 只有 2，但 Cap 还有 4——说明它"记得"自己底层数组从位置 1 往后还有 3 个元素的空间（只是当前不暴露）。这意味着如果对它做 append，只要不超出 cap，**仍然会在原底层数组上操作**。

原始 `s` 的 Data、Len、Cap 全程不变——四个 case 里截取出的都是新 header，赋值给了函数内的局部变量，不影响外部。

总结一下：

> 截取切片就是"换个角度看同一块内存"。Data 指针往后挪、Len 缩短、Cap 缩短，没有数据拷贝。所有截取出来的切片共享底层数组，一个改了元素，其他都看得见。

## 4. 删除元素：append 拼接的陷阱

```golang
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func PrintSlice(s *[]int) {
	ss := (*reflect.SliceHeader)(unsafe.Pointer(s))

	fmt.Printf("slice struct: %+v, slice is %v\n", ss, s)
}

func main() {
	s := []int{0, 1, 2, 3, 4}

	_ = s[4]
	PrintSlice(&s)
	// 删除第一个元素，从0开始计数
	// [0,1) + [2, len(s))
	s1 := append(s[:1], s[2:]...)
	{
		// 拷贝元素
		// 0, 1, 2, 3, 4
		// 0, 2, 3, 4, 4
	}

	PrintSlice(&s1)
	PrintSlice(&s)

	// 访问原切片
	_ = s[4]
	// 访问从原切片中删除了一个元素的切片
	_ = s1[4]

}

```

输出：

```bash
slice struct: &{Data:824634392576 Len:5 Cap:5}, slice is &[0 1 2 3 4]
slice struct: &{Data:824634392576 Len:4 Cap:5}, slice is &[0 2 3 4]
slice struct: &{Data:824634392576 Len:5 Cap:5}, slice is &[0 2 3 4 4]
panic: runtime error: index out of range [4] with length 4

goroutine 1 [running]:
main.main()
        /home/vect/golang/priciple/03-slice/demo4/demo4_delete.go:35 +0x1b6
exit status 2
```



代码要做的事：从 `[0, 1, 2, 3, 4]` 里删除索引 1 的元素（值 `1`），得到 `[0, 2, 3, 4]`。

Go 没有内置的删除切片元素的方法，惯用写法是 `s = append(s[:i], s[i+1:]...)`。

拆解这个过程：

```text
原始 s:  [0, 1, 2, 3, 4]  len=5  cap=5
底层:    +---+---+---+---+---+
         | 0 | 1 | 2 | 3 | 4 |
         +---+---+---+---+---+
         ^               ^
    s.Data            s.Data+4*8

s[:1]:  取前 1 个元素 [0]
         Data 同 s, Len=1, Cap=5  ← 注意 cap 还是 5，有 4 个空位

s[2:]:  从索引 2 开始取到底 [2,3,4]
         Data = s.Data + 2*8, Len=3, Cap=3

append(s[:1], s[2:]...):
         s[:1] 还剩 4 个 cap 空位，能装下 3 个元素，不扩容！
         在底层数组的位置 1、2、3 依次写入 2、3、4
```

执行后底层数组的变化：

```text
操作前:   +---+---+---+---+---+
          | 0 | 1 | 2 | 3 | 4 |
          +---+---+---+---+---+

写入后:   +---+---+---+---+---+
          | 0 | 2 | 3 | 4 | 4 |   ← 位置 4 的旧值 4 没被覆盖
          +---+---+---+---+---+
            ^~~~~~~~~~~~^
            s1 看到的范围  len=4
            ^~~~~~~~~~~~~~~~^
            s  看到的范围    len=5
```

这就解释了输出：

- `s1` = `[0 2 3 4]`，len=4，cap=5
- `s` = `[0 2 3 4 4]`，len=5——原切片底层数组被 append 就地修改了，尾部多了一个 `4`
- `s[4]` 访问成功（s 的 len=5），但 `s1[4]` **panic**（s1 的 len=4）

两个关键问题被暴露出来：

1. **原切片被"污染"了**。`append` 没有扩容，直接在共享的底层数组上写，`s` 看到的内容跟着变了。
2. **新切片 len 变小了**。`s1` 是"逻辑上删除了一个元素"的切片，它的 len=4，访问 `s1[4]` 直接越界 panic——即使底层数组那个位置确实有值。

总结一下：

> 用 append 拼接来删除元素，本质是把后面的元素往前拷贝，覆盖掉要删的那个位置。如果原切片 cap 够大，不会触发扩容，操作就在原底层数组上发生——原切片的内容也会被连带修改。删除后新切片的 len 少 1，按 len 访问才是安全的，不要以为底层数组还有值就能越界访问。

## 5. append 的行为：有容量走原地，没容量走扩容

```golang
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func PrintSlice(s *[]int) {
	ss := (*reflect.SliceHeader)(unsafe.Pointer(s))

	fmt.Printf("slice struct: %+v, slice is %v\n", ss, s)
}

func case1() {
	s1 := make([]int, 3, 3)
	s1 = append(s1, 1)

	PrintSlice(&s1)
}

func case2() {
	s1 := make([]int, 3, 4)
	s2 := append(s1, 1)

	PrintSlice(&s1)
	PrintSlice(&s2)
}

func case3() {
	s1 := make([]int, 3, 3)
	s2 := append(s1, 1)

	PrintSlice(&s1)
	PrintSlice(&s2)
}

func main() {
	case1()
	case2()
	case3()
}

```

```bash
slice struct: &{Data:824633811472 Len:4 Cap:6}, slice is &[0 0 0 1]
slice struct: &{Data:824633827584 Len:3 Cap:4}, slice is &[0 0 0]
slice struct: &{Data:824633827584 Len:4 Cap:4}, slice is &[0 0 0 1]
slice struct: &{Data:824633819424 Len:3 Cap:3}, slice is &[0 0 0]
slice struct: &{Data:824633811520 Len:4 Cap:6}, slice is &[0 0 0 1]
```

三个 case 对比了 append 的两种情况。

| Case  | 原始切片       | append 后      | 是否扩容                    | Data 是否变 |
| ----- | -------------- | -------------- | --------------------------- | ----------- |
| case1 | `len=3, cap=3` | `len=4, cap=6` | **是**（3 < 256，扩容到 6） | 变          |
| case2 | `len=3, cap=4` | `len=4, cap=4` | **否**（4 ≤ 4，原地追加）   | 不变        |
| case3 | `len=3, cap=3` | `len=4, cap=6` | **是**（3 < 256，扩容到 6） | 变          |

case2 最值得关注：

```text
s1 := make([]int, 3, 4)   // len=3, cap=4, 还有一个预留空位
s2 := append(s1, 1)       // 不扩容，直接在预留空位写 1

底层数组:
         +---+---+---+---+
         | 0 | 0 | 0 |   |   ← s1 创建后的状态 (cap=4, len=3 只暴露前 3 个)
         +---+---+---+---+
           s1 可见 ↑
         
         +---+---+---+---+
         | 0 | 0 | 0 | 1 |   ← append 后 (位置 3 填入了 1)
         +---+---+---+---+
           s1 可见 ↑   ↑ s2 可见
```

`s1` 和 `s2` 的 Data 地址相同（都是 `0x...0800`），因为 **append 没有扩容，直接用了预留空间**。`s1` 的 len 是 3，看不到第 4 个元素；`s2` 的 len 是 4，能看到。

case3 和 case1 本质一样——cap 已满，append 必须扩容，`s2` 拿到新数组，跟 `s1` 彻底分家。

对比三组 Data 地址：

- case2（未扩容）：s1.Data == s2.Data → 共享底层
- case1 / case3（扩容）：s1.Data 和新的 s1/s2.Data 完全不同 → 独立底层

总结：

> `append` 返回的切片不一定和原切片共享底层数组。能不能共享，取决于原切片的 cap 是否还有余量。这就是为什么官方一直强调 `s = append(s, ...)`——你永远不知道 append 会不会换底层数组，不接收返回值就等于把新数据丢了。

## 6. 深拷贝：copy 才是真正的复制

```golang
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func PrintSlice(s *[]int) {
	ss := (*reflect.SliceHeader)(unsafe.Pointer(s))

	fmt.Printf("slice struct: %+v, slice is %v\n", ss, s)
}

func main() {
	s1 := []int{1, 2, 3}
	s2 := make([]int, len(s1))

	copy(s2, s1)

	PrintSlice(&s1)
	PrintSlice(&s2)

}

```

输出：

```golang
slice struct: &{Data:824634392576 Len:3 Cap:3}, slice is &[1 2 3]
slice struct: &{Data:824634392600 Len:3 Cap:3}, slice is &[1 2 3]
```

```go
s1 := []int{1, 2, 3}
s2 := make([]int, len(s1))   // 先分配一个等长的独立切片
copy(s2, s1)                  // 把 s1 的元素逐个拷贝到 s2 的底层数组
```

输出里 `s1.Data` 和 `s2.Data` 是不同的地址——两个切片各自拥有独立的底层数组，互不影响。

```text
s1:                          s2:
+------------------+         +------------------+
| Data: 0x...9328  |--+      | Data: 0x...9352  |--+
| Len:  3          |  |      | Len:  3          |  |
| Cap:  3          |  |      | Cap:  3          |  |
+------------------+  |      +------------------+  |
                      v                             v
              +---+---+---+                 +---+---+---+
              | 1 | 2 | 3 |                 | 1 | 2 | 3 |
              +---+---+---+                 +---+---+---+
              数组 A (独立)                   数组 B (独立)
```

`copy(dst, src)` 的行为要点：

- 拷贝的元素数量 = `min(len(dst), len(src))`
- 只拷贝元素值，不共享底层数组
- 如果 dst 比 src 短，src 多出来的元素不会拷过去；如果 dst 比 src 长，多出来的位置保持原值

对比：截取切片 `s2 := s1[:]` 仍然共享底层数组，改 `s2` 会影响 `s1`，这不是深拷贝。

总结：

> 要想得到一个和原切片**完全独立**的副本，必须先 `make` 分配新切片，再 `copy` 拷贝元素。截取只是创建了一个新视图，底层还是同一块内存。

## 7. 三索引切片：控制 cap 防止误伤

前面讲截取时说过，`s[low:high]` 的 cap 是 `oldCap - low`。这意味着子切片 append 时可能会一直写到父切片的数据区域：

```go
parent := make([]int, 5, 5)
child := parent[:3]       // len=3, cap=5 — 还有 2 个空位！
child = append(child, 1)  // 不扩容，直接在 parent 的底层数组位置 3 写 1
// parent 变成了 [0,0,0,1,0]，被"误伤"了
```

Go 提供了**三索引切片** `s[low:high:max]` 来限制子切片的 cap：

```go
child := parent[:3:3]     // len=3, cap=3 — cap 被限制住了
child = append(child, 1)  // 触发扩容，分配到新数组，parent 安全
```

公式：

- `newLen = high - low`
- `newCap = max - low`
- 约束：`0 ≤ low ≤ high ≤ max ≤ cap(s)`

```text
parent := make([]int, 5, 5)   // [0, 0, 0, 0, 0]  len=5 cap=5

parent[:3]    → Data=parent.Data, Len=3, Cap=5   ← 危险，append 可能覆盖
parent[:3:3]  → Data=parent.Data, Len=3, Cap=3   ← 安全，append 立即扩容
parent[:3:5]  → Data=parent.Data, Len=3, Cap=5   ← 显式允许共享 2 个空位
```

实际意义：**如果子切片只是读数据，不用三索引也行；但如果子切片后续可能 append，用三索引限制 cap 可以防止意外修改父切片。** 这在并发场景下尤其重要——你以为两个切片互不影响，其实 append 在共享底层上写了。

### 总结一下

> 三索引切片 `s[low:high:max]` 可以精确控制子切片的容量，防止 append 时意外污染父切片。`max` 决定了子切片能"看到"多远的底层数组空间。

## 8. nil slice vs 空 slice：面试必问的区分题

```go
var s1 []int           // nil slice
s2 := make([]int, 0)   // 空 slice，非 nil
s3 := []int{}          // 空 slice，非 nil
```

从 header 视角看：

```text
nil slice (var s []int):
+------------------+
| Data: 0 (nil)    |
| Len:  0          |
| Cap:  0          |
+------------------+

空 slice (make([]int, 0)):
+------------------+
| Data: 0x...      |  ← 有一个合法的地址（zerobase，Go 运行时的全局零值基地址）
| Len:  0          |
| Cap:  0          |
+------------------+
```

哪些地方会体现出区别？

| 场景           | nil slice            | 空 slice                  |
| -------------- | -------------------- | ------------------------- |
| `s == nil`     | `true`               | `false`                   |
| `len(s)`       | `0`                  | `0`                       |
| `cap(s)`       | `0`                  | `0`                       |
| `for range s`  | 0 次循环             | 0 次循环                  |
| `append`       | 可以，正常用         | 可以，正常用               |
| `json.Marshal` | `"null"`             | `"[]"`                    |
| 等值比较       | 只能跟 nil 比        | 只能跟 nil 比（切片不能互相 `==`） |

关键结论：**大部分操作对 nil slice 和空 slice 行为一致，但序列化结果是不同的。API 设计中通常用 nil slice 表示"无数据"，用空 slice 表示"数据为空列表"。**

### 总结一下

> nil slice 的 Data 字段是 0，空 slice 的 Data 指向合法的内存地址（zerobase）。日常使用 `append` / `len` / `range` 没区别，但 JSON 序列化一个返回 `null` 一个返回 `[]`。

## 9. for range 遍历切片：值拷贝的坑

```go
s := []int{1, 2, 3}
for _, v := range s {
    v = v * 2   // 没用！v 是元素的拷贝，改它不影响 s
}
fmt.Println(s)  // [1 2 3]
```

`for range` 遍历切片时，`v` 是元素的**值拷贝**，修改 `v` 不影响原切片。如果想原地修改，用索引：

```go
for i := range s {
    s[i] = s[i] * 2  // 通过索引直接改底层数组
}
fmt.Println(s)  // [2 4 6]
```

还有一个容易忽略的点：**取 `v` 的地址每次都一样**。因为 `v` 在整个循环过程中是同一个变量，只是每次被赋予不同的值：

```go
for _, v := range s {
    fmt.Printf("%p\n", &v)  // 三次输出的地址相同！
}
```

如果你想拿到每个元素真正的地址，应该用 `&s[i]`，不能用 `&v`。

### 总结一下

> `for range` 的迭代变量是元素的值拷贝，修改它不影响原切片。取 `&v` 的地址也是同一个地址，不是每个元素的真实地址。要改元素或取地址，用索引 `s[i]`。

## 10. 截取大切片与内存泄漏

前面的内容反复强调过：截取不复制数据，子切片和父切片共享底层数组。这带来一个问题——如果父切片很大，你只截取了一小段，但**整块底层数组都会被 GC 标记为"仍在使用"**：

```go
big := make([]byte, 1<<30)   // 1GB
small := big[100:200]         // 只用了 100 个字节
// 但 big 的底层 1GB 数组不会被 GC，因为 small 还引着它
```

即使 `big` 本身出了作用域不再使用，只要 `small` 还活着，底层那 1GB 就收不回来。

修复方式——把想要的那一小段**真正拷贝出来**：

```go
small := make([]byte, 100)
copy(small, big[100:200])  // 现在 small 有自己的 100 字节底层数组
// big 如果不再被引用，1GB 可以被 GC 了
```

### 总结一下

> 从大切片截取小切片，小切片仍然持有对大底层数组的引用，导致整块内存无法回收。解决方法是 `make` + `copy` 把数据真正拷贝出来，断开对大数组的引用。

---

## 易错点

1. **把切片传给函数后，在函数里 append 却指望外部看到**。如果触发了扩容，外部切片完全不受影响。要对外部切片做 append，应该传 `*[]int` 指针，或者把新切片 return 回去。

2. **截取出来的切片仍持有对原底层数组的引用**。比如从一个百万元素的大切片截一小段出来，虽然新切片的 len 很小，但 cap 可能很大（`cap = oldCap - low`），底层大数组不会被 GC。如果只想要那一小段，应该 `make` + `copy`。

3. **用 append 拼接方式删除元素后，原切片的内容也变了**。案例 4 中 `s` 从 `[0,1,2,3,4]` 变成 `[0,2,3,4,4]`——因为 append 在原底层数组上直接覆写了。这不是 bug，但要知道自己在干什么。

4. **把 nil slice 和空 slice 搞混**。`var s []int`（nil slice，Data 为 0，len=cap=0）和 `s := make([]int, 0)`（空 slice，Data 有地址指向 zerobase，len=cap=0）不一样。`json.Marshal` 对 nil slice 输出 `null`，对空 slice 输出 `[]`，这在 API 设计中是两回事。

5. **`copy` 的拷贝数量由较短的那个切片决定**。如果 `copy(dst, src)` 的 dst 比 src 短，src 多余的元素会被丢弃，不会自动扩容 dst。

6. **`for range` 的迭代变量 v 是值拷贝**，修改它不影响原切片。取 `&v` 拿到的永远是同一个临时变量的地址，不是元素的真实地址。要取元素地址用 `&s[i]`。

7. **子切片 append 可能误伤父切片**。截取出的子切片 cap 如果大于 len，append 时会在共享底层数组上直接写入，覆盖父切片的数据。想避免就用三索引切片 `s[low:high:max]` 限制子切片容量。

8. **切片不能直接用 `==` 比较**。两个切片（除了跟 nil 比较）不能直接用 `==` 判断是否相等，因为底层数组可能共享、元素可能含不可比较类型。如果需要比较，用 `reflect.DeepEqual` 或自己写循环。

---

## 快问快答

### Q1：切片是引用类型吗？

答：可以说是，也可以说不是——要看你问的是哪个层面。

从行为上看，切片传参时共享底层数组，改元素外部看得见，这很像"引用"。

但从底层实现看，切片的 header 是值拷贝传到函数里的（Data、Len、Cap 三个字段都是值）。所以 Go 官方说切片是值类型，只是 header 里包含了一个指向底层数组的指针。

面试时可以说：**Go 里没有"引用类型"这个说法。切片 header 是值，但 header 里有指向底层数组的指针，所以共享底层数组。**

### Q2：`s = append(s, ...)` 为什么一定要接收返回值？

答：因为 `append` 之后底层数组可能变了。如果原切片 cap 不够，append 会分配新数组、拷走数据、追加新元素，然后返回的切片指向新数组。不接收返回值的话，你手里的旧切片还是指向旧数组，新数据等于白追加了。

### Q3：`s = append(s, x)` 是不是线程安全的？

答：不是。Go 里任何操作都不是线程安全的，除非明确文档说明。多个 goroutine 并发 append 同一个切片会导致数据竞争，需要加锁或用 channel 串行化。

### Q4：`nil` 切片和空切片，什么时候需要区分对待？

答：大部分时候不用，`append`、`len`、`range` 行为都一样。但序列化 JSON 时，nil slice → `null`，空 slice → `[]`。写 REST API 时，如果前端期望数组字段始终是 `[]` 而不是 `null`，初始化时就要用 `make([]T, 0)` 而不是 `var s []T`。

### Q5：对一个 nil 切片取 `len` 或 `cap` 会 panic 吗？

答：不会。`len(nilSlice)` 返回 0，`cap(nilSlice)` 也返回 0。nil slice 底层 header 的三个字段都是零值（Data=0, Len=0, Cap=0），调用 `len`/`cap` 是安全的。只有对 nil 切片做 `s[0]` 这种索引操作才会 panic。

### Q6：截取切片 `s[1:3]` 后，新切片的 cap 是多少？

答：`cap(s) - 1 = cap(s) - low`。截取不会复制数据，新切片的 Data 指针往后挪了 `low` 个元素，能看到的底层数组范围自然少了 `low` 个位置。

### Q7：为什么 `copy` 不会自动扩容 dst？

答：Go 的设计哲学是让程序员明确控制内存分配。`copy` 只拷贝 `min(len(dst), len(src))` 个元素，不会偷偷分配新内存。如果想让 dst 刚好装下 src 的所有元素，先确保 `len(dst) >= len(src)`，或者在 `make` 时把 len 设对。

### Q8：扩容后旧的底层数组会怎样？

答：扩容时 `growslice` 分配了新数组，旧数组如果没有其他切片引用，就会被 GC 回收。但如果还有其他切片也指向旧数组（比如截取出来的子切片），旧数组就会一直活着——这也是截取大切片可能导致内存泄漏的原因。

### Q9：`var s []int`、`s := []int{}`、`s := make([]int, 0)` 哪种写法最好？

答：场景不同选择不同：

- `var s []int`：最推荐，零值初始化，不分配任何内存
- `s := make([]int, 0)`：需要明确语义——"这是一个已经初始化但为空的列表"，或者在 JSON API 中需要区分 `null` 和 `[]`
- `s := []int{}`：效果和 `make([]int, 0)` 一样，但不够直接，一般不推荐

大部分场景用 `var s []int` 就够了，后面 `append` 时自然会分配内存。

### Q10：三索引切片 `s[low:high:max]` 什么时候用？

答：当你截取了一个子切片，并且子切片后面会被 append 的时候。限制 cap 可以确保子切片 append 时触发扩容而不是在父切片的底层数组上直接写，防止父子切片的数据互相污染。并发场景尤其重要。

---

## 一句话总结

**切片 header 是值拷贝，底层数组是共享的——记住这一点，切片 90% 的行为都能推理出来。**

