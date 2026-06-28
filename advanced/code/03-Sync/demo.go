package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// sync 全家桶，一个文件跑完所有并发原语
var divider = strings.Repeat("-", 60)

func pause() {
	time.Sleep(1500 * time.Millisecond)
}

// chan struct{} 是零成本通知
// 空结构体不占用内存，专门用于发送信号：“我完事了”、“你可以继续了”

func demo0_ChannelSignal() {
	fmt.Println(divider)
	fmt.Println("Demo 0: chan struct{} - 用 channel 发信号")
	fmt.Println(divider)

	done := make(chan struct{})	// struct{} 不占用内存，纯信号

	go func() {
		fmt.Println("	🐢	干活中...")
		time.Sleep(500 * time.Millisecond)
		fmt.Println("	✔	干完活了，发信号！")
		done <-struct{}{}	// 往channel中塞一个空结构体，意思是搞定了
	}()

	fmt.Println("	⏳	主 goroutine 等待...")
	<-done	// 收到信号才继续
	fmt.Println("   🎉 收到信号，继续执行！")
}

// 1. sync.WaitGroup - 等到人齐再触发
//	场景：旅行团等人齐
//  Add(n)	-> 登记 n 个人到
//  Done()	-> 一个人到了
// 	Wait()	-> 一直等到所有人到齐
func demo1_WaitGroup() {
	fmt.Println(divider)
	fmt.Println("Demo 1: WaitGroup - 等所有人到齐")
	fmt.Println(divider)

	var wg sync.WaitGroup
	const travelers = 5

	for i := 1; i <= travelers; i++ {
		wg.Add(1)	// 登记：又来了一个人要等
		go func(id int) {
			defer wg.Done()	// 不管怎么退出，最后必须说明“我到了！”
			delay := time.Duration(300 + rand.Intn(700)) * time.Millisecond
			time.Sleep(delay)
			fmt.Printf("	👴 旅客 %d 到了 （路上花了 %v）\n", id, delay)
		}(i)
	}
	fmt.Println("   ⏳ 导游: 大家先别走，等人齐...")
	wg.Wait() // 阻塞，等所有人 Done
	fmt.Println("   🚌 导游: 人齐了，出发！")
}

// 2. sync.Mutex - 厕所门锁
//	场景：1000 个人同时 count++，没锁就会有数据竞争
//  Lock()	-> 进门加锁
// 	Unlock()-> 解锁出门
func demo2_Mutex() {
	fmt.Println(divider)
	fmt.Println("🔒 Demo 2: Mutex — 互斥锁")
	fmt.Println(divider)

	var (
		counter int
		mu		sync.Mutex
		wg 		sync.WaitGroup
	)

	const workers = 1000

	// ❌ 如果把 mu.Lock() / mu.Unlock() 注释掉，counter 大概率 != 1000
	// 因为 counter++ 不是原子操作（读-加-写 三步）

	for i := 0; i < workers; i++ {
		wg.Add(1) 
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("   ✅ 预期: %d,  实际: %d  (相等说明锁生效)\n", workers, counter)
}

// 3. sync.RWMutex - 图书馆
//	可以有很多人同时看书（读锁），但是只能有一个人写书（写锁）
//	有人在写的时候，不能看；有人在看的时候，不能写
func demo3_RWMutex() {
	fmt.Println(divider)
	fmt.Println("📚 Demo 3: RWMutex — 读写锁(图书馆模式)")
	fmt.Println(divider)

	var (
		book string = "《C++ 从入门到入土》 第 1 版"
		rw	 sync.RWMutex
		wg   sync.WaitGroup
	)

	// 读者：可以很多人同时读
	reader := func(id int) {
		defer wg.Done()
		rw.RLock()
		fmt.Printf("	👀 读者 %d 正在看: %s\n", id, book)
		time.Sleep(300 * time.Millisecond)
		rw.RUnlock()
	}

	// 写者：独占，写的时候别人不能读也不能写
	writer := func(id int) {
		defer wg.Done()
		rw.Lock()
		book = fmt.Sprintf("《C++ 从入门到入土》第 %d 版,", id + 1)
		fmt.Printf("   ✍️  作者 %d 写完了新版本: %s\n", id, book)
		time.Sleep(500 * time.Millisecond)
		rw.Unlock()
	}

	// 先来五个读者
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go reader(i)
	}
	time.Sleep(100 * time.Millisecond) // 等读者先拿到读锁

	// 再来两个写者
	wg.Add(2)
	go writer(1)
	go writer(2)

	wg.Wait()
	fmt.Printf("   📖 最终版本: %s\n", book)
}

