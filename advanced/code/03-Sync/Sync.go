package main

import (
	"fmt"
	"sync"
	"time"
	"strconv"
	"sync/atomic"
)

// =============================================================
// 1. sync.WaitGroup 使用
// =============================================================

func demo1() {
	// 空的 struct 做参数，用于通知，表示收到或者发出一个消息，没有实际意义 
	ch := make(chan struct{}, 10) 
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("num:%d\n", i)
			// struct{}是一个类型->空结构体类型 struct{}{}是这个类型的实例
			ch<- struct{}{}
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-ch
	}
	fmt.Println("end")
}

// (wg *WaitGroup) Add(delta int)	计数器加delta
// (wg *WaitGroup) Done()			计数器减1
// (wg *WaitGroup) Wait()			阻塞等待到计数器减成0

func demo1_WaitGroup() {
	var wg sync.WaitGroup

	myGo := func() {
		defer wg.Done()
		fmt.Println("myGo!")
	}

	wg.Add(10)

	for i := 0; i < 10; i++ {
		go myGo()
	}

	wg.Wait()
	fmt.Println("end!!!")
}

// =============================================================
// 2. sync.Once 使用
// 很多逻辑只需要执行一次，比如配置文件的加载
// sync.Once 可以在代码任意位置初始化和调用，线程安全
// 作用：延迟初始化，在第一次用到它的时候初始化一次，之后就留在内存中
// =============================================================
type Config struct{}	// 声明配置结构体

var instance *Config
var once sync.Once

// 获取配置结构体
func IintConfig() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

// =============================================================
// 3. sync.Lock 使用
// =============================================================
func demo2_nolock() {
	var (
		num int
		wg = sync.WaitGroup{}
	)

	add := func() {
		defer wg.Done()
		num += 1
	}

	var n = 10 * 10 * 10 * 10 * 10
	wg.Add(n)

	for i := 0; i < n; i++ {
		go add()
	}
	wg.Wait()

	fmt.Println(num == n)
}

func dmeo2_lock() {
	var (
		num int
		wg = sync.WaitGroup{}
		guard = sync.Mutex{}
	)

	add := func() {
		guard.Lock()
		defer wg.Done()
		num += 1
		guard.Unlock()
	}

	var n = 10 * 10 * 10 * 10 * 10
	wg.Add(n)

	for i := 0; i < n; i++ {
		go add()
	}
	wg.Wait()

	fmt.Println(num == n)
}

// 读写锁，将读和写分开，一般用于大量读、少量写的情况
// 1. 同时只能有一个 goroutine 能够获得写锁
// 2. 通知可以有读个 goroutine 获得读锁
// 3. 同时只能存在写锁定或读锁定->读和写互斥
func demo2_rwlock() {
	cnt := 0

	read := func(mr *sync.RWMutex, i int) {
		fmt.Printf("goroutine %d reader start\n", i)
		mr.RLock()
		fmt.Printf("goroutine %d reading count: %d\n", i, cnt)
		time.Sleep(time.Millisecond)
		fmt.Printf("goroutine %d reader over\n", i)
		mr.RUnlock()
	}

	write := func(mr *sync.RWMutex, i int) {
		fmt.Printf("goroutine %d writer start\n", i) 
		mr.Lock()
		cnt++
		fmt.Printf("goroutine %d writing count: %d\n", i, cnt)
		time.Sleep(time.Millisecond)
		fmt.Printf("goroutine %d writer over\n", i)
		mr.Unlock()
		
	}

	var mr sync.RWMutex
	for i := 1; i <= 3; i++ {
		go write(&mr, i)
	}

	for i := 1; i <= 3; i++ {
		go read(&mr, i)
	}

	time.Sleep(time.Second)
	fmt.Println("final count:", cnt)
}

// =============================================================
// 4. 死锁
// 两个及以上的 goroutine 在执行过程中，因为争夺共享资源处在相互等待的状态
// 如果没有外部干涉将会一直处于阻塞状态
// =============================================================

// 4.1. Lock/Unlock不成对
// fatal error: all goroutines are asleep - deadlock!
// goroutine 1 [sync.Mutex.Lock]:
// internal/sync.runtime_SemacquireMutex(0x490013?, 0x78?, 0x414d57?)
//         /home/vect/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/runtime/sema.go:95 +0x25
// internal/sync.(*Mutex).lockSlow(0xc000086020)
//         /home/vect/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/internal/sync/mutex.go:149 +0x15d
// internal/sync.(*Mutex).Lock(...)
//         /home/vect/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/internal/sync/mutex.go:70
// sync.(*Mutex).Lock(...)
//         /home/vect/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/src/sync/mutex.go:46
// main.demo4_no_couple_lock.func1({{}, {0x1, 0x0}})
//         /home/vect/golang/advanced/code/03-Sync/Sync.go:170 +0x5e
// main.demo4_no_couple_lock()
//         /home/vect/golang/advanced/code/03-Sync/Sync.go:178 +0x8b
// main.main()
//         /home/vect/golang/advanced/code/03-Sync/Sync.go:189 +0xf
// exit status 2
func demo4_no_couple_lock() {
	copyMutex := func(mu sync.Mutex) {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("ok")
	}

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	// 将带有锁结构的变量赋值给其他变量，锁的状态会赋值
	copyMutex(mu)
}

