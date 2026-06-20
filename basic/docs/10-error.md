# Go error

## 这一章要记住什么

这一章主要讲四个点：

- Go 用 `error` 返回值表达错误。
- `error` 本质上是一个接口。
- `errors.New` 每次都会创建新的错误对象。
- 自定义错误可以携带业务字段，比如 code 和 msg。

---

## 1. error 基本使用

代码里：

```go
func getPositiveSelfAdd(num int) (int, error) {
	if num <= 0 {
		return -1, fmt.Errorf("num is not a positive number")
	}
	return num + 1, nil
}
```

返回值是：

```text
int, error
```

正常时：

```go
return num + 1, nil
```

出错时：

```go
return -1, fmt.Errorf("num is not a positive number")
```

```text
调用函数
   |
   v
返回 result, err
   |
   +-- err == nil：正常
   |
   +-- err != nil：出错
```

## 总结一下

Go 的错误处理是显式的。

函数把错误作为返回值交出去，调用方用 `if err != nil` 判断。

---

## 2. error 是接口

Go 内置的 `error` 可以理解成：

```go
type error interface {
	Error() string
}
```

只要一个类型实现了 `Error() string` 方法，它就可以当作 error 使用。

```text
error 接口
+----------------+
| Error() string |
+----------------+
        |
        v
实现这个方法的类型都可以当 error
```

## 总结一下

`error` 不是特殊魔法类型，它就是一个只有一个方法的接口。

---

## 3. errors.New 的比较

代码里：

```go
err3 := errors.New("hello")
err4 := errors.New("hello")
fmt.Println(err3 == err4)
fmt.Println(err3.Error() == err4.Error())
```

`errors.New("hello")` 每次都会创建一个新的错误对象。

所以：

```text
err3 == err4              -> false
err3.Error()==err4.Error() -> true
```

可以理解成：

```text
err3
+----------------+
| *errorString A |
| "hello"        |
+----------------+

err4
+----------------+
| *errorString B |
| "hello"        |
+----------------+

内容一样，但不是同一个对象
```

## 总结一下

直接比较两个 `errors.New` 产生的 error，比较的是错误对象是否相等，不是只比较错误字符串。

如果只是比较文本，可以用 `Error()` 字符串，但业务代码里更推荐用哨兵错误或 `errors.Is` 这类方式。

---

## 4. 自定义 error

代码里：

```go
type MyError struct {
	code int
	msg  string
}

func (m MyError) Error() string {
	return fmt.Sprintf("code=%d, msg=%v", m.code, m.msg)
}
```

`MyError` 实现了 `Error() string`，所以它实现了 `error` 接口。

创建错误：

```go
func NewError(code int, msg string) error {
	return MyError{
		code: code,
		msg:  msg,
	}
}
```

取出业务字段：

```go
if e, ok := err.(MyError); ok {
	return e.code
}
```

```text
error 接口变量
+----------------+
| 动态类型 MyError|
| code / msg     |
+----------------+
        |
        v
类型断言成 MyError
        |
        v
读取 code 和 msg
```

## 总结一下

自定义错误适合携带更多业务信息。

如果只返回字符串，调用方只能看到文本；如果返回结构体错误，调用方可以取出 code、msg 等字段。

---

## 易错点

1. `error` 是接口，不是结构体。
2. `err != nil` 是 Go 错误处理的常见入口。
3. `errors.New("hello")` 调两次，得到的是两个不同错误对象。
4. 类型断言要带 `ok`，避免断言失败导致 `panic`。
5. 当前代码里 `getPositiveSelfAdd(1)` 调了两次，第二次不会触发错误分支；如果想看错误，应传 `0` 或负数。

---

## 快问快答

### Q1：Go 里的 error 是什么？

答：

`error` 是一个接口，只要求实现 `Error() string` 方法。

### Q2：为什么 Go 要把 error 放到返回值里？

答：

这样调用方必须显式处理错误，错误路径更清楚。

### Q3：为什么两个 `errors.New("hello")` 不相等？

答：

因为每次都会创建新的错误对象，虽然错误文本一样，但对象不是同一个。

### Q4：自定义 error 有什么用？

答：

可以携带业务字段，比如错误码、错误信息、模块名，调用方可以通过类型断言取出来。

---

## 一句话总结

Go 的错误就是实现了 `Error() string` 的值；简单错误用标准库，自定义错误用结构体承载更多业务信息。