// 4. sync.Once - 只做一次
//	场景：全局配置初始化、数据库连接池创建
//  保证里面的函数无论被多少个 goroutine 调用，只执行一次
func demo4_Once() {
	fmt.Println(divider)
	fmt.Println("1️⃣  Demo 4: Once — 只初始化一次")
	fmt.Println(divider)

	var (
		once	sync.Once
		config	string
	)

	loadConfig := func() {
		time.Sleep(600 * time.Millisecond)
		config = "⚙️  配置加载完毕 (数据库地址、Redis地址...)"
	}

	// 模拟 10 个 goroutine 同时想要初始化配置
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("	goroutine %d: 我需要配置！\n", id)
			once.Do(loadConfig)
			fmt.Printf("	goroutine %d: 拿到配置 -> %s\n", id, config)
		}(i)
	}

	wg.Wait()
}

// 5. sync.Map - 并发安全的map
//	普通 map 并发写会 panic: concurrent map writes
//	sync.Map 专门为两种场景优化：
//  a. key 只写一次但读很多次
//	b. 多个 goroutine 读写不同的 key
// 	方法：Store Load LoadOrStore Delete Range

func demo5_SyncMap() {
	fmt.Println(divider)
	fmt.Println("🗺️  Demo 5: sync.Map — 并发安全的 Map")
	fmt.Println(divider)

	var sm sync.Map

	// 1. 存 Store
	fmt.Println(" Store 写入")
	sm.Store("🌏", "Earth")
	sm.Store("🌙", "Moon")
	sm.Store("☀️", "Sun")
	fmt.Println("      写入了 3 个星球 🌍 🌙 ☀️")

	// 2. 取 Load
	fmt.Println(" Load 读取")
	if v, ok := sm.Load("🌏"); ok {
		fmt.Printf("	🌏 = %v\n", v)
	} 
	if v, ok := sm.Load("🪐"); ok {
		fmt.Printf("      🪐 = %v\n", v)
	} else {
		fmt.Println("      🪐 不存在！")
	}

	// 有就取，没有就存 LoadOrStore
	fmt.Println(" LoadOrStore 读取或写入") 
	v, loaded := sm.LoadOrStore("🌙", "Moon2")
	fmt.Printf("      🌙: val=%v, 之前存在? %v (因为是已存在的，不覆盖)\n", v, loaded)

	v, loaded = sm.LoadOrStore("🪐", "Saturn")
	fmt.Printf("      🪐: val=%v, 之前存在? %v (不存在，存进去了)\n", v, loaded)

	// 遍历 Range 
	fmt.Println("   Range 遍历")
	sm.Range(func(key, value any) bool {
		fmt.Printf("      %v → %v\n", key, value)
		return true
	})

	// 删除 Delete 
	fmt.Println("   Delete 删除")
	sm.Delete("🌙")
		sm.Range(func(key, value any) bool {
		fmt.Printf("      %v → %v\n", key, value)
		return true
	})
}

