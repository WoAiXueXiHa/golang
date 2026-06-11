# 01-go-language-struct：先把 Go 程序的骨架看懂

这个目录最值得记住一句话：**Go 代码不是按单个文件组织的，而是按 package 组织的。** `hello/main.go` 是入口，`myMath/` 是被调用的包，`Add` 和 `Mul` 分散在两个文件里，但因为同属 `package mymath`，编译时会被当成一个整体。

## 关键代码怎么看

`hello/main.go` 里这句最重要：

```go
import (
	"fmt"
	myMath "learn/golang/01-go-language-struct/myMath"
)
```

这里同时出现了三个概念，初学时很容易混：

- `learn/golang/01-go-language-struct/myMath` 是导入路径，由 `go.mod` 里的模块名加目录路径组成。
- `package mymath` 是包名，写在 `myMath1.go` 和 `myMath2.go` 顶部。
- `myMath` 是导入别名，调用时写 `myMath.Add`。

真正优秀的 Go 写法会尽量让这三个东西少制造歧义。这个目录里目录叫 `myMath`，包名叫 `mymath`，又用别名 `myMath` 修正调用体验，能跑，但不够顺。更推荐目录名也叫 `mymath`，包名也叫 `mymath`，调用时自然就是 `mymath.Add(...)`。

`myMath/myMath1.go` 和 `myMath/myMath2.go` 说明了另一个重点：

```go
package mymath

func Add(x, y int) int {
	return x + y
}
```

`Add` 和 `Mul` 都是大写开头，所以能被 `main` 包调用。Go 没有 `public/private` 关键字，**首字母大写就是导出，小写就是包内私有**。这个规则以后会影响结构体字段、接口方法、JSON 序列化、跨包调用，必须形成肌肉记忆。

## 必须掌握的点

**1. `package main` + `func main()` 才是可执行程序入口。**  
普通工具包不会写 `package main`。真实后端服务里常见入口是 `cmd/server/main.go`，里面只负责启动配置、日志、数据库、路由，真正业务逻辑放到别的包。

**2. 一个目录通常就是一个包。**  
同一目录下的 `.go` 文件应该声明同一个 package。你可以把一个包拆成多个文件，但不要把一个目录塞成多个包，这会让编译和阅读都变混乱。

**3. import 路径不是磁盘绝对路径。**  
它是模块路径加相对目录路径。这里的 `learn/golang/...` 能成立，是因为根目录的 `go.mod` 声明了模块名。

**4. 导出规则很重要。**  
`Add` 能跨包调用，`add` 不能。后端项目里如果 `User.Name` 是大写，JSON、ORM、其他包才能访问；如果写成 `name`，外部包就看不到。

## 用一个形象例子理解

把 Go 项目想成一栋办公楼：

- `go.mod` 是楼盘地址。
- 每个目录是一个部门。
- `package` 是部门门牌。
- 大写开头的函数和字段，是这个部门对外开放的窗口。
- 小写开头的是部门内部资料，外人不能直接拿。

你现在这个 `myMath` 目录像是门牌写着 `mymath`，但楼层导航叫 `myMath`。能找到，但对后来维护的人不够友好。

## 和 Go 后端开发的关系

这个 demo 是后端项目分层的最小雏形：入口包调用业务包。真实项目会扩展成：

```text
cmd/server/main.go
internal/handler
internal/service
internal/repository
internal/model
```

优秀 Go 后端工程师不是把所有代码写进 `main.go`，而是让 `main.go` 像总开关，业务包各司其职。

## 更像工程代码的写法

- 把目录 `myMath` 改成全小写 `mymath`。
- import 时尽量不使用别名，除非包名冲突或包名本身不清晰。
- 给导出函数加简短注释，例如 `// Add returns the sum of x and y.`
- 保持入口文件薄一点，复杂逻辑放到独立包里。

## 复习时问自己

1. **为什么 `Add` 能被 `main` 包调用？**  
   因为它首字母大写，是导出标识符。

2. **`package main` 有什么特殊？**  
   它表示这个包会编译成可执行程序，并且需要 `func main()` 作为入口。

3. **同一个包拆成多个文件后，文件之间要互相 import 吗？**  
   不需要。同目录同包的文件共享同一个包级命名空间。

4. **真实后端项目里，为什么不应该把所有逻辑都写进 main？**  
   因为入口应该负责组装，业务逻辑应该放在 handler/service/repository 等包里，方便测试和维护。
