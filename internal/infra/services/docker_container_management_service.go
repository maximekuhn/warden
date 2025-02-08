package services

import (
	"context"
	"errors"
	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

const imageName = "warden-mc"

type DockerContainerManagementService struct {
	cli *client.Client
}

func NewDockerContainerManagementService() (*DockerContainerManagementService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerContainerManagementService{
		cli: cli,
	}, nil
}

// For now, this implementation of services.ContainerManagementService does not build the
// minecraft server docker image if not present. Use this function to check if the image already
// exists.
func (d *DockerContainerManagementService) EnsureMinecraftServerImageExists() (bool, error) {
	images, err := d.cli.ImageList(context.TODO(), image.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, image := range images {
		for _, repotag := range image.RepoTags {
			imgName := strings.Split(repotag, ":")[0]
			if imageName == imgName {
				return true, nil
			}
		}
	}
	return false, nil
}

func (d *DockerContainerManagementService) StartMinecraftServer(
	ctx context.Context,
	uow transaction.UnitOfWork,
	serverID valueobjects.MinecraftServerID,
) error {
	return errors.New("not yet implemented")
}
