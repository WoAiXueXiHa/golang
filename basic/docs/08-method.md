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
func Add(x, y int) int {
    return x + y
}
```

方法是这样：

```go
func (c *Counter) AddPtr() {
    c.Count++
}
```

`(c *Counter)` 叫接收者（receiver）。它表示这个方法属于 `Counter` 这个类型。

如果去掉接收者，方法签名就和普通函数一模一样：

```go
// 方法的"函数形态"
func AddPtr(c *Counter) {
    c.Count++
}
```

```text
普通函数：
Add(x, y)
   |
   v
函数独立存在，不属于任何类型

方法：
c.AddPtr()
   |
   v
AddPtr 通过接收者绑定到了 Counter 类型上
```

**总结：** 方法 = 函数 + 接收者。接收者决定了这个方法属于哪个类型，也决定了方法里面操作的是值副本还是原对象。

---

## 2. 值接收者和指针接收者

代码里有两个方法：

```go
func (c *Counter) AddPtr() {   // 指针接收者
    c.Count++
}

func (c Counter) AddValue() {  // 值接收者
    c.Count++
}
```

核心区别：

- **指针接收者 `*Counter`**：方法拿到的是原对象的地址，改的就是原对象。
- **值接收者 `Counter`**：方法拿到的是原对象的副本，改副本不影响原对象。

```text
调用 c.AddPtr()

指针接收者拿到地址
+----------+        +----------+
| 地址     | -----> | Count: 1 |
+----------+        +----------+
                         |
                         v
                     Count: 2  ← 原对象被修改

调用 c.AddValue()

值接收者拿到副本
+----------+        +----------+
| 原对象   |        | 副本     |
| Count: 2 |        | Count: 2 |
+----------+        +----------+
                         |
                         v
                     副本变 3

原对象还是 Count: 2  ← 副本的修改不影响原对象
```

**总结：** 值接收者改副本，指针接收者改原对象。方法要修改结构体内部状态时，用指针接收者；结构体很大时也优先用指针接收者避免复制开销。只读且对象很小时，值接收者足够。

---

## 3. 编译器的自动转换（语法糖）

代码里这句可以正常运行：

```go
c.AddPtr()
```

虽然 `c` 是值类型 `Counter`，而 `AddPtr` 的接收者是 `*Counter`，但 Go 编译器会自动帮你取地址：

```go
(&c).AddPtr()  // 编译器实际做的转换
```

反过来也一样——指针变量也能调用值接收者方法：

```go
pc.AddValue()
// 编译器转换为：
(*pc).AddValue()
```

```text
值变量调用指针方法：
c.AddPtr()
   |
   v  编译器自动取地址
(&c).AddPtr()

指针变量调用值方法：
pc.AddValue()
   |
   v  编译器自动解引用
(*pc).AddValue()
```

**总结：** 调用方法时 Go 会自动帮你取地址或解引用，写起来很舒服。但判断一个类型是否实现接口时，看的是方法集，不能只看"能不能调用"——这一点在后面方法集章节会详细讲。

---

## 4. 组合代替继承

Go 没有 `extends` 关键字。想让一个结构体复用另一个结构体的字段和方法，就直接把它嵌进去：

```go
type Engine struct {
    Power int
}

func (e *Engine) Start() {
    fmt.Println("引擎启动...")
}

type Car struct {
    Name string
    Engine        // 匿名嵌入，也叫"组合"
}
```

组合之后，内层字段和方法会被**提升**到外层，可以直接访问：

```go
fmt.Println(myCar.Power)  // 等价于 myCar.Engine.Power
```

```text
Car
+----------------------+
| Name: "Audi"         |
| Engine               |
| +------------------+ |
| | Power: 500       | |
| | Start()          | |
| +------------------+ |
| Start() ← 自己也有  |
+----------------------+

