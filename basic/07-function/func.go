package main

import (
	"errors"
	"fmt"
	"strings"
)

// 1. 函数骨架，支持多返回值---->错误处理
// go 中的错误必须显示处理 if-err 哲学
func RegisterUser(username string, password string) (string, error) {
	if len(username) < 3 {
		return "", errors.New("用户名长度不能小于3位")
	}
	if len(password) < 6 {
		return "", errors.New("密码长度不能小于6位")
	}

	welcomeMsg := username + ", 欢迎加入!"
	return welcomeMsg, nil
}

// 2. 传参的本质：值传递
// go 只有值传递，传结构体会完全拷贝一份副本，占用新内存
// 传指针复制的是地址的副本，指向同一个地址
// 游戏角色血量的修改
type Player struct {
	Name string
	HP   int
}

// 错误：值传递
func DamageByValue(p Player) {
	p.HP -= 20
	fmt.Printf("in function (value): %s 的血量被改为 %d\n", p.Name, p.HP)
}

// 正确：指针传递
func DamageByPtr(p *Player) {
	p.HP -= 20
	fmt.Printf("in function (ptr): %s 的血量被改为 %d\n", p.Name, p.HP)
}

// 3. 函数作为另外一个函数的参数
// 只要符合这个特征的函数，都可以作为参数，起别名了
type LogFilter func(msg string) bool

// 工厂函数
func CreateFilter(keyword string) LogFilter {
	// 直接返回匿名函数，它把 keyword 变量“打包”带走了
	return func(msg string) bool {
		// 无论 keyword 是什么，直接丢进 Contains 里去匹配！
		return strings.Contains(msg, keyword)
	}
}

func ProcessLogs(logs []string, f LogFilter) {
	for _, log := range logs {
		// 调用传进来的函数 f， 判断日志是否符合规则
		if f(log) {
			fmt.Println("发现目标日志 -> ", log)
		}
	}
}

// 过滤包含 ERROR 的日志
func IsErrorLog(msg string) bool {
	return msg == "ERROR"
}

func IsDebugLog(msg string) bool {
	return msg == "DEBUG"
}

func main() {
	fmt.Println("------------ 函数支持多返回值，if-err 处理 ----------------")
	msg, err := RegisterUser("Golang", "123456")
	if err != nil {
		fmt.Println("注册失败，原因：", err)
		return
	}
	// 此时 err 必定为空，安全使用msg
	fmt.Println("注册成功：", msg)

	fmt.Println("------------ 值传递和指针传递 ----------------")
	boss := Player{
		Name: "BOSS",
		HP:   100,
	}

	fmt.Println("1. 尝试值传递")
	DamageByValue(boss)
	fmt.Println("Main function: 此时BOSS的实际血量: ", boss.HP)

	fmt.Println("2. 尝试指针传递")
	DamageByPtr(&boss)
	fmt.Println("Main function: 此时BOSS的实际血量: ", boss.HP)

	fmt.Println("------------ 函数作为另外一个函数的参数 ----------------")
	allLogs := []string{
		"INFO",
		"ERROR",
		"DEBUG",
	}
	fmt.Println("开始过滤日志:")
	ProcessLogs(allLogs, IsErrorLog) // 把 IsErrorLog 这个函数名，像变量一样传递了
	ProcessLogs(allLogs, IsDebugLog)

	fmt.Println("------------ 匿名函数和闭包 ----------------")
	// 4. 匿名函数，常用匿名函数来做延迟处理和封装临时逻辑
	// 闭包
	// 一般函数： func funcName(参数列表) 返回类型列表 {}
	// 匿名函数:  func(参数列表) 返回类型列表{}
	// func(参数列表) 返回类型列表 可以用 type 关键字起别名
	discount := func(price float64) float64 {
		return price * 0.8
	}
	// discount 存了 func(price float64) float64 这个函数类型
	// 有括号就是调用这个函数，没括号就是没调用，我有它这个函数类型的意思
	fmt.Println("打折后的价格：", discount(100))

	// 立刻执行匿名函数
	func() {
		tmpMsg := "我是临时变量，执行完之后就被销毁"
		fmt.Println(tmpMsg)
	}() // 括号触发函数执行

	// 闭包：匿名函数打包了自己作用域外的变量，让这个变量的生命周期延长
	factor := 2
	multiplier := func(val int) int {
		factor++
		return val * factor // 闭包捕获了外面的变量 factor
	}
	fmt.Println("第一次调用multiplier，结果： ", multiplier(5))
	fmt.Println("第二次调用multiplier，结果： ", multiplier(5))
}
