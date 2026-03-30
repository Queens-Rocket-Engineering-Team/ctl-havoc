package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func FindSingleContainerID(c *client.Client, ctx context.Context, projectName, serviceName string, status string) (string, error) {
	containers, err := c.ContainerList(ctx, container.ListOptions{Filters: filters.NewArgs(
		filters.Arg("label", fmt.Sprintf("com.docker.compose.project=%s", projectName)),
		filters.Arg("label", fmt.Sprintf("com.docker.compose.service=%s", serviceName)),
		filters.Arg("status", status),
	)})

	if err != nil {
		return "", fmt.Errorf("Failed to list containers: %v", err)
	}

	// We expect only one container to match the criteria for a given service
	switch len(containers) {
	case 0:
		return "", fmt.Errorf("No containers found for service %s in project %s with status %s", serviceName, projectName, status)
	case 1:
		return containers[0].ID, nil
	default:
		return "", fmt.Errorf("Multiple containers found for service %s in project %s with status %s; expected exactly one", serviceName, projectName, status)
	}
}
