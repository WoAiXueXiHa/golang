package main

import "fmt"

// 1. 方法的本质：带暗号的普通函数
// 在函数名前加一个接收者
// func (接收变量名 接收类型) funcName(...) ... {}

type Counter struct {
	Count int
}

// 指针接收
func (c *Counter) AddPtr() {
	c.Count++
}

// 值接收
func (c Counter) AddValue() {
	c.Count++
}

// 无论定义的方法是值还是指针
// 无论手里拿的是值变量还是指针变量
// go 编译器在编译时会自动转换成匹配的类型

// 2. 组合代替继承
// 想让一个结构体拥有另外一个结构体的字段和方法
// 直接把另一个结构体无命名地写进来，这就是组合（匿名嵌入）
// 被嵌入地结构体地所有字段和方法，都会自动提升到外层
// 如果外层重写了同名方法，就会 override
type Engine struct {
	Power int
}

func (e *Engine) Start() {
	fmt.Println("引擎启动...")
}

type Car struct {
	Name string
	Engine
}

// override
func (c *Car) Start() {
	fmt.Printf("%s 跑车正在启动...\n", c.Name)
}

// 3. 方法集
// 值类型变量T只包含值接收者，why？一份独立的数据，如果调用指针方法，意味着要取地址
// 但这个变量值不能取地址，就会出问题
// 指针类型变量*T，全能型，值和指针方法都包含

type Payer interface {
	Pay(amount int)
}

type WeChatPay struct {
	Balance int
}

// 支付要扣钱，必须使用指针
func (w *WeChatPay) Pay(amount int) {
	w.Balance -= amount
	fmt.Printf("微信支付成功，扣除 %d 元， 余额 %d 元\n", amount, w.Balance)
}

func main() {
	fmt.Println("----------- 方法基本使用 ----------")
	c := Counter{Count: 1} // 语法糖：值类型变量
	// 编译器自动把 c.AddPtr() 转换为 (&c).AddPtr()
	c.AddPtr()
	fmt.Println("指针方法执行后: ", c.Count)

	// 编译器自动把 (&c).AddValue() 转换为 c.AddValue()
	c.AddValue()
	fmt.Println("值方法执行后: ", c.Count)

	fmt.Println("----------- 组合 ----------")
	myCar := Car{
		Name:   "Audi",
		Engine: Engine{Power: 500},
	}
	// 字段提升：可以直接使用 Engine 的字段
	fmt.Println("马力：", myCar.Power)
	// override 有限调用外层自己的方法
	myCar.Start()
	// 显式调内部被隐藏的方法
	myCar.Engine.Start()

	fmt.Println("----------- 方法集 ----------")
	wxWallet := WeChatPay{Balance: 100}
	// var p1 Payer = wxWallet
	// 会报错，WeChatPay 的方法集是空的，Pay作为指针接收者，没有资格冒充 Payer 接口
	var p2 Payer = &wxWallet
	p2.Pay(20)
}
