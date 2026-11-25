package config

import (
	"fmt"
	"sync"
	"system-design/lld/task-scheduler/model"
)

var (
	once sync.Once
	pool *WorkerPool
)

type WorkerPool struct {
	workers     int
	workerQueue chan model.Runnable
	wg          sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
	once.Do(func() {
		pool = &WorkerPool{
			workers:     workers,
			workerQueue: make(chan model.Runnable),
		}

		for i := range pool.workers {
			pool.wg.Add(1)
			fmt.Printf("Worker routine initialized, ID: %v\n", i)

			go func(i int) {
				defer pool.wg.Done()
				for run := range pool.workerQueue {
					run.Run()
				}
			}(i)
		}
	})
	return pool
}

func (w *WorkerPool) Add(run model.Runnable) {
	w.workerQueue <- run
}

func (w *WorkerPool) Wait() {
	w.wg.Wait()
	close(w.workerQueue)
}
