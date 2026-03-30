package faults

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Queens-Rocket-Engineering-Team/ctl-havoc/docker"
	"github.com/docker/docker/client"
)

type FaultDockerKill struct {
	docker      *client.Client
	ProjectName string
	ServiceName string
	killed      bool
}

func NewDockerKill(projectName, serviceName string) (*FaultDockerKill, error) {
	client, err := docker.NewClient()

	if err != nil {
		return nil, fmt.Errorf("Failed to create Docker client: %v", err)
	}

	return &FaultDockerKill{
		docker:      client,
		ProjectName: projectName,
		ServiceName: serviceName,
		killed:      false,
	}, nil
}

func (f *FaultDockerKill) Inject(ctx context.Context) error {
	// Find container ID by service name (ServiceName)
	containerID, err := docker.FindSingleContainerID(f.docker, ctx, f.ProjectName, f.ServiceName, "running")

	if err != nil {
		return fmt.Errorf("Failed to find container for service %s in project %s: %v", f.ServiceName, f.ProjectName, err)
	}

	slog.Info("Killing container with SIGKILL", "serviceName", f.ServiceName, "containerID", containerID)
	if err := f.docker.ContainerKill(ctx, containerID, "SIGKILL"); err != nil {
		return fmt.Errorf("Failed to kill container %s: %v", containerID, err)
	}
	f.killed = true

	return nil
}

func (f *FaultDockerKill) Cleanup(ctx context.Context) error {
	// No cleanup needed for DockerKill since the container is already killed
	return nil
}

func (f *FaultDockerKill) Name() string {
	return fmt.Sprintf("docker-kill/%s/%s", f.ProjectName, f.ServiceName)
}

func (f *FaultDockerKill) Duration() time.Duration {
	return time.Duration(0) // No duration since the container is killed immediately
}

func (f *FaultDockerKill) Close() error {
	return f.docker.Close()
}
