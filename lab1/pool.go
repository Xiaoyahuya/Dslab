package lab1

import (
	"sync"
)

// Task 是一个函数类型，代表要执行的任务
type Task func()

type WorkerPool struct {
	mu       sync.Mutex
	cond     *sync.Cond    // 条件变量，用于通知 worker 有新任务
	tasks    []Task        // 任务队列
	capacity int           // 最大 worker 数量
	running  int           // 当前正在工作的 worker 数量
	shutdown bool          // 是否正在关闭
	wg       sync.WaitGroup // 等待所有任务完成
}

func NewWorkerPool(capacity int) *WorkerPool {
	p := &WorkerPool{
		tasks:    make([]Task, 0),
		capacity: capacity,
	}
	// sync.Cond 需要绑定一个锁
	p.cond = sync.NewCond(&p.mu)
	return p
}

// Submit 提交一个任务到队列
func (p *WorkerPool) Submit(t Task) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.shutdown {
		return // 如果已关闭，不再接受新任务
	}

	// TODO: 1. 将任务 t 加入 p.tasks
	
	// TODO: 2. 唤醒一个正在等待的 Worker (Signal)
	// 提示：如果有空闲 Worker 在 Wait()，叫醒它；
	// 如果当前 worker 数 < capacity，我们需要启动一个新的 worker 协程
	if p.running < p.capacity {
		p.running++
		p.wg.Add(1)
		go p.workerLoop()
	} else {
		p.cond.Signal()
	}
}

// workerLoop 是 Worker 的工作循环
func (p *WorkerPool) workerLoop() {
	defer p.wg.Done()

	for {
		p.mu.Lock()

		// TODO: 3. 循环等待任务
		// 如果队列为空 且 没有 shutdown，则调用 p.cond.Wait()
		// 注意：Wait() 会自动释放锁，唤醒后重新加锁
		for len(p.tasks) == 0 && !p.shutdown {
			p.cond.Wait()
		}

		// TODO: 4. 处理退出条件
		// 如果队列为空 且 p.shutdown 为 true，说明所有活干完了，退出循环
		if len(p.tasks) == 0 && p.shutdown {
			p.mu.Unlock()
			return
		}

		// TODO: 5. 取出任务
		// task := p.tasks[0]
		// p.tasks = p.tasks[1:]
		
		p.mu.Unlock()

		// 执行任务 (不要持有锁执行任务！)
		// task()
	}
}

// Shutdown 优雅关闭：不再接受新任务，但执行完已有的任务
func (p *WorkerPool) Shutdown() {
	p.mu.Lock()
	p.shutdown = true
	p.cond.Broadcast() // 唤醒所有睡眠的 worker 让它们检测退出条件
	p.mu.Unlock()

	p.wg.Wait() // 等待所有 worker 退出
}