# 07-function：函数的本质是边界、组合和状态

这一章抓四件事：**值传递、函数值、闭包、defer/error。**  
前两个是基础，后两个是 Go 后端工程里真正天天遇到的东西。

## 这份代码最该看什么

```go
func swapByValue(x, y int) {
	x, y = y, x
}

func swapByPtr(x, y *int) {
	*x, *y = *y, *x
}
```

`swapByValue` 改的是副本。  
`swapByPtr` 能改外部变量，但 Go 仍然是值传递：传进去的是“地址值”的副本。

面试会逼你说清楚：**Go 有没有引用传递？**  
没有。Go 只有值传递。

slice、map、chan 看起来像引用，是因为它们的值内部包含指针。复制它们时，复制的是 header，不是底层数据。

## 函数也是值

```go
getSquareRoot := func(x float64) float64 {
	return math.Sqrt(x)
}
```

函数可以赋值给变量。

```go
type fc func(int) int

func callBack(x int, f fc) {
	res := f(x)
	fmt.Println("回调返回值：", res)
}
```

函数也可以作为参数。这是 Go 里做“可插拔逻辑”的基础。

真实工程里更常见的是：

```go
func Retry(times int, fn func() error) error {
	var err error
	for i := 0; i < times; i++ {
		if err = fn(); err == nil {
			return nil
		}
	}
	return err
}
```

调用方把“具体做什么”传进来，`Retry` 只负责重试流程。

## 闭包：带状态的函数

代码里的闭包：

```go
func getNumber() func() int {
	i := 0
	return func() int {
		i += 10
		return i
	}
}
```

`i` 本来是 `getNumber` 的局部变量。函数返回后它应该消失，但返回的匿名函数捕获了它，所以它继续存活。

闭包的本质：**函数 + 被捕获的外部变量。**

通俗理解：普通函数像一次性柜台，办完就清场；闭包像带抽屉的柜台，下一次来还能看到上次放进去的东西。

工程里闭包常见场景：

```go
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// before
		if r.Header.Get("token") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

`next` 被内部函数捕获，这就是 HTTP 中间件的基本形态。

## 闭包常见坑

闭包捕获的是变量，不是当时的值。旧版本 Go 里循环变量尤其容易出问题。现在 Go 1.22 起 range 变量语义已经改进，但你仍然要理解这个坑。

安全写法：

```go
for _, item := range items {
	item := item
	go func() {
		process(item)
	}()
}
```

或者：

```go
for _, item := range items {
	go func(v Item) {
		process(v)
	}(item)
}
```

面试官问这个，不是要你背版本变化，而是看你是否理解“闭包捕获变量”。

## defer：函数退出前的清理钩子

这个目录代码没写 defer，但函数章节必须补。

```go
f, err := os.Open(path)
if err != nil {
	return err
}
defer f.Close()
```

defer 在当前函数返回前执行。多个 defer 后进先出。

```go
defer fmt.Println("1")
defer fmt.Println("2")
// 输出 2 再输出 1
```

工程上，defer 用来保证资源释放：

- 文件关闭
- 数据库 rows 关闭
- 锁释放
- panic recover
- tracing span 结束

锁的典型写法：

```go
mu.Lock()
defer mu.Unlock()
```

简洁，但在超高频循环里要注意 defer 成本和释放时机。

## error 返回模式

Go 后端函数最常见签名：

```go
func DoSomething(...) error
func FindUser(...) (*User, error)
```

调用：

```go
user, err := repo.FindUser(id)
if err != nil {
	return nil, err
}
```

这不是啰嗦，是 Go 的显式控制流。优秀工程代码会让 error 带上下文：

```go
return nil, fmt.Errorf("find user %d: %w", id, err)
```

`%w` 用于错误包装，后续可以用 `errors.Is`、`errors.As` 判断。

## 后端实践怎么用

**事务封装：**

```go
func WithTx(db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}
```

调用方只写业务，事务函数负责 begin/commit/rollback。

**中间件：**

函数接收函数，返回函数，本质就是闭包。

**测试替身：**

简单依赖可以直接传函数，不一定非要上复杂 mock。

```go
type FindUserFunc func(id int64) (*User, error)
```

## 本目录代码可改进点

- `fc` 命名太随意，工程里用 `IntFunc`、`HandlerFunc`、`Mapper` 更清楚。
- `cb` 太抽象，应该表达动作。
- `nxetNumber` 拼错，应为 `nextNumber`。
- 建议补充 `error`、`defer`、多返回值示例。
- `swapByPtr` 是学习指针好例子，但真实业务里更常见的是返回新值，而不是修改入参。

## 面试拷打

1. **Go 是值传递还是引用传递？**  
   只有值传递。指针、slice、map 传进去也是复制值。

2. **slice 传函数后修改元素为什么外部可见？**  
   因为 slice header 被复制，但底层数组共享。

3. **闭包为什么能记住变量？**  
   函数捕获了外部变量，变量会逃逸并继续存活。

4. **defer 执行顺序？**  
   后进先出，函数返回前执行。

5. **为什么 error 要显式返回？**  
   Go 让失败路径可见，避免异常式隐藏控制流。

6. **什么时候用闭包，什么时候用 struct？**  
   小状态、小逻辑用闭包；状态多、方法多、生命周期复杂，用 struct。
