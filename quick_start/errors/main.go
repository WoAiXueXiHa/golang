package main

import (
	"errors"
	"fmt"
)

var (
	// 哨兵错误（sentinel error）：用一个固定的 error 值代表一种稳定的业务错误类型。
	// 后续不要用 err.Error() 字符串判断错误，而是用 errors.Is 判断是不是这个错误。
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrDBTimeOut    = errors.New("db timeout")
)

// TaskError 是自定义错误类型，用来在底层错误之外补充业务上下文。
// 这里额外记录了：执行了什么操作、操作的是哪个任务、底层原始错误是什么。
type TaskError struct {
	Op     string
	TaskID int64
	Err    error
}

// Error 实现 error 接口。
// 注意：这里负责把错误格式化成人能读懂的文本，但程序判断错误类型不要依赖这段文本。
func (e *TaskError) Error() string {
	return fmt.Sprintf("op=%s task_id=%d: %v", e.Op, e.TaskID, e.Err)
}

// Unwrap 暴露被包装的底层错误。
// 有了 Unwrap，errors.Is / errors.As 才能顺着错误链继续往里找。
func (e *TaskError) Unwrap() error {
	return e.Err
}

func findTask(taskID int64) error {
	err := queryTaskFromDB(taskID)
	if err != nil {
		// 返回错误时，不要只返回底层错误；最好补充当前函数知道的上下文。
		// %w 会包装错误并保留错误链，所以外层还能用 errors.Is 找到 ErrTaskNotFound。
		return &TaskError{
			Op:     "find_task",
			TaskID: taskID,
			Err:    fmt.Errorf("query db failed: %w", err),
		}
	}
	return nil
}

func queryTaskFromDB(taskID int64) error {
	// 模拟数据库查询：1001 这个任务不存在。
	// 真实项目里这里可能是 SQL 查询、RPC 调用或第三方 SDK 调用。
	if taskID == 1001 {
		return ErrTaskNotFound
	}
	return nil
}

func validateTask(title string, duration int) error {
	// 一次校验可能同时出现多个错误，用切片先收集起来。
	var errs []error

	if title == "" {
		// 这里仍然用 %w 包装 ErrInvalidInput，方便调用方统一判断“是否是非法输入”。
		errs = append(errs, fmt.Errorf("title empty: %w", ErrInvalidInput))
	}

	if duration <= 0 {
		errs = append(errs, fmt.Errorf("duration must be positive: %w", ErrInvalidInput))
	}

	// errors.Join 会把多个错误合成一个错误：
	// 1. 打印时会逐行输出每个错误。
	// 2. errors.Is 会递归检查所有子错误。
	// 3. 如果 errs 为空，Join 返回 nil。
	return errors.Join(errs...)
}

func main() {
	fmt.Println("===== 1. 普通错误 =====")
	// errors.New 创建一个最基础的错误，只有错误消息，没有错误链和额外上下文。
	err := errors.New("something wrong")
	fmt.Println(err)

	fmt.Println("\n===== 2. 包装错误并保留错误链 =====")
	err = findTask(1001)
	fmt.Println("完整错误: ", err)

	fmt.Println("\n===== 3. errors.Is 判断错误类型 =====")
	// errors.Is 不要求 err == ErrTaskNotFound。
	// 只要错误链中的任意一层包装了 ErrTaskNotFound，就能识别出来。
	if errors.Is(err, ErrTaskNotFound) {
		fmt.Println("判断结果：任务不存在")
	}

	fmt.Println("\n===== 4. errors.As 提取自定义错误信息 =====")
	var taskErr *TaskError
	// errors.As 用来从错误链里找“某一种错误类型”，找到后会赋值给 taskErr。
	// 这适合读取自定义错误里的结构化字段，例如 Op、TaskID。
	if errors.As(err, &taskErr) {
		fmt.Println("操作：", taskErr.Op)
		fmt.Println("任务ID: ", taskErr.TaskID)
	}

	fmt.Println("\n===== 5. %v 不会保留错误链 =====")
	// %v 只是把底层错误格式化成字符串，错误链会断掉。
	// 所以输出文本看起来包含 task not found，但 errors.Is 已经识别不到 ErrTaskNotFound。
	err2 := fmt.Errorf("query failed: %v", ErrTaskNotFound)
	fmt.Println("完整错误:", err2)
	fmt.Println("errors.Is 能否识别:", errors.Is(err2, ErrTaskNotFound))

	fmt.Println("\n===== 6. errors.Join 合并多个错误 =====")
	err3 := validateTask("", -1)
	fmt.Println("完整错误:")
	fmt.Println(err3)

	// 即使 err3 里合并了多个错误，errors.Is 仍然能识别里面是否包含 ErrInvalidInput。
	if errors.Is(err3, ErrInvalidInput) {
		fmt.Println("判断结果: 参数非法")
	}

}
