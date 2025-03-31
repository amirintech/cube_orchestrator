package main

import (
	"fmt"
	"time"

	"github.com/amirintech/cube_orchestrator/manager"
	"github.com/amirintech/cube_orchestrator/node"
	"github.com/amirintech/cube_orchestrator/task"
	"github.com/amirintech/cube_orchestrator/worker"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	t := task.Task{
		ID:     uuid.New(),
		Name:   "task-1",
		State:  task.Pending,
		Image:  "image-1",
		Memory: 1_024,
		Desk:   1,
	}
	te := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Pending,
		Task:      &t,
		CreatedAt: time.Now(),
	}
	fmt.Printf("TASK: %+v\n", t)
	fmt.Printf("TASK Event: %+v\n", te)

	w := worker.Worker{
		Name:      "worker-1",
		Queue:     *queue.New(),
		DB:        make(map[uuid.UUID]*task.Task),
		TaskCount: 0,
	}
	fmt.Printf("WORKER: %+v\n", w)
	w.CollectStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		PendingTasks: *queue.New(),
		TaskDB:       make(map[string][]*task.Task),
		EventDB:      make(map[string][]*task.TaskEvent),
		Workers:      []string{w.Name},
	}
	fmt.Printf("MANAGER: %+v\n", m)
	m.SelectWorker()
	m.UpdateTasks()
	m.AssignWork()

	n := node.Node{
		Name:            "node-1",
		IP:              "192.168.8.1",
		Cores:           8,
		Memory:          65_536,
		MemoryAllocated: 0,
		Disk:            24,
		DiskAllocated:   0,
		Role:            "worker",
		TaskCount:       0,
	}
	fmt.Printf("NODE: %+v\n", n)
}
