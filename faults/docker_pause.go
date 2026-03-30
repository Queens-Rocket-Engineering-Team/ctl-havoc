package faults

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Queens-Rocket-Engineering-Team/ctl-havoc/docker"
	"github.com/docker/docker/client"
)

type FaultDockerPause struct {
	docker      *client.Client
	ProjectName string
	ServiceName string
	Timeout     time.Duration
}

func NewDockerPause(projectName, serviceName string, timeout time.Duration) (*FaultDockerPause, error) {
	client, err := docker.NewClient()

	if err != nil {
		return nil, fmt.Errorf("Failed to create Docker client: %v", err)
	}

	return &FaultDockerPause{
		docker:      client,
		ProjectName: projectName,
		ServiceName: serviceName,
		Timeout:     timeout,
	}, nil
}

func (f *FaultDockerPause) Inject(ctx context.Context) error {
	// Find container ID by service name (ServiceName)
	containerID, err := docker.FindSingleContainerID(f.docker, ctx, f.ProjectName, f.ServiceName, "running")

	if err != nil {
		return fmt.Errorf("Failed to find container for service %s in project %s: %v", f.ServiceName, f.ProjectName, err)
	}

	slog.Info("Pausing container", "serviceName", f.ServiceName, "containerID", containerID, "timeout", f.Timeout)
	if err := f.docker.ContainerPause(ctx, containerID); err != nil {
		return fmt.Errorf("Failed to pause container %s: %v", containerID, err)
	}

	return nil
}

func (f *FaultDockerPause) Cleanup(ctx context.Context) error {
	// Find container ID by service name (ServiceName)
	containerID, err := docker.FindSingleContainerID(f.docker, ctx, f.ProjectName, f.ServiceName, "paused")

	if err != nil {
		return fmt.Errorf("Failed to find container for service %s in project %s: %v", f.ServiceName, f.ProjectName, err)
	}

	slog.Info("Unpausing container", "serviceName", f.ServiceName, "containerID", containerID, "timeout", f.Timeout)
	if err := f.docker.ContainerUnpause(ctx, containerID); err != nil {
		return fmt.Errorf("Failed to unpause container %s: %v", containerID, err)
	}

	return nil
}

func (f *FaultDockerPause) Name() string {
	return fmt.Sprintf("docker_pause/%s/%s", f.ProjectName, f.ServiceName)
}

func (f *FaultDockerPause) Duration() time.Duration {
	return f.Timeout
}

func (f *FaultDockerPause) Close() error {
	return f.docker.Close()
}
