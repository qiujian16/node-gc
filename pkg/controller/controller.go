package controller

import (
	"context"
	"fmt"
	"time"

	dockerapi "github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

type GCController struct {
	client *docker.Client
	period time.Duration
}

func NewGCController(period time.Duration) (*GCController, error) {
	client, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &GCController{
		client: client,
		period: period,
	}, nil
}

func (gc *GCController) Run() {
	for range time.Tick(gc.period) {
		if err := gc.cleanImages(); err != nil {
			fmt.Printf(err.Error())
		}
		if err := gc.cleanContainers(); err != nil {
			fmt.Printf(err.Error())
		}
	}
}

func (gc *GCController) cleanContainers() error {
	ctx := context.Background()
	opts := dockerapi.ContainerListOptions{
		All: true,
	}

	containers, err := gc.client.ContainerList(ctx, opts)
	if err != nil {
		return err
	}

	for _, container := range containers {
		if container.State == "error" || container.State == "exited" {
			rmOpts := dockerapi.ContainerRemoveOptions{
				Force: true,
			}
                        fmt.Printf("Remove container %s", container.ID)
			err := gc.client.ContainerRemove(ctx, container.ID, rmOpts)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (gc *GCController) cleanImages() error {
	ctx := context.Background()
	opts := dockerapi.ImageListOptions{
		All: true,
	}

	images, err := gc.client.ImageList(ctx, opts)
	if err != nil {
		return err
	}

	for _, image := range images {
                if len(image.RepoTags) == 0 {
                        rmOpts := dockerapi.ImageRemoveOptions{
                                Force: true,
                        }
                        fmt.Printf("Remove image %s", image.ID)
                        _, err := gc.client.ImageRemove(ctx, image.ID, rmOpts)
                        if err != nil {
                                return err
                        }
                }
	}
	return nil
}
