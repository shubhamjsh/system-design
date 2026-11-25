package task

import (
	"container/heap"
	"fmt"
	"sync"
	"system-design/lld/task-scheduler/config"
	"time"
)

type Scheduler interface {
	Schedule(task *Task)
}

var (
	schedulerOnce sync.Once
	taskScheduler *TaskScheduler
)

type TaskScheduler struct {
	mtx       *sync.Mutex
	cond      *sync.Cond
	pool      *config.WorkerPool
	taskQueue *TaskHeap
}

func NewTaskScheduler(pool *config.WorkerPool) *TaskScheduler {
	schedulerOnce.Do(func() {
		taskScheduler = &TaskScheduler{
			pool:      pool,
			taskQueue: NewTaskHeap(),
			mtx:       &sync.Mutex{},
		}
		taskScheduler.cond = sync.NewCond(taskScheduler.mtx)
		go taskScheduler.Execute()
	})
	return taskScheduler
}

func (ts *TaskScheduler) Schedule(task *Task) {
	ts.mtx.Lock()
	defer ts.mtx.Unlock()

	fmt.Printf("Task: %v Scheduled\n", task.Name)

	ts.taskQueue.Push(task)
	ts.cond.Signal()
}

func (ts *TaskScheduler) Stop(task *Task) {
	ts.mtx.Lock()
	defer ts.mtx.Unlock()

	for i, item := range *ts.taskQueue {
		if task == item {
			fmt.Printf("Task: %v Stopped\n", task.Name)
			heap.Remove(ts.taskQueue, i)
			break
		}
	}
}

func (ts *TaskScheduler) Execute() {
	for {
		ts.mtx.Lock()
		if ts.taskQueue.Len() == 0 {
			ts.cond.Wait()
		}

		task := heap.Pop(ts.taskQueue).(*Task)
		if task.ExecutionTime.After(time.Now()) {
			heap.Push(ts.taskQueue, task)

			ts.mtx.Unlock()
			time.Sleep(time.Until(task.ExecutionTime))
			continue
		}
		ts.mtx.Unlock()
		ts.pool.Add(task)
	}
}
