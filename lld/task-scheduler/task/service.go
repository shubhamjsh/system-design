package task

import (
	"sync"
)

type TaskService struct {
	mtx           sync.Mutex
	taskMap       map[int]*Task
	TaskScheduler *TaskScheduler
}

func NewTaskSvc(taskScheduler *TaskScheduler) *TaskService {
	return &TaskService{
		TaskScheduler: taskScheduler,
		taskMap:       make(map[int]*Task),
	}
}

func (ts TaskService) CreateTask(name, creator string, retry int, tType TaskType,
	priority TaskPriority, initialDelay, recurringDelay int) *Task {
	ts.mtx.Lock()
	defer ts.mtx.Unlock()

	task := NewTask(len(ts.taskMap)+1, retry, name, creator, tType, priority, initialDelay, recurringDelay)
	ts.taskMap[task.ID] = task
	return task
}

func (ts *TaskService) StopTask(taskID int) bool {
	ts.mtx.Lock()
	defer ts.mtx.Unlock()
	task, ok := ts.taskMap[taskID]
	if !ok {
		return true
	}
	ts.TaskScheduler.Stop(task)
	return true
}