// 4.2 循环等待
// fatal error: all goroutines are asleep - deadlock!
func demo4_loop_wait() {
	var mu1, mu2 sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2) 
	go func() {
		defer wg.Done()
		mu1.Lock()
		defer mu1.Unlock()
		time.Sleep(1 * time.Second)

		mu2.Lock()
		defer mu2.Unlock()
	}()

	go func() {
		defer wg.Done()
		mu2.Lock()
		defer mu2.Unlock()
		time.Sleep(1 * time.Second)

		mu1.Lock()
		defer mu1.Unlock()
	}()

	wg.Wait()
}

// =============================================================
// 5. sync.Map 使用
// =============================================================
// go 内置的 Map 不是并发安全的
// fatal error: concurrent map writes
func demo5_show_map() {
	var m = make(map[string]int)

	getVal := func(key string) int {
		return m[key]
	}

	setVal := func(key string, value int) {
		m[key] = value
	}

	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(num int) {
			defer wg.Done()
			key := strconv.Itoa(num)
			setVal(key, num)
			fmt.Printf("key:%v, val:%v\n", key, getVal(key))
		}(i)
	}
	wg.Wait()
}

// 所以需要加锁
func demo5_lock_map() {
	var m = make(map[string]int)
	var mu sync.Mutex

	getVal := func(key string) int {
		return m[key]
	}

	setVal := func(key string, value int) {
		m[key] = value
	}

	wg := sync.WaitGroup{}

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(num int) {
			defer func() {
				wg.Done()
				mu.Unlock()
			}()

			key := strconv.Itoa(num)
			mu.Lock()
			setVal(key, num)
			fmt.Printf("key:%v, val:%v\n", key, getVal(key))
		}(i)
	}
	wg.Wait()
}

// 1.9 版本引入了并发安全版的sync.Map
// 没有提供获取 map 数量的方法，遍历时需要自行计算
// 非并发情况下，使用 map 效率会更高，因为要保证并发安全，一定会有性能损失
func demo5_show_Map() {
	var m sync.Map
	// 写入 Store
	m.Store("name", "zhangsan")
	m.Store("age", 18)

	// 读取 Load
	age, _ := m.Load("age")
	fmt.Println(age)

	// 遍历 Range
	m.Range(func(key, value interface{}) bool {
		fmt.Printf("key: %v, val: %v\n", key, value)
		return true
	})

	// 删除 Delete
	m.Delete("age")
	age, ok := m.Load("age")
	fmt.Println(age, ok)

	// 读取或写入 LoadOrStore
	m.LoadOrStore("name", "lisi")
	name, _ := m.Load("name")
	fmt.Println(name)
}

// =============================================================
// 6. sync/atomic 使用
// atomic和mutex的区别：
// 1、使用方式：mutex保护共享资源，atomic针对变量操作
// 2、底层实现：mutex由OS调度器实现，atomic操作底层有硬件指令支持->保证在CPU上执行不被中断
// atomic提供了这些方法:
// func AddT(addr *T, delta T)(new T)
// func StoreT(addr *T, val T)
// func LoadT(addr *T) (val T)
// func SwapT(addr *T, new T) (old T)
// func CompareAndSwapT(addr *T, old, new T) (swapped bool)
// =============================================================


func demo6() {
	var sum int32 = 0
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt32(&sum, 1)
		}()
	}
	wg.Wait()
	fmt.Printf("sum is %d\n", sum)
}

func demo6_atomic_struct() {
	type Student struct {
		Name string
		Age int
	}

	st1 := Student {
		Name: "zhangsan",
		Age:   19,
	}

	st2 := Student {
		Name: "lisi",
		Age:   20,
	}
	
	st3 := Student {
		Name: "wangwu",
		Age:   21,
	}

	var v atomic.Value
	v.Store(st1)
	fmt.Println(v.Load().(Student))

	old := v.Swap(st2)
	fmt.Printf("after swap: v=%v\n", v.Load().(Student))
	fmt.Printf("after swap: old%v\n", old)

	swapped := v.CompareAndSwap(st1, st3)
	fmt.Println("compare st1 and v\n", swapped,v)

	swapped = v.CompareAndSwap(st2, st3)
	fmt.Println("compare st2 and v\n", swapped, v)
}

// =============================================================
// 7. sync.Pool使用
// 内存池组件，实现对象复用，避免创建相同的对象
// New() 构造函数，指定缓存的类型
// Get() 取对象
// Put() 放对象
// =============================================================
func demo7_show_Pool() {
	type Student struct {
		Name		string
		Age 		int
	}

	pool := sync.Pool {
		// 初始化一个 sync.Pool 对象，返回 Student 指针
		New: func() interface{} {
			return &Student {
				Name: "zhangsan",
				Age:   18,
			}
		},
	}

	st := pool.Get().(*Student)
	println(st.Name, st.Age)
	fmt.Printf("addr is %p\n", st)

	// 修改
	st.Name = "hghhhhhh"
	st.Age = 34

	// 回收
	pool.Put(st)

	st1 := pool.Get().(*Student)
	println(st1.Name, st1.Age)
	fmt.Printf("addr1 is %p\n", st1)
}

func main() {
	// demo1()
	// demo1_WaitGroup()
	// demo2_nolock()
	// dmeo2_lock()	
	// demo2_rwlock()
	// demo4_no_couple_lock()
	// demo4_loop_wait()
	// demo5_show_map()
	// demo5_lock_map()
	// demo5_show_Map()
	// demo6()
	// demo6_atomic_struct()
	demo7_show_Pool()
}