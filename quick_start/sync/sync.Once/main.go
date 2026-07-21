package main

import (
	"fmt"
	"sync"
)

// 多个 goroutine 同时调用，只执行一次
// 初始化数据库连接、初始化全局对象
// init()：程序启动时自动执行
// sync.Once：调用once.Do() 时才执行，懒加载资源
func demo1() {
	var once sync.Once
	var config map[string]string

	loadConfig := func() {
		fmt.Println("Load configuration")
		config = map[string]string{
			"env": "dev",
		}
	}

	GetConfig := func() map[string]string {
		once.Do(loadConfig) // 执行对应函数
		return config
	}

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(GetConfig()["env"])
		}()
	}
	wg.Wait()
}

func main() {
	demo1()
}
