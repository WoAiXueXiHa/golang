package main

import "fmt"

// 父结构体
type People struct {
	Name string
	Age  int
}

// 指针接收者方法，绑定*People
func (this *People) GetName() string {
	return this.Name
}

// Student 匿名嵌入People，Go组合（替代继承）
type Student struct {
	ID    int
	Score float64
	People // 匿名嵌入，不是普通字段，提升字段+方法
}

// Student重写GetName，覆盖嵌入的People同名方法
func (this *Student) GetName() string {
	return this.Name
}

// 值接收者方法：绑定Student值类型
func (st Student) GetScore() float64 {
	return st.Score
}

// 指针接收者方法：绑定*Student指针类型
func (this *Student) SetScore(score float64) {
	this.Score = score
}

func main() {
	st := Student{
		ID:    1,
		Score: 99.0,
		People: People{
			Name: "zhangsan",
			Age:  21,
		},
	}
	fmt.Printf("学生st的姓名是: %s\n", st.GetName())

	// 方法调用演示
	fmt.Printf("设置前，学生st分数：%.2f\n", st.GetScore())
	st.SetScore(99.5)
	fmt.Printf("设置后，学生st分数：%.2f\n", st.GetScore())

	stu := &Student{
		ID:    2,
		Score: 97.0,
		People: People{
			Name: "lisi",
			Age:  22,
		},
	}
	fmt.Printf("学生stu的姓名是: %s\n", stu.GetName())
	fmt.Printf("学生stu分数：%.2f\n", stu.GetScore())
}