package agent

import (
	"context"
	"sync"
)

type WorkerPool struct {
	workerCount int
	taskQueue   chan func()
	workerWg    sync.WaitGroup
}

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		taskQueue:   make(chan func()),
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	wp.workerWg.Add(wp.workerCount)
	for i := 0; i < wp.workerCount; i++ {
		go func() {
			defer wp.workerWg.Done()
			for {
				select {
				case task := <-wp.taskQueue:
					task()
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
	wp.workerWg.Wait()
}

func (wp *WorkerPool) SubmitTask(task func()) {
	wp.taskQueue <- task
}
