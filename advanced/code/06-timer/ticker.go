package main

import (
	"fmt"
	"time"
)

// =================================================
// 场景：微服务心跳上报 —— 每 2 秒向注册中心发一次心跳
// 知识点：time.NewTicker、Stop
// =================================================
func demo1_Ticker_Heartbeat() {
	fmt.Println("=== demo1: 服务心跳上报（NewTicker / Stop）===")

	sendHeartbeat := func(serviceID string) {
		fmt.Printf("  💓 [%s] 心跳上报 → 注册中心\n", serviceID)
	}

	serviceID := "user-service-01"
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop() // 退出时停止 Ticker，释放内部 goroutine

	// 模拟收到关闭信号
	stopCh := make(chan struct{})

	// 10 秒后模拟服务关闭
	go func() {
		time.Sleep(10 * time.Second)
		close(stopCh)
	}()

	heartbeatCount := 0
	fmt.Printf("[%s] 启动，开始上报心跳...\n\n", serviceID)

	for {
		select {
		case <-ticker.C:
			heartbeatCount++
			sendHeartbeat(serviceID)

		case <-stopCh:
			fmt.Printf("\n🛑 [%s] 收到关闭信号，停止心跳\n", serviceID)
			fmt.Printf("   本次运行期间共发送 %d 次心跳\n", heartbeatCount)
			return
		}
	}
}

// =================================================
// 场景：一个 Ticker 驱动多个周期任务
//       - 每 1 秒：检查服务健康状态（快速检查）
//       - 每 3 秒：拉取最新配置（较重操作）
// 知识点：Ticker 驱动 + 计数器分频
// =================================================
func demo2_Ticker_MultiTask() {
	fmt.Println("\n=== demo2: 单 Ticker 驱动多周期任务 ===")

	checkHealth := func() {
		fmt.Println("  🩺 [1s] 健康检查 → 所有指标正常")
	}

	fetchConfig := func() {
		fmt.Println("  ⚙️  [3s] 拉取配置 → 配置已更新")
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	done := make(chan struct{})
	go func() {
		time.Sleep(10 * time.Second)
		close(done)
	}()

	counter := 0
	fmt.Println("监控系统启动，开始定时任务...\n")

	for {
		select {
		case <-ticker.C:
			counter++

			// 每秒都执行的轻量任务
			checkHealth()

			// 每 3 次 tick 执行一次（1 tick = 1 秒 → 3 秒一次）
			if counter%3 == 0 {
				fetchConfig()
			}

			fmt.Printf("  ── tick #%d 完成 ──\n\n", counter)

		case <-done:
			fmt.Printf("🛑 监控系统关闭，共执行 %d 轮\n", counter)
			return
		}
	}
}

// =================================================
// 场景：定期执行数据库备份，但可以随时响应紧急停止信号
// 知识点：Ticker + select + done channel 优雅退出
// =================================================
func demo3_Ticker_GracefulShutdown() {
	fmt.Println("\n=== demo3: 定时备份 + 优雅退出（Ticker + select）===")

	backupDB := func() {
		start := time.Now()
		time.Sleep(800 * time.Millisecond) // 模拟备份耗时
		fmt.Printf("  💾 数据库备份完成（耗时 %v）\n", time.Since(start).Round(time.Millisecond))
	}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// 模拟外部发来的紧急停止信号（比如收到 SIGTERM）
	emergencyStop := make(chan struct{})

	go func() {
		time.Sleep(8 * time.Second) // 8 秒后模拟收到停止信号
		fmt.Println("\n⚠️  收到 SIGTERM 信号！")
		close(emergencyStop)
	}()

	backupCount := 0
	fmt.Println("备份服务启动，每 3 秒备份一次...\n")

	for {
		select {
		case <-ticker.C:
			backupCount++
			backupDB()
			fmt.Printf("  下次备份将在 3 秒后...\n\n")

		case <-emergencyStop:
			fmt.Printf("🛑 紧急停止！当前正在进行的备份会完成...\n")
			// 注意：如果备份正在执行，这里会等它完成才返回
			// 因为 select 同一时刻只执行一个 case
			fmt.Printf("   共完成 %d 次备份，安全退出\n", backupCount)
			return
		}
	}
}

func main() {
	demo1_Ticker_Heartbeat()
	demo2_Ticker_MultiTask()
	demo3_Ticker_GracefulShutdown()
}
