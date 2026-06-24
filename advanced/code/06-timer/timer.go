package main

import (
	"fmt"
	"time"
)

// =================================================
// 场景：调用外部支付 API，设置 3 秒超时
// 知识点：time.NewTimer、Stop、Reset
// =================================================
func demo1_Timer_Stop_Reset() {
	fmt.Println("=== demo1: API 超时控制（NewTimer / Stop / Reset）===")

	callPaymentAPI := func(orderID string, delay time.Duration) string {
		time.Sleep(delay) // 模拟 API 调用耗时
		return "payment_ok"
	}

	processOrder := func(orderID string, apiDelay time.Duration) {
		fmt.Printf("[订单 %s] 开始处理，预估 API 耗时 %v\n", orderID, apiDelay)

		// 创建 Timer：最多等 3 秒
		timer := time.NewTimer(3 * time.Second)

		// 用 goroutine 异步调用 API，结果通过 channel 返回
		resultCh := make(chan string, 1)
		go func() {
			resultCh <- callPaymentAPI(orderID, apiDelay)
		}()

		select {
		case result := <-resultCh:
			// API 在超时前返回了 → 停掉 Timer，不会触发超时
			if !timer.Stop() {
				// Stop 返回 false 说明 Timer 已经触发了，
				// 需要排空 channel 防止 goroutine 泄漏
				<-timer.C
			}
			fmt.Printf("[订单 %s] ✅ 支付成功: %s\n", orderID, result)

		case <-timer.C:
			// 超时了！但也许还能抢救一下？
			fmt.Printf("[订单 %s] ⏰ 首次超时！尝试延长等待时间...\n", orderID)

			// Reset：重新计时，再给 2 秒
			timer.Reset(2 * time.Second)

			select {
			case result := <-resultCh:
				if !timer.Stop() {
					<-timer.C
				}
				fmt.Printf("[订单 %s] ✅ 延迟恢复，支付成功: %s\n", orderID, result)
			case <-timer.C:
				fmt.Printf("[订单 %s] ❌ 最终超时，订单取消\n", orderID)
			}
		}
		fmt.Println()
	}

	// 订单 A：API 很快（0.5 秒），Timer 被 Stop
	processOrder("A-001", 500*time.Millisecond)

	// 订单 B：API 很慢（4 秒），第一次超时后 Reset 再等，最终还是超时
	processOrder("B-002", 6*time.Second)

	// 订单 C：API 稍慢（3.5 秒），第一次超时后用 Reset 等到了结果
	processOrder("C-003", 3500*time.Millisecond)
}

// =================================================
// 场景：用户会话过期清理 —— 用户 5 秒不操作就触发清理，
//       但如果用户重新操作了就取消清理
// 知识点：time.AfterFunc —— 到点自动执行回调
// =================================================
func demo2_AfterFunc() {
	fmt.Println("=== demo2: 会话过期清理（AfterFunc）===")

	type Session struct {
		UserID     string
		LastActive time.Time
	}

	startSessionMonitor := func(session *Session) (renew func(), logout chan string) {
		logout = make(chan string, 1)

		// AfterFunc 返回的 Timer 可以 Stop —— 这就是它的关键：
		//   - 如果用户重新操作了 → Stop 掉，不清理
		//   - 如果用户真的 5 秒没操作 → 回调执行，踢下线
		var idleTimer *time.Timer

		renew = func() {
			if idleTimer != nil {
				// 用户有操作 → 停掉旧的倒计时
				stopped := idleTimer.Stop()
				if !stopped {
					<-idleTimer.C // 排空
				}
			}
			session.LastActive = time.Now()
			fmt.Printf("  [会话 %s] 🔄 用户操作刷新，重新倒计时 5 秒\n", session.UserID)

			// 重新启动 5 秒倒计时
			idleTimer = time.AfterFunc(5*time.Second, func() {
				logout <- session.UserID
				fmt.Printf("  [会话 %s] ⏰ 5 秒无操作，自动踢下线！\n", session.UserID)
			})
		}

		// 初始启动倒计时
		renew()
		return
	}

	// ——— 用户 zhangsan：一直活跃，从不被踢 ———
	session1 := &Session{UserID: "zhangsan"}
	renew1, logout1 := startSessionMonitor(session1)

	// ——— 用户 lisi：2 秒后操作一次，然后不再操作，最终被踢 ———
	session2 := &Session{UserID: "lisi"}
	renew2, logout2 := startSessionMonitor(session2)

	// 用 done channel 优雅等待
	done := make(chan struct{})
	go func() {
		for {
			select {
			case user := <-logout1:
				fmt.Printf("🚪 系统通知：%s 被强制下线\n", user)
			case user := <-logout2:
				fmt.Printf("🚪 系统通知：%s 被强制下线\n", user)
				close(done)
				return
			}
		}
	}()

	// 模拟用户行为时间线
	time.Sleep(2 * time.Second)
	renew2() // lisi 在第 2 秒操作了一次

	time.Sleep(2 * time.Second)
	renew1() // zhangsan 在第 4 秒操作

	time.Sleep(3 * time.Second)
	renew1() // zhangsan 在第 7 秒操作（保持活跃）

	<-done
	fmt.Println()
}

// =================================================
// 场景：从下游服务拉取数据，最多等 2 秒，超时走降级逻辑
// 知识点：time.After —— 本质是 NewTimer(d).C 的语法糖
// =================================================
func demo3_After() {
	fmt.Println("=== demo3: 下游服务调用超时降级（time.After）===")

	fetchFromDownstream := func(serviceName string, delay time.Duration) string {
		time.Sleep(delay)
		return fmt.Sprintf("<%s 的真实数据>", serviceName)
	}

	queryService := func(serviceName string, delay time.Duration) {
		fmt.Printf("→ 查询 %s（预估耗时 %v）...\n", serviceName, delay)

		resultCh := make(chan string, 1)
		go func() {
			resultCh <- fetchFromDownstream(serviceName, delay)
		}()

		select {
		case data := <-resultCh:
			fmt.Printf("✅ %s 返回: %s\n", serviceName, data)

		case <-time.After(2 * time.Second):
			// time.After(d) 等价于 NewTimer(d).C
			// 但它不返回 Timer 对象，没法 Stop ——
			// Timer 要等到超时才会被 GC 回收
			// 适合只用一次、不在乎这点内存的场景
			fmt.Printf("⚠️  %s 超时！使用缓存数据降级\n", serviceName)
		}
		fmt.Println()
	}

	// 服务 A：很快回来
	queryService("用户画像服务", 500*time.Millisecond)

	// 服务 B：超时
	queryService("推荐服务", 5*time.Second)

	// ——— time.After 的坑 ———
	fmt.Println("💡 注意：time.After 在循环里用会导致内存泄漏！")
	fmt.Println("   每次调用都会新建一个 Timer，select 没走到超时分支时")
	fmt.Println("   那些 Timer 要等到超时才会被 GC。")
	fmt.Println("   循环场景应该用 time.NewTimer + Reset 复用。")
}

func main() {
	demo1_Timer_Stop_Reset()
	demo2_AfterFunc()
	demo3_After()
}
