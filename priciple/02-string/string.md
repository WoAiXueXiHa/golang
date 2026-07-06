# 什么是 string

## 这一章要记住什么

- string 底层是 `{str unsafe.Pointer, len int}`，str 指向底层字节数组，len 存的是**字节数**不是字符数
- string 不可变（immutable），所以在只读内存区域分配，线程安全、哈希稳定、子串可共享底层内存
- `for range` 遍历字符串自动按 UTF-8 解码为 rune，普通 `for i` 逐字节遍历，中文场景差异巨大
- string 和 `[]byte` 互转**一定发生内存拷贝**，想避免就用 `unsafe`（但要清楚风险）
- 字符串拼接首选 `strings.Builder`（零拷贝 `String()`），已知切片用 `strings.Join`（一次分配），禁止循环里用 `+` 或 `fmt.Sprintf`
- **子串 `s[:n]` 共享底层内存**，大字符串截一小段会导致整个大字符串无法 GC，Go 1.18+ 用 `strings.Clone` 解决

---

在 `src/builtin/builtin.go` 中这样定义：

```golang
// string is the set of all strings of 8-bit bytes, conventionally but not
// necessarily representing UTF-8-encoded text. A string may be empty, but
// not nil. Values of string type are immutable.
type string string
```

- 字符串是所有 8bit 字节的集合，但不一定是 UTF-8 编码的文本
- 字符串可以为空，但是不能为 `nil`
- 字符串类型的值是不可变的

本质是一个字符数组，每个字符在存储时都对应一个整数，也可能对应多个整数

对于 C 语言的 string，每个字符串结尾必须加 `\0`，表示这个字符串结束了

Go 不是这样设计，Go 使用一个 len （int类型）存这个字符串的总字节数

在`src/runtime/string.go`文件中，对 string 结构体进行了定义：

```golang
type stringStruct struct {
    str unsafe.Pointer
    len int
}
```

- str 指针指向字符串首地址
- len 表示字符串的长度

> **注意**：`stringStruct` 是 runtime 内部使用的结构体，用户代码无法直接访问

len 表示的是这个字符串占用的字节数

一个常见误区是以为 `len` 返回的是字符个数，实际上它返回的是**底层占用的字节数**。对于中文字符（UTF-8 编码下每个中文字符占 3 个字节），差异非常明显：

```go
package main

import (
    "fmt"
    "unicode/utf8"
)

func main() {
    s := "你好"
    fmt.Println(len(s))                    // 6（字节数），不是 2
    fmt.Println(utf8.RuneCountInString(s)) // 2（字符数）
}
```

要获取实际的字符数量，需要使用 `utf8.RuneCountInString`。

看代码：

```golang
package main

import "fmt"

func main() {
    word := "Hello, World"
    for _, v := range word {
        fmt.Printf("%d\n", v)
    }
}
```

输出：

