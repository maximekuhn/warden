package services

import (
	"context"
	"errors"

	"github.com/docker/docker/client"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

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

func (d *DockerContainerManagementService) StartMinecraftServer(
	ctx context.Context,
	uow transaction.UnitOfWork,
	serverID valueobjects.MinecraftServerID,
) error {
	return errors.New("not yet implemented")
}
