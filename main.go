package main

import (
	"fmt"
	"os"
	"time"

	"github.com/amirintech/cube_orchestrator/manager"
	"github.com/amirintech/cube_orchestrator/node"
	"github.com/amirintech/cube_orchestrator/task"
	"github.com/amirintech/cube_orchestrator/worker"
	"github.com/docker/docker/client"
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

	fmt.Printf("create a test container\n")
	dockerTask, createResult := createContainer()
	if createResult.Error != nil {
		fmt.Printf("%v", createResult.Error)
		os.Exit(1)
	}
	time.Sleep(time.Second * 5)
	fmt.Printf("stopping container %s\n", createResult.ContainerId)
	_ = stopContainer(dockerTask, createResult.ContainerId)
}

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres",
		Env: []string{
			"POSTGRES_USER=sugar",
			"POSTGRES_PASSWORD=cookies_and_brownies",
		},
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c,
	}

	result := d.Run()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, nil
	}

	fmt.Printf(
		"Container %s is running with config %v\n", result.ContainerId, c)
	return &d, &result
}

func stopContainer(d *task.Docker, id string) *task.DockerResult {
	result := d.Stop(id)
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil
	}

	fmt.Printf(
		"Container %s has been stopped and removed\n", result.ContainerId)
	return &result
}
