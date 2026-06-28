package main

import (
	"fmt"
	"time"
	"sync"
	"sync/atomic"
)

// worker: 执行任务的 goroutine
// task: 具体任务
// Pool: 池子

type Task struct {
	f func() error 	// 具体的任务逻辑
}

func NewTask(funcArg func() error) *Task {
	return &Task {
		f: funcArg,
	}
}

type Pool struct {
	RunningWorkers		int64
	Capacity			int64
	JobCh				chan *Task
	sync.Mutex
}

func NewPool(cap int64, taskNum int) *Pool {
	return &Pool {
		Capacity: cap,
		JobCh:	  make(chan *Task, taskNum),
	}
}

func (p *Pool) GetCap() int64 {
	return p.Capacity
}

func (p *Pool) incRunning() {
	atomic.AddInt64(&p.RunningWorkers, 1)
}

func (p *Pool) decRunning() {
	atomic.AddInt64(&p.RunningWorkers, -1)
}

func (p *Pool) GetRunningWorkers() int64 {
	return atomic.LoadInt64(&p.RunningWorkers)
}

func (p *Pool) run() {
	p.incRunning()
	go func() {
		defer func() {
			p.decRunning()
		}()

		for task := range p.JobCh {
			task.f()
		}
	}()
}

func (p *Pool) AddTask(task *Task) {
	p.Lock()
	defer p.Unlock()

	if p.GetRunningWorkers() < p.GetCap() {
		p.run()
	}

	p.JobCh <-task
}

func main() {
	pool := NewPool(3, 10)

	for i := 0; i < 10; i++ {
		pool.AddTask(NewTask(func() error {
			fmt.Printf("I am Task\n")
			return nil
		}))
	}

	time.Sleep(1e9)
}