// 6. sync.Cond 条件变量，等信号再干活
//	场景：运动员等裁判发令
//	Wait()		-> 等着，先解锁，被唤醒后重新加锁
// 	Signal()	-> 叫醒一个等着的人
//  Broadcast()	-> 广播，叫醒所有人
//  注意：Wait() 必须放在 for 循环里，不能用 if，防止虚假唤醒
func demo6_Cond() {
	fmt.Println(divider)
	fmt.Println("🏃 Demo 6: Cond — 条件变量 (发令枪模式)")
	fmt.Println(divider)

	var mu sync.Mutex
	ready := false
	cond := sync.NewCond(&mu)

	const runners = 5
	for i := 1; i <= runners; i++ {
		go func(id int) {
			fmt.Printf("   🏃 选手 %d: 就位，等待发令枪...\n", id)
			cond.L.Lock()	// 进门
			for !ready {
				cond.Wait()	// 解锁 + 睡觉 + 被叫醒
			}
			cond.L.Unlock()	// 出门

			fmt.Printf("   🏃 选手 %d: 冲啊！！！\n", id)
		}(i)
	} 

	time.Sleep(time.Second)
	fmt.Println("   📢 裁判: 各就位——预备——")
	
	cond.L.Lock()
	ready = true
	cond.L.Unlock()

	cond.Broadcast()	
	fmt.Println("   🔫 砰！")
}

// 7. sync.Pool - 对象池，用完还回去
//	场景：频繁创建/销毁的对象，例如 buffer，用 Pool 复用
//	Get() 借一个（池子空了就 New 一个）
//  Put() 还回去
//  注意：Pool 里的对象随时可能被 GC 清掉，不要依赖池子里一定有东西

func demo7_Pool() {
	fmt.Println(divider)
	fmt.Println("🏊 Demo 7: Pool — 对象池复用")
	fmt.Println(divider)

	type Bullet struct {
		ID  	int
		Used 	bool
	}

	var bulletID int32

	pool := sync.Pool {
		New: func() any {
			id := atomic.AddInt32(&bulletID, 1)
			fmt.Printf("      🏭 工厂新建设备 #%d\n", id)
			return &Bullet {
				ID:	int(id),
			}
		},
	}

	// 模拟：借 3 个，还 2 个，再借 3 个（其中有复用的）
	fmt.Println("	--- 第一轮：借 3 个设备 ---")
	b1 := pool.Get().(*Bullet)
	_ = pool.Get().(*Bullet) // b2 暂不使用，丢弃
	b3 := pool.Get().(*Bullet)

	fmt.Println("	--- 还回: 1 和 3---")
	pool.Put(b1)
	pool.Put(b3)

	fmt.Println("	--- 第二轮：再借 3 个设备 ---")
	b4 := pool.Get().(*Bullet)	// 可能拿到还回去的 1 或 3
	b5 := pool.Get().(*Bullet)	// 可能拿到剩下还回去的
	b6 := pool.Get().(*Bullet)	// 池子空了，新建

	fmt.Printf("   借到: #%d, #%d, #%d\n", b4.ID, b5.ID, b6.ID)

	fmt.Println("   💡 第4/5个设备复用了第1/3个的内存，不用重新分配")

}

