package main

import (
	"system-design/lld/task-scheduler/config"
	"system-design/lld/task-scheduler/task"
	"time"
)

func main() {
	pool := config.NewWorkerPool(5)

	taskScheduler := task.NewTaskScheduler(pool)
	taskService := task.NewTaskSvc(taskScheduler)

	task1 := taskService.CreateTask("Pikachu Capture", "Team Rocket", 3, task.FixedRate, task.High, 60, 180)
	task2 := taskService.CreateTask("Defeat Vilgax", "Ben 10", 3, task.FixedDelay, task.Low, 30, 120)
	task3 := taskService.CreateTask("Eat Ramen", "Naruto", 3, task.OneTime, task.Low, 10, 0)

	taskScheduler.Schedule(task1)
	taskScheduler.Schedule(task2)
	taskScheduler.Schedule(task3)

	go func() {
		time.Sleep(time.Duration(40) * time.Second)
		taskService.StopTask(task2.ID)
	}()

	pool.Wait()
}