![image-20260704152406209](https://gitee.com/binary-whispers/pic/raw/master///20260704152408999.png)

`for range` 遍历字符串时，Go 会自动按 **rune（Unicode 码点）** 解码，`v` 是解码后的 rune 值（int32），索引 `i` 是当前 rune 在字符串中的**字节偏移量**。相比之下，普通 `for i := 0; i < len(s); i++` 是逐字节遍历，遇到中文会得到乱码的单个字节。

## rune 类型与三种遍历方式对比

`rune` 是 Go 的内置类型，本质是 `int32` 的别名，代表一个 **Unicode 码点（Code Point）**。它和 `byte`（`uint8` 的别名，代表一个字节）是两个层面的东西：

```go
type rune = int32   // 四个字节，能装下所有 Unicode 字符
type byte = uint8   // 一个字节，只够 ASCII
```

遍历字符串有三种姿势，行为完全不同：

```go
s := "你好Go"

// 方式一：逐字节遍历 —— 每步拿到一个 byte
for i := 0; i < len(s); i++ {
    fmt.Printf("%x ", s[i])  // e4 bd a0 e5 a5 bd 47 6f
}

// 方式二：逐 rune 遍历 —— 每步拿到一个完整字符
for _, r := range s {
    fmt.Printf("%c ", r)  // 你 好 G o
}

// 方式三：带索引的 range —— i 是字节偏移量，r 是 rune
for i, r := range s {
    fmt.Printf("s[%d]=%c ", i, r)  // s[0]=你 s[3]=好 s[6]=G s[7]=o
}
```

注意方式三中索引不是连续的 0,1,2,3，而是 0,3,6,7——"你"占 3 字节，"好"占 3 字节，"G"和"o"各占 1 字节。

### 什么时候转 `[]rune`？

当需要按**字符位置**（而非字节位置）修改或访问时，必须先转 `[]rune`：

```go
s := "你好世界"
r := []rune(s)    // 分配新内存，拷贝：['你', '好', '世', '界']
r[1] = '嗨'       // 修改第二个字符
s = string(r)     // 转回去："你好嗨界"
```

如果不转 `[]rune` 直接 `s[1]`，拿到的是字节 `0xbd`，不是"好"字，改它也只会破坏 UTF-8 编码。

### 总结一下

> `rune` = `int32` = Unicode 码点。`for range` 按 rune 遍历，索引是字节偏移量。需要按字符位置增删改时，先转 `[]rune` 再操作。`len()` 返回的是字节数，要用 `utf8.RuneCountInString` 或 `len([]rune(s))` 拿字符数。

---

以下是底层原理图：

![image-20260704155932499](https://gitee.com/binary-whispers/pic/raw/master///20260704155935913.png)

值得注意的是，**Go 认为字符串内容是不会被修改的，所以会把字符串分配到只读内存区域**。这样设计有几个关键原因：

1. **线程安全**：不可变意味着任意多个 goroutine 并发读取同一字符串时无需加锁
2. **哈希稳定**：string 作为 map 的 key 时，其哈希值不会改变，保证了 map 的正确性
3. **子串共享内存**：`s[1:3]` 这种取子串的操作是 O(1) 的，新字符串直接复用原串的底层内存，无需拷贝

字符串变量可以指向同一块底层内存，共享底层内容，如下图所示：

![image-20260704154127004](https://gitee.com/binary-whispers/pic/raw/master///20260704154128583.png)

正因为是共享底层内存的，如果允许通过 s1 修改内容，s2 也会随之变化，这样的风险无法预知，所以 Go 从根本上禁止了这种操作

如果非要修改，可以**给变量赋新的值，让其指针指向新的内存空间**：

![image-20260704154637568](https://gitee.com/binary-whispers/pic/raw/master///20260704154639314.png)

以上是 string 的一些基本性质

## 子串共享内存与内存泄漏

因为 string 不可变，取子串 `s[low:high]` 不需要拷贝数据——新字符串的 `str` 指针直接指向原字符串底层数组的偏移位置，O(1) 完成。

```go
s := "hello, world"
sub := s[0:5]  // "hello" —— 和 s 共享底层 13 字节的内存
```

```text
s:   +------------------+
     | str: 0x...100    |---+
     | len: 13          |   |
     +------------------+   |
                            v
sub: +------------------+   0x...100  0x...10D
     | str: 0x...100    |   +--------------------------------+
     | len: 5           |   | h | e | l | l | o | , |  | w | ...
     +------------------+   +--------------------------------+
                            sub 能看到的范围 ↑ (len=5)
                            整个底层数组仍被 s 或 sub 引用 ↑
```

这在大多数场景下是好事——省内存、速度快。但有一个坑：**如果从一个大字符串截一小段，且原字符串被 GC 了，底层大数组仍然被小子串"锚定"，无法回收**。

```go
// 假设从一个大文件读了 1GB 的内容
content := readHugeFile()      // 1GB
firstLine := content[:100]     // 只取前 100 字节
// 即使 content 出了作用域，那 1GB 底层数组仍然被 firstLine 引用着
```

解决方案——强制拷贝一份：

```go
// 方式一：string ← []byte ← string（两次拷贝，但有效）
firstLine := string([]byte(content[:100]))

// 方式二：Go 1.18+ 直接用 strings.Clone（推荐）
firstLine := strings.Clone(content[:100])
```

`strings.Clone` 内部就是做了 `string([]byte(s))` 的事，但语义更清晰——"我要一份独立的副本"。

### 总结一下

> 子串和原串共享底层内存（O(1) 操作）。大串截小串会导致整个大串无法 GC。想断开引用，Go 1.18+ 用 `strings.Clone`，之前用 `string([]byte(s))`。

---

# string 和 []byte的转换

还有种方式，将字符串强转为切片，通过索引修改切片，再转换回字符串：

```golang
package main

import "fmt"

func main() {
    s := "Hello"
    strByte := []byte(s)
    strByte[0] = 'h'
    fmt.Println(string(strByte))
}

```

输出：

```text
hello
```

需要注意的是：**源字符串并没有发生变化，我们得到的只是 s 字符串的一个拷贝**

## 转化原理

string 和 []byte 的转化会发生一次拷贝，申请一块新的切片空间

byte 切片转为 string 的过程：

- 新申请切片内存空间，构建内存地址为 addr， 长度为 len
- 构建 string 对象，指针地址为 addr， len 字段赋值为 len
- 将源切片中数据拷贝到新申请的 string 中指针指向的内存空间

![image-20260704161209935](https://gitee.com/binary-whispers/pic/raw/master///20260704161212024.png)

string 转为 byte 切片的过程：

- 新申请切片内存空间
- 将 string 中指针指向内存区域的内容拷贝到新切片

![image-20260704161326296](https://gitee.com/binary-whispers/pic/raw/master///20260704161328186.png)

## unsafe 零拷贝转换（面试高频）

标准的 `string(dat)` 和 `[]byte(s)` 都会发生内存拷贝，在需要极致性能的场景（比如高频调用、大字符串处理），可以用 `unsafe` 做到零拷贝：

```go
import "unsafe"

// []byte → string：零拷贝，仅构造 header
func BytesToString(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}

// string → []byte：零拷贝，仅构造 header
func StringToBytes(s string) []byte {
    return *(*[]byte)(unsafe.Pointer(&s))
}
```

原理很简单——string 和 slice 的 header 内存布局高度相似：

```text
string header:                    slice header:
+------------------+              +------------------+
| str unsafe.Pointer|              | Data unsafe.Pointer|
+------------------+              +------------------+
| len int           |              | Len int           |
+------------------+              +------------------+
                                   | Cap int           |
                                   +------------------+
```

就是直接把 string header 的 16 字节强转为 slice header 的 24 字节（多出来的 Cap 会从相邻内存读取，等于自动补了一个 `Cap = len`），反过来也一样。

**但这里有三个必须说清楚的致命风险：**

1. **修改零拷贝得到的 `[]byte` 会导致 panic 或更严重的未定义行为**。string 的底层内存在只读区，通过零拷贝拿到的 `[]byte` 虽然能编译通过，但一旦 `dat[0] = 'x'`，就会尝试写只读内存，直接 `SIGSEGV` 崩溃。

2. **原 `[]byte` 被 GC 回收后，零拷贝出的 string 变成悬空指针**。如果你把 `BytesToString(buf)` 的结果存下来用，而 `buf` 的底层数组被 GC 了，string 的 `str` 指针就指向了已回收的内存。

3. **标准库中 `strings.Builder.String()` 用的就是这个技巧**，但它的前提是 `Builder` 内部持有 `[]byte`，保证底层数组不会在 string 存活期间被 GC。

```go
// strings.Builder 源码中的做法（简化版）
func (b *Builder) String() string {
    // buf 是 Builder 持有的 []byte，不会被 GC
    return *(*string)(unsafe.Pointer(&b.buf))
}
```

### 什么时候可以考虑用？

- 临时使用、用完即扔的场景（`string(buf)` 用作 map key 查询，用完就不管了）
- 你确信原 `[]byte` 的生命周期覆盖了 string 的使用周期
- 你**绝对不会**修改零拷贝拿到的 `[]byte`

**大多数业务代码不需要也不应该用 unsafe 转换**——标准转换的内存开销在绝大部分场景下可以忽略，换来的是安全和可读性。

### 总结一下

> unsafe 零拷贝通过强转 header 实现，本质是欺骗类型系统。从 `[]byte` 转来的 string 不能改内容；从 string 转来的 `[]byte` 修改会崩溃。只有标准库这类对生命周期有完全掌控的场景才安全使用。

---

# 字符串拼接

字符串拼接会有内存的拷贝，存在性能损耗，常见有以下方式：

- +操作符
- fmt.Sprintf
- bytes.Buffer
- strings.Builder
- append
- string.Join

使用代码测试一下：

```golang
package main

import (
    "bytes"
    "fmt"
    "strings"
    "testing"
)

// 基础配置：拼接 1000 个短字符串
const (
    loopCount = 1000
    subStr    = "go"
)

// 1. + 操作符
func BenchmarkPlus(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var s string
        for j := 0; j < loopCount; j++ {
            s += subStr // 每次都会产生新字符串，旧字符串变垃圾，高频触发内存拷贝
        }
    }
}

// 2. fmt.Sprintf
func BenchmarkSprintf(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var s string
        for j := 0; j < loopCount; j++ {
            s = fmt.Sprintf("%s%s", s, subStr) // 内部涉及接口反射和动态分配，最慢
        }
    }
}

// 3. bytes.Buffer
func BenchmarkBytesBuffer(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var buf bytes.Buffer
        for j := 0; j < loopCount; j++ {
            buf.WriteString(subStr)
        }
        _ = buf.String() // 最后一次性转换为 string
    }
}

// 4. strings.Builder
func BenchmarkStringsBuilder(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var builder strings.Builder
        for j := 0; j < loopCount; j++ {
            builder.WriteString(subStr)
        }
        _ = builder.String() // 底层通过 unsafe 转换，零拷贝指针，性能极高
    }
}

// 5. append (切片转字符串)
func BenchmarkAppend(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var buf []byte
        for j := 0; j < loopCount; j++ {
            buf = append(buf, subStr...)
        }
        _ = string(buf) // 这一步依然会发生一次内存拷贝
    }
}

// 6. strings.Join
func BenchmarkStringsJoin(b *testing.B) {
    // 先准备好切片数据
    slice := make([]string, loopCount)
    for i := 0; i < loopCount; i++ {
        slice[i] = subStr
    }

    b.ResetTimer() // 重置时间，扣除准备切片的耗时
    for i := 0; i < b.N; i++ {
        _ = strings.Join(slice, "") // 内部提前计算总长度并预分配内存，适合已知切片拼接
    }
}

```

输出：

```bash
[vect@ubuntu-dev ~/golang/priciple/02-string/demo3]$ go test -bench=. -benchmem main_test.go

goos: linux

goarch: amd64

cpu: Intel(R) Xeon(R) Gold 6148 CPU @ 2.40GHz

BenchmarkPlus-2                     3909            288534 ns/op       1063873 B/op         999 allocs/op

BenchmarkSprintf-2                  3244            375520 ns/op       1080060 B/op        1999 allocs/op

BenchmarkBytesBuffer-2            158611              7444 ns/op          6080 B/op           7 allocs/op

BenchmarkStringsBuilder-2         339829              3072 ns/op          5368 B/op          10 allocs/op

BenchmarkAppend-2                 525447              2490 ns/op          7416 B/op          11 allocs/op

BenchmarkStringsJoin-2            134302              8968 ns/op          2048 B/op           1 allocs/op

PASS

ok      command-line-arguments  7.393s
```

分析：

1. **`+` 和 `Sprintf` 直接崩掉**：
   * `BenchmarkPlus` 的 `999 allocs/op` 说明 1000 次循环里，**几乎每拼接一次都在堆上申请一次内存**。
   * `BenchmarkSprintf` 的 `1999 allocs/op` 翻倍了，因为除了拼接，还要承担格式化参数逃逸到堆上的额外分配，耗时最长（375us）。
2. **`StringsJoin` 内存控制无敌**：
   * `1 allocs/op` 证明了它的**一次性预分配**。无论拼接多少，只申请一次。
3. **`StringsBuilder` 相比 `Buffer` 的优势**：
   * `Builder` 耗时（3072 ns）只有 `Buffer`（7444 ns）的一半。这就是最后一步**零拷贝**省下来的 CPU 开销。
4. **`Append` 速度最快的原因**：
   * `2490 ns/op` 拿了第一，这是因为内置的 `append` 有运行时（runtime）专门的汇编级别优化，且切片扩容策略和底层容量对齐极度灵敏。但看内存（7416 B/op）能发现，它最后强转 string 多拷贝了一次，所以内存占用比 Builder 略大。
   * 一个值得注意的细节：1000 次循环却只产生了 10~11 次内存分配，这是因为 `[]byte` 的扩容是指数级增长的 —— 容量小于 1024 时每次翻倍，超过后每次增加 25%。所以实际扩容次数远小于循环次数。

总结：

| **拼接方式**                  | **耗时 (ns/op)** | **内存分配次数 (allocs/op)** | **底层核心原理**                                             | **适用场景与局限**                                           |
| ----------------------------- | ---------------- | ---------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| **`BenchmarkAppend`**         | 2490 ns          | 11 次                        | 手动维护 `[]byte` 切片，利用 runtime 内置的 `append` 进行快速扩容。最后 `string(buf)` **触发一次全量内存拷贝**。 | **偏底层字节处理**。当后续还需要对字节切片进行微调、或是纯字节流操作时适用。 |
| **`BenchmarkStringsBuilder`** | 3072 ns          | 10 次                        | 底层同样是 `[]byte`。`String()` 时利用 `unsafe.Pointer` 直接共享底层数组指针，**零拷贝返回**。若提前知道总长度，调用 `Grow(n)` 预分配可进一步减少扩容次数。 | **绝大多数动态/循环拼接的首选**。不知道最终长度，需要不断往里塞字符串的通用高频场景。 |
| **`BenchmarkBytesBuffer`**    | 7444 ns          | 7 次                         | 经典的字节缓冲区。最后 `buf.String()` 会**重新申请一块新内存**，把所有字节拷贝过去变成不可变 string。 | **I/O 混合场景**。多用于既要拼接字符串，又要和 `io.Reader/Writer`（如网络、文件）做交互的地方。 |
| **`BenchmarkStringsJoin`**    | 8968 ns          | **1 次 👑**                   | 1. 遍历计算所有单项的精确总长度；2. 一次性 `make` 足额空间；3. 拷贝数据并零拷贝转为 string。 | **已有切片数据、或可预知长度**。数据原本就在 `[]string` 里，或者能提前算好长度，用它内存最干净（只有 1 次分配）。 |
| **`BenchmarkPlus`**           | 288534 ns        | 999 次                       | 每次 `+=` 都在堆上开辟新空间，把老 string 和新短串拷贝过去。循环中会导致复杂度退化为 $O(N^2)$。 | **2-3 个已知串简单拼接**。禁止在循环体内使用。单行 `a + b + c` 编译器会优化，效率很高。 |
| **`BenchmarkSprintf`**        | 375520 ns        | 1999 次                      | 内部依赖 `reflect`（反射）动态解析占位符，参数会发生隐式转换并**逃逸到堆上**，伴随大量高频分配。 | **复杂的格式化输出 / 日志**。性能极差，纯粹为了代码可读性服务，高频或循环拼接中绝对不能用。 |

---

## 易错点

1. **`len(s)` 返回字节数不是字符数**。中文字符在 UTF-8 下占 3 个字节，`len("你好")` = 6。要拿字符数用 `utf8.RuneCountInString` 或 `len([]rune(s))`。

2. **循环里用 `+=` 拼接字符串**。每次 `+=` 都会分配新内存并拷贝，1000 次就是 O(N²) 的复杂度。循环拼接永远用 `strings.Builder`。

3. **通过索引 `s[i]` 拿到的是 byte，不是字符**。`s := "你好"; s[0]` 不是 `'你'`，而是 UTF-8 编码的第一个字节 `0xe4`。要按字符索引必须先转 `[]rune`。

4. **零拷贝 `unsafe` 转换后修改 `[]byte` 会崩**。把 string 零拷贝转成的 `[]byte` 仍然指向只读内存，写操作直接 SIGSEGV。反过来，`[]byte` 零拷贝转成 string 后如果原 `[]byte` 被 GC，string 就变成了悬空指针。

5. **子串不拷贝数据，大串截小串导致内存泄漏**。`small := huge[:100]` 后 `huge` 的底层数据不能被 GC。Go 1.18+ 用 `strings.Clone` 解决。

6. **`for range` 和 `for i` 遍历行为不同**。前者按 rune 解码、索引是字节偏移量；后者逐字节遍历。面试经常拿出来比较。

7. **`string` 不能为 nil，但可以是空字符串**。`var s string`（空字符串 `""`）和 nil 不同。判断空字符串用 `s == ""` 或 `len(s) == 0`，不要写成 `s == nil`——编译都过不了。

---

## 快问快答

### Q1：Go 的 string 底层结构是什么样的？

答：就是 16 字节的结构体——一个 `unsafe.Pointer` 指向底层字节数组，一个 `int` 存长度。没有 `cap` 字段（string 不可变，不需要扩容）。

### Q2：为什么 string 设计成不可变的？

答：三个原因：① 线程安全——多 goroutine 读同一字符串无需加锁；② 哈希稳定——string 作为 map key 时哈希值不会变；③ 子串可以安全共享底层内存，O(1) 取子串无需拷贝。

### Q3：`len(s)` 和 `len([]rune(s))` 有什么区别？

答：`len(s)` 返回的是底层字节数（英文 1 字节、中文 3 字节），`len([]rune(s))` 返回的是字符数（中英文都算 1 个字符）。中文 "hello你好" 的 `len()` = 11，`len([]rune())` = 7。

### Q4：`string` 和 `[]byte` 互转一定发生拷贝吗？

答：标准方式一定发生拷贝。string 在只读内存，`[]byte` 在堆上，完全不同的内存区域，必须分配+拷贝。不拷贝只能用 `unsafe` 强转 header，但有只读内存写入崩溃和悬空指针的风险，标准库 `strings.Builder` 就是用 unsafe 做到的零拷贝 `String()`。

### Q5：循环拼接字符串，`strings.Builder` 为什么比 `+` 快那么多？

答：`+` 每次拼接都 `malloc` + `copy` 全量数据，1000 次循环就是 1000 次分配、O(N²) 的拷贝量。`Builder` 底层维护一个 `[]byte`，空间不够才扩容（指数增长），最后 `String()` 用 unsafe 零拷贝返回。1000 次拼接 Builder 约 10 次分配，`+` 约 1000 次。

### Q6：`strings.Builder` 和 `bytes.Buffer` 有什么区别？

答：① Builder 的 `String()` 用 unsafe 零拷贝，Buffer 的 `String()` 会拷贝一次内存；② Builder 不可复用（没有 Reset 的... 其实有 `Reset()`），Buffer 有完整的读写接口；③ Builder 更轻量，只做字符串拼接；Buffer 实现了 `io.Reader/Writer`，适合 I/O 场景。

### Q7：能不能对 string 做 `s[0] = 'x'` 这样的赋值？

答：不能，编译错误——string 是不可变类型，Go 语言层面禁止了索引赋值。想"修改"，只能转 `[]byte` 或 `[]rune`，改完再转回 string——但这个过程产生的是新字符串，原字符串没变。

### Q8：`strings.Clone` 是干什么的？

答：Go 1.18 引入，作用是把子串的数据**强制拷贝一份**，让新字符串拥有独立的底层内存。用于解决"大串截小串导致内存泄漏"的问题。内部实现就是 `string([]byte(s))`。

### Q9：string 能跟 nil 比较吗？

答：不能，编译错误。string 不是引用类型，不能为 nil。空字符串的判断用 `s == ""` 或 `len(s) == 0`。

### Q10：`for _, v := range s` 中 v 是什么类型？索引 i 代表什么？

答：v 是 `rune` 类型（`int32`），代表 Unicode 码点。i 是当前 rune 在字符串中的**字节偏移量**，不是字符位置索引。比如 "你好Go" 中 '好' 的 i=3（"你" 占 3 字节），不是 1。

---

## 一句话总结

**string 不可变、底层只有指针和长度、len 是字节数不是字符数——记住这三个特性，剩下的转换、遍历、拼接性能差异都能从这三个特性推导出来。**