// 8. sync/atomic - CPU 级别的原子操作
//	比 Mutex 更轻量，适合简单的计数器/状态切换
// 	Add Load Store Swap CompareAndSwap
func demo8_Atomic() {
	fmt.Println(divider)
	fmt.Println("⚛️  Demo 8: atomic — 无锁原子操作")
	fmt.Println(divider)

	var (
		counter		int64
		wg 			sync.WaitGroup
	)

	const workers = 1000

	// Add 原子加
	fmt.Println("   atomic.AddInt64 — 并发安全地加")
	for i := 0; i < workers; i++ {
		wg.Add(1) 
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()
	fmt.Printf("      counter = %d (预期 %d)\n", counter, workers)

	// Store Load 原子存取
	fmt.Println("    atomic.Store / Load — 原子读写")
	var flag int32
	atomic.StoreInt32(&flag, 42)
	val := atomic.LoadInt32(&flag)
	fmt.Printf("      flag = %d\n", val)

	// Swap 原子交换（返回旧值）
	fmt.Println("   atomic.Swap — 原子交换")
	old := atomic.SwapInt32(&flag, 100)
	fmt.Printf("      旧值=%d, 新值=%d\n", old, atomic.LoadInt32(&flag))
	
	// CompareAndSwap CAS 如果还是旧值，就换成新值
	fmt.Println("   atomic.CompareAndSwap — CAS 乐观锁")
	swapped := atomic.CompareAndSwapInt32(&flag, 100, 200)
	fmt.Printf("      CAS(flag, 100, 200) → %v (换成功了，flag 现在是 %d)\n", swapped, flag)

	swapped = atomic.CompareAndSwapInt32(&flag, 999, 300) // flag 不是 999，所以失败
	fmt.Printf("      CAS(flag, 999, 300) → %v (换失败了，flag 还是 %d)\n", swapped, flag)

	// atomic.Value 原子地存/取任意类型
	fmt.Println("   atomic.Value — 原子存取任意类型")
	type Config struct{ DB string }
	var cfg atomic.Value
	cfg.Store(&Config{
		DB: "mysql:://localhost",
	})

	c := cfg.Load().(*Config)
	fmt.Printf("      DB 地址: %s\n", c.DB)
}

// 综合场景：限量抢购
//		Once		-> 初始化库存
//		WaitGroup 	-> 等所有用户抢完
//		atomic 		-> 无锁扣库存
// 		Pool 		-> 复用订单结构体

func demo9_FlashSale() {
	fmt.Println(divider)
	fmt.Println("🛒 Demo 9: 综合 — 限量秒杀")
	fmt.Println(divider)

	var (
		stock 		int32 = 10
		success 	int32
		wg 			sync.WaitGroup
		once 		sync.Once
	)

	const users = 100

	// 初始化逻辑只跑一次
	initSale := func() {
		fmt.Println("   🎬 秒杀活动初始化完成！")
		fmt.Printf("   📦 库存: %d 件 | 👥 参与人数: %d\n", stock, users)
	}

	for i := 1; i <= users; i++ {
		wg.Add(1)
		go func(uid int) {
			defer wg.Done()

			once.Do(initSale)

			// CAS 乐观锁扣库存
			for {
				cur := atomic.LoadInt32(&stock)
				if cur <= 0 {
					return 
				}

				if atomic.CompareAndSwapInt32(&stock, cur, cur-1) {
					atomic.AddInt32(&success, 1)
					fmt.Printf("   ✅ 用户 %d 抢到了！剩余库存: %d\n", uid, cur-1)
					return
				}
			}
		}(i)
 	}
	wg.Wait()
	fmt.Printf("\n   🏁 秒杀结束！成功: %d 人，剩余库存: %d\n", success, stock)
}

func main() {
	fmt.Println(`
╔══════════════════════════════════════════════════════╗
║     🎯 sync 包全家桶 Demo                            ║
║     每个 demo 间隔 1.5s，看清楚再继续                  ║
╚══════════════════════════════════════════════════════╝
`)
	fmt.Println("📋 菜单:")
	fmt.Println("   0. chan struct{} — 零成本信号通知")
	fmt.Println("   1. WaitGroup  — 等人到齐")
	fmt.Println("   2. Mutex      — 互斥锁（厕所门）")
	fmt.Println("   3. RWMutex    — 读写锁（图书馆）")
	fmt.Println("   4. Once       — 只做一次")
	fmt.Println("   5. sync.Map   — 并发安全 Map")
	fmt.Println("   6. Cond       — 条件变量（发令枪）")
	fmt.Println("   7. Pool       — 对象池复用")
	fmt.Println("   8. atomic     — 原子操作")
	fmt.Println("   9. 综合: 秒杀系统 🛒")
	fmt.Println()

	demo0_ChannelSignal()
	pause()

	demo1_WaitGroup()
	pause()

	demo2_Mutex()
	pause()

	demo3_RWMutex()
	pause()

	demo4_Once()
	pause()

	demo5_SyncMap()
	pause()

	demo6_Cond()
	pause()

	demo7_Pool()
	pause()

	demo8_Atomic()
	pause()

	demo9_FlashSale()

	fmt.Println("\n" + divider)
	fmt.Println("🏁 全部 demo 跑完！sync 包的核心原语都在这里了。")
	fmt.Println(divider)
}