myCar.Power        → myCar.Engine.Power   (字段提升)
myCar.Start()      → Car.Start() 优先    (同名方法，外层优先)
myCar.Engine.Start() → Engine.Start()    (显式调用内层方法)
```

如果外层自己也有同名方法（比如 `Car.Start` 和 `Engine.Start`），优先调用外层的——这不叫 override，Go 里没有继承链，只是"外层方法遮蔽了内层方法"。

**总结：** Go 靠组合而非继承来复用能力。匿名嵌入后，内层的字段和方法被提升到外层；同名时外层优先，可以显式指定内层来调用被遮蔽的方法。

---

## 5. 方法集

代码里有一个接口：

```go
type Payer interface {
    Pay(amount int)
}
```

`WeChatPay` 的 `Pay` 方法是用指针接收者定义的：

```go
func (w *WeChatPay) Pay(amount int) {
    w.Balance -= amount
}
```

这就导致了一个关键结果：

- `*WeChatPay` 实现了 `Payer` ✅
- `WeChatPay` 没有实现 `Payer` ❌

```text
类型 WeChatPay 的方法集
+----------------+
| (空)           |  ← 值类型只包含值接收者方法
+----------------+

类型 *WeChatPay 的方法集
+----------------+
| Pay(amount)    |  ← 指针类型包含值接收者方法 + 指针接收者方法
+----------------+

Pay 是指针接收者方法
所以只有 *WeChatPay 实现了 Payer 接口
```

```go
var p2 Payer = &wxWallet  // ✅ 编译通过
var p1 Payer = wxWallet   // ❌ 编译错误！
```

为什么值类型不能包含指针方法？因为值类型变量是一份独立的数据——如果允许值类型调用指针方法，意味着要取这个值的地址；但在某些场景下（比如字面量、map 中的值），这个地址不可靠，Go 索性在设计上就不让值类型的方法集包含指针接收者方法。

编译器允许 `wxWallet.Pay(20)` 这种调用（语法糖帮你取了地址），但接口赋值不看语法糖，只看方法集。

**总结：** 方法集 = 类型"真正拥有"的方法列表。值类型只包含值接收者方法，指针类型同时包含值和指针接收者方法。接口赋值时，Go 严格检查方法集，语法糖帮不了忙。

---

## 6. 汇编视角：方法底层到底发生了什么

上面讲了这么多"值接收者改副本"、"指针接收者改原对象"、"编译器自动转换"、"接口动态分发"——这些不是靠信念相信的，汇编代码里写得明明白白。

下面以 `method.go` 的反汇编代码为基础，逐一对照验证。[^1]

[^1]: 反汇编命令 `go tool compile -S method.go`，为便于阅读保留了核心指令，省略了 FUNCDATA/PCDATA 等元数据。完整反汇编文件见 `basic/code/08-method/method_analysis.s`。

### 6.1 指针接收者：拿到地址，直接写入原内存

```asm
; main.(*Counter).AddPtr — 指针接收者方法
; method.go:14  func (c *Counter) AddPtr()
TEXT  main.(*Counter).AddPtr(SB), NOSPLIT|NOFRAME|ABIInternal, $0-8
  MOVQ  AX, main.c+8(SP)    ; ① 把接收者指针存到栈上
  TESTB AL, (AX)             ; ② nil 指针检查
  TESTB AL, (AX)             ;     (两次用于并发安全)
  MOVQ  (AX), CX             ; ③ CX = *AX  ← 通过地址读 Count 的值
  INCQ  CX                   ; ④ CX = CX + 1
  MOVQ  CX, (AX)             ; ⑤ *AX = CX  ← 通过地址写回
  RET
```

**逐条解读：**

1. `AX` 里放的是接收者指针——也就是 `c` 的地址。`ABIInternal` 意味着参数通过寄存器传递，不再走栈。
2. `NOSPLIT|NOFRAME` 说明这个方法极短（只有 19 字节），不需要栈帧——连 `PUSHQ BP` 都省了。
3. `MOVQ (AX), CX` — 这条指令是"指针接收者改原对象"的铁证：它先**解引用**指针拿到原始值，放到 CX。
4. `INCQ CX` — 在 CX 上自增。
5. `MOVQ CX, (AX)` — 然后把结果**写回指针指向的地址**。

```text
AX → 0xc0000120a8  (Counter 变量的地址)
      |
      v
      MOVQ (AX), CX  → CX = [0xc0000120a8] = 1
      INCQ CX        → CX = 2
      MOVQ CX, (AX)  → [0xc0000120a8] = 2  ← 原内存被修改！
