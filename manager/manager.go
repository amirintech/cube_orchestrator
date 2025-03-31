package manager

import (
	"fmt"

	"github.com/amirintech/cube_orchestrator/task"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Manager struct {
	PendingTasks  queue.Queue
	TaskDB        map[string][]*task.Task
	EventDB       map[string][]*task.TaskEvent
	Workers       []string
	WorkerTaskMap map[string][]*uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

func (m *Manager) SelectWorker() {
	fmt.Println("Selecting worker...")
}

func (m *Manager) UpdateTasks() {
	fmt.Println("Updating tasks...")
}

func (m *Manager) AssignWork() {
	fmt.Println("Assign wrok...")
}
