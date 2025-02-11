package services

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/maximekuhn/warden/internal/domain/repositories"
	"github.com/maximekuhn/warden/internal/domain/services"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

const (
	imageName   = "warden-mc"
	exposedPort = "25565/tcp"
)

type DockerContainerManagementService struct {
	cli      *client.Client
	portRepo repositories.PortRepository
}

func NewDockerContainerManagementService(portRepo repositories.PortRepository) (*DockerContainerManagementService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerContainerManagementService{
		cli:      cli,
		portRepo: portRepo,
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
	serverPort, found, err := d.portRepo.GetByServerID(ctx, uow, serverID)
	if err != nil {
		return err
	}
	if !found {
		return services.ErrServerNotFound
	}

	portBindings := nat.PortMap{
		nat.Port(exposedPort): []nat.PortBinding{
			{HostIP: "0.0.0.0", HostPort: fmt.Sprint(serverPort)},
		},
	}
	containerName := getContainerName(serverID)
	resp, err := d.cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        imageName,
			ExposedPorts: nat.PortSet{nat.Port(exposedPort): struct{}{}},
		},
		&container.HostConfig{PortBindings: portBindings},
		nil,
		nil,
		containerName,
	)
	if err != nil {
		return err
	}

	if err := d.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	// TODO: wait for server to actually start (check logs)

	return nil
}

func (d *DockerContainerManagementService) StopMinecraftServer(
	ctx context.Context,
	uow transaction.UnitOfWork,
	serverID valueobjects.MinecraftServerID,
) error {
	// TODO: filter by name
	containers, err := d.cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return err
	}
	containerID := ""
	// for some reason, actual name starts with a '/'
	containerName := fmt.Sprintf("/%s", getContainerName(serverID))
	for _, c := range containers {
		if slices.Contains(c.Names, containerName) {
			containerID = c.ID
			break
		}
	}
	if containerID == "" {
		return errors.New("container is not running")
	}
	// TODO: check container status
	return d.cli.ContainerStop(ctx, containerID, container.StopOptions{})
}

func getContainerName(id valueobjects.MinecraftServerID) string {
	return fmt.Sprintf("warden-%s", id.Value().String())
}
