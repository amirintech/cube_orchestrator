package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Task struct {
	ID            uuid.UUID
	Name          string
	State         State
	Image         string
	Memory        int
	Desk          int
	ExposedPorts  nat.PortSet
	PortPindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	StopTime      time.Time
}

type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Task      *Task
	CreatedAt time.Time
}