```

这就是"指针接收者改原对象"的物理本质：**写操作的目标地址就是调用方那块内存。**

### 6.2 值接收者：拿到副本，改了个寂寞

```asm
; main.Counter.AddValue — 值接收者方法
; method.go:19  func (c Counter) AddValue()
TEXT  main.Counter.AddValue(SB), NOSPLIT|NOFRAME|ABIInternal, $0-8
  MOVQ  AX, main.c+8(SP)    ; ① 把值存到栈上 (这是副本！)
  INCQ  AX                   ; ② AX = AX + 1  ← 直接改寄存器
  MOVQ  AX, main.c+8(SP)    ; ③ 把新值存回栈  ← 只改了栈上的副本
  RET
```

**逐条解读：**

1. `AX` 里放的不是地址，而是 `Count` 的**值本身**（比如 `1`）。值接收者把整个结构体（这里只有一个 int 字段）通过寄存器传进来。
2. `INCQ AX` — 直接在寄存器上自增。注意这里没有任何解引用操作！
3. `MOVQ AX, main.c+8(SP)` — 把自增后的值写回**栈上的参数位置**。这个位置在函数返回后就释放了。

```text
调用方栈                          AddValue 的栈
+-------------------+            +-------------------+
| c.Count = 1       | --复制--> | main.c+8(SP) = 1  |
+-------------------+            +-------------------+
                                       |
                                       v INCQ AX
                                 main.c+8(SP) = 2  ← 只改了这里
                                       |
                                       v RET → 栈帧销毁，2 跟着消失

调用方 c.Count 纹丝不动，还是 1
```

**和指针版本的对比：**

| 操作 | `AddPtr` (指针) | `AddValue` (值) |
|------|----------------|-----------------|
| AX 存的是什么 | 地址 | 值 |
| 读字段 | `MOVQ (AX), CX` 解引用 | AX 就是值，无需解引用 |
| 改字段 | `MOVQ CX, (AX)` 写回原地址 | `INCQ AX` 改寄存器，写回栈 |
| nil 检查 | 有 (`TESTB AL, (AX)`) | 无（值类型不可能 nil） |
| 大小 | 19 字节 | 14 字节 |

**值接收者比指针接收者还少了 5 字节**——因为它不需要解引用，不需要 nil 检查。但如果结构体很大（比如几十个字段），值接收者每次调用都要复制整个结构体，开销反而更大。

### 6.3 编译器自动转换：汇编里看得清清楚楚

代码里有这两行：

```go
c.AddPtr()    // c 是 Counter 值类型，AddPtr 要 *Counter
c.AddValue()  // c 是 Counter 值类型，AddValue 要 Counter
```

在 `main.main` 的汇编中：

```asm
; c.AddPtr() — 值变量调指针方法
; method.go:73
  LEAQ  main.c+32(SP), AX       ; AX = &c  ← 编译器自动取了地址！
  CALL  main.(*Counter).AddPtr(SB)

; c.AddValue() — 值变量调值方法
; method.go:77
  MOVQ  main.c+32(SP), AX       ; AX = c.Count (值本身)
  CALL  main.Counter.AddValue(SB)
```

**看，`c.AddPtr()` 这一行高级语言的调用，编译器干了两件事：**

1. `LEAQ main.c+32(SP), AX` — 用 `LEAQ`（Load Effective Address）取了 `c` 的地址
2. `CALL main.(*Counter).AddPtr` — 调用指针版本的函数

而 `c.AddValue()` 则是直接把值传给 AX，调用值版本。**一切转换在编译期完成，零运行时开销。**

### 6.4 编译器自动生成的方法包装器

你可能会问：`(*Counter).AddValue` 存在吗？`AddValue` 是用 `Counter`（值）定义的，不是 `*Counter`。但如果写 `(&c).AddValue()` 呢？

编译器会自动生成一个**包装方法**（WRAPPER），在反汇编的末尾能找到：

```asm
; 编译器自动生成的包装器，源文件里不存在这段代码
; <autogenerated>:1
TEXT  main.(*Counter).AddValue(SB), DUPOK|WRAPPER|ABIInternal, $24-8
  MOVQ  AX, main.c+32(SP)          ; 存指针
  TESTQ AX, AX                     ; 检查指针是否为 nil
  JNE   notnil                     ; 不为 nil，继续
  CALL  runtime.panicwrap(SB)      ; nil！panic
