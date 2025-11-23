package lab1

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	// 限制只能有 3 个并发 Worker
	pool := NewWorkerPool(3)
	var counter int32
	
	taskCount := 100

	// 提交 100 个任务
	for i := 0; i < taskCount; i++ {
		pool.Submit(func() {
			time.Sleep(10 * time.Millisecond) // 模拟耗时
			atomic.AddInt32(&counter, 1)
		})
	}

	// 优雅关闭
	pool.Shutdown()

	if atomic.LoadInt32(&counter) != int32(taskCount) {
		t.Fatalf("期望执行 %d 个任务，实际执行 %d", taskCount, counter)
	}
}