package task

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
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

type Config struct {
	Name          string
	AttachStdin   bool
	AttachStdout  bool
	AttachStderr  bool
	ExposedPorts  nat.PortSet
	CMD           []string
	Image         string
	CPU           float64
	Memory        int64
	Disk          int64
	Env           []string
	RestartPolicy string
}

type Docker struct {
	Client *client.Client
	Config Config
}

type DockerResult struct {
	Error       error
	Action      string
	ContainerId string
	Result      string
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	// pull image
	reader, err := d.Client.ImagePull(
		ctx, d.Config.Image, image.PullOptions{})
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}

	// write image to stdout
	bytesWritten, err := io.Copy(os.Stdout, reader)
	if err != nil {
		log.Printf("Error reading image pull response: %v\n", err)
		return DockerResult{Error: err}
	}

	// define container configs
	containerConfig := &container.Config{
		Tty:          true,
		Env:          d.Config.Env,
		Image:        d.Config.Image,
		ExposedPorts: d.Config.ExposedPorts,
	}
	restartPolicy := container.RestartPolicy{
		Name: container.RestartPolicyMode(d.Config.RestartPolicy),
	}
	resources := container.Resources{}
	containerHostConfig := &container.HostConfig{
		RestartPolicy:   restartPolicy,
		Resources:       resources,
		PublishAllPorts: true,
	}

	// create container
	cont, err := d.Client.ContainerCreate(
		ctx, containerConfig, containerHostConfig, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container %s: %v\n", cont.ID, err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerStart(ctx, cont.ID, container.StartOptions{})
	if err != nil {
		log.Printf("Error starting container %s: %v\n", cont.ID, err)
		return DockerResult{Error: err}
	}

	return DockerResult{
		Error:       nil,
		Action:      "pull",
		ContainerId: cont.ID,
		Result:      strconv.FormatInt(bytesWritten, 10),
	}
}

func (d *Docker) Stop(id string) DockerResult {
	log.Printf("Attempting to stop container %v", id)
	ctx := context.Background()

	// stop container
	err := d.Client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		log.Printf("Error stopping container %s: %v\n", id, err)
		return DockerResult{Error: err}
	}

	// release resources
	err = d.Client.ContainerRemove(ctx, id, container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	})
	if err != nil {
		log.Printf("Error removing container %s: %v\n", id, err)
		return DockerResult{Error: err}
	}

	return DockerResult{
		Action: "stop",
		Result: "success",
		Error:  nil,
	}
}