notnil:
  TESTB AL, (AX)                   ; 二次 nil 屏障
  MOVQ  (AX), AX                   ; ★ AX = *AX  解引用，取出值
  MOVQ  AX, main..autotmp_1+8(SP)  ; 把值作为参数传递
  CALL  main.Counter.AddValue(SB)  ; 调用真正的值接收者方法
  ...
  RET
```

```text
你的代码:  (&c).AddValue()
              |
              v
编译器自动生成包装器:
  1. 检查指针不是 nil → 是 nil 就 panic
  2. *指针 解引用拿到值副本
  3. 调用真正的 Counter.AddValue(值)
```

这个包装器标了 `DUPOK`（重复允许），意味着每个包都生成一份，链接时只保留一份。

### 6.5 接口动态分发：itab 与方法指针

看 `main.Payer.Pay`——接口的包装方法：

```asm
; main.Payer.Pay — 接口方法分发
; <autogenerated>:1
TEXT  main.Payer.Pay(SB), DUPOK|WRAPPER|ABIInternal, $24-24
  MOVQ  AX, main.~p0+32(SP)       ; AX = 接口值的 itab 指针（类型表）
  MOVQ  BX, main.~p0+40(SP)       ; BX = 接口值的数据指针（指向 WeChatPay 实例）
  MOVQ  CX, main.amount+48(SP)    ; CX = amount 参数
  TESTB AL, (AX)                   ; nil 检查
  MOVQ  24(AX), DX                ; ★ DX = itab[3] = Pay 方法的函数指针
  MOVQ  BX, AX                    ; 重排参数：receiver 放到第一参数位
  MOVQ  CX, BX                    ; amount 放到第二参数位
  CALL  DX                        ; ★ 间接调用！跳转到实际方法的地址
  ...
  RET
```

```text
接口值在运行时的结构（两个指针）:

  Payer 接口变量
  +--------+--------+
  | data   | itab   |
  +--------+--------+
      |        |
      v        v
  WeChatPay  +-----------------+
  实例       | InterfaceType*  |  ← itab[0]
             | ConcreteType*   |  ← itab[1]
             | Hash            |  ← itab[2]
             | Pay 的函数指针  |  ← itab[3] ★ 动态分发的关键
             +-----------------+

CALL DX 时，DX 的值来自 itab[3]——而这个指针在编译期就已经确定
是 main.(*WeChatPay).Pay 的地址。
```

在反汇编的 itab 段中，可以看到这个静态表：

```asm
; go:itab.*main.WeChatPay,main.Payer
go:itab.*main.WeChatPay,main.Payer SRODATA dupok size=32
  rel 0+8  t=R_ADDR type:main.Payer+0           ; itab[0]: Payer 接口类型
  rel 8+8  t=R_ADDR type:*main.WeChatPay+0       ; itab[1]: *WeChatPay 具体类型
  rel 24+8 t=RelocType(-32767) main.(*WeChatPay).Pay+0  ; itab[3]: Pay 函数指针
```

**关键信息：** itab 是编译器生成的只读静态表。当代码写出 `var p Payer = &wxWallet` 时，编译器把 `*WeChatPay` 的 itab 指针和数据指针一起打包进接口值。后续调用 `p.Pay(20)` 时，CPU 执行的是 `CALL DX`，DX 来自 itab——这就是 Go 接口动态分发的完整链路。

### 6.6 方法集在 type descriptor 中的物理证据

反汇编里，每个类型都有一个 `type` 描述符，记录了它的方法集。来看看 `Counter` 的值类型和指针类型的差异：

```asm
; type:main.Counter — 值类型的方法集
type:main.Counter SRODATA size=136
  ...
  rel 120+4 t=R_ADDROFF type:.namedata.AddValue.+0     ; ← 只有 AddValue
  rel 124+4 t=R_METHODOFF type:func()+0
  rel 128+4 t=R_METHODOFF main.(*Counter).AddValue+0
  rel 132+4 t=R_METHODOFF main.Counter.AddValue+0
  ; 注意：没有 AddPtr！

; type:*main.Counter — 指针类型的方法集
type:*main.Counter SRODATA size=104
  ...
  rel 72+4  t=R_ADDROFF type:.namedata.AddPtr.+0       ; ← 有 AddPtr
  rel 76+4  t=R_METHODOFF type:func()+0
  rel 80+4  t=R_METHODOFF main.(*Counter).AddPtr+0
  rel 84+4  t=R_METHODOFF main.(*Counter).AddPtr+0
  rel 88+4  t=R_ADDROFF type:.namedata.AddValue.+0     ; ← 也有 AddValue
  rel 92+4  t=R_METHODOFF type:func()+0
  rel 96+4  t=R_METHODOFF main.(*Counter).AddValue+0
  rel 100+4 t=R_METHODOFF main.(*Counter).AddValue+0
```

```text
type:main.Counter (值)       type:*main.Counter (指针)
+---------------------+     +---------------------------+
| AddValue ✓          |     | AddPtr  ✓                 |
| AddPtr   ✗          |     | AddValue ✓ (自动生成包装器)|
+---------------------+     +---------------------------+

结论：值类型的方法集只登记了 AddValue；
      指针类型的方法集同时登记了 AddPtr 和 AddValue。
```

**这就是"为什么 `var p Payer = wxWallet` 会报错"的底层答案：**

`Payer` 接口要求 `Pay(amount int)`。编译器在检查 `WeChatPay` 是否实现 `Payer` 时，会去读 **`type:main.WeChatPay`** 的方法集——里面没有 `Pay`（因为 `Pay` 定义在 `*WeChatPay` 上）。而 `type:*main.WeChatPay` 的方法集里有 `Pay`。所以只有 `*WeChatPay` 实现了 `Payer`，值类型 `WeChatPay` 没有。

这不是"编译器不够聪明"，而是**设计者故意为之**：值类型的方法集里不收指针接收者方法，保证你拿到的是一个明确、可寻址的实体，而不是一个临时副本。

---

## 易错点

| # | 易错场景 | 为什么 |
|---|---------|--------|
| 1 | `c.AddPtr()` 能调用，不代表 `Counter` 实现了需要指针方法的接口 | 调用有语法糖，接口赋值只看方法集 |
| 2 | 值接收者方法里修改字段，不影响原对象 | 汇编已证明：改的是寄存器/栈副本 |
| 3 | 指针接收者常用于修改状态或结构体较大时 | 避免复制开销，且可以保证修改到位 |
| 4 | 匿名嵌入不是继承 | Go 没有父类子类，只是字段/方法提升 + 同名遮蔽 |
| 5 | 接口赋值看方法集，不看语法糖 | 编译器在接口检查时不帮你取地址 |
| 6 | 给接口变量赋值时，接收者决定谁能"上岗" | `*T` 能当 `T` 的接口，反过来不行（如果方法是指针接收者） |

---

## 快问快答

### Q1：Go 方法和函数有什么区别？

方法就是带接收者的函数。从反汇编看，它们的 TEXT 符号名、参数传递方式完全一致，唯一的区别是方法多了一个隐式的 receiver 参数，并且被注册到类型的描述符里。

### Q2：值接收者和指针接收者怎么选？

需要修改原对象 → 指针接收者。结构体较大不想复制 → 指针接收者。只读且对象很小（几个 int 这种） → 值接收者足够。从汇编可以看到，值接收者版本更短（少了解引用和 nil 检查），但付出的是复制成本。

### Q3：为什么 `var p Payer = wxWallet` 会报错？

因为 `Pay` 定义在 `*WeChatPay` 上（指针接收者），只有 `*WeChatPay` 的方法集包含 `Pay`。编译器做接口检查时，去 `type:main.WeChatPay` 的方法集里找 `Pay`，找不到。`wxWallet.Pay(20)` 能调是编译器帮你取了地址——接口检查不吃这套。

### Q4：匿名嵌入的字段同名怎么办？

外层优先。如果 `Car` 自己有 `Start()`，`myCar.Start()` 调的是 `Car.Start`，不是 `Engine.Start`。想调引擎的？写 `myCar.Engine.Start()`。这跟继承的 override 不一样——Go 没有虚函数表，选择在编译期就确定了。

---

## 一句话总结

Go 方法就是带接收者的函数；组合负责能力复用；方法集在编译期写死在 type descriptor 里，决定谁有资格实现接口。汇编不会说谎——指针接收者写原内存，值接收者改栈副本，接口调用靠 itab 表做间接跳转。
