package services

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
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
	cli                 *client.Client
	portRepo            repositories.PortRepository
	persistenceHostPath string
}

func NewDockerContainerManagementService(portRepo repositories.PortRepository, persistenceHostPath string) (*DockerContainerManagementService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerContainerManagementService{
		cli:                 cli,
		portRepo:            portRepo,
		persistenceHostPath: persistenceHostPath,
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

	c, found, err := d.containerExists(ctx, serverID)
	if err != nil {
		return err
	}

	containerID := ""
	if !found {
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
			&container.HostConfig{
				PortBindings: portBindings,
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: getVolumeName(d.persistenceHostPath, serverID),
						Target: "/home/ubuntu/paper",
						BindOptions: &mount.BindOptions{
							CreateMountpoint: true,
						},
					},
				},
			},
			nil,
			nil,
			containerName,
		)
		containerID = resp.ID
		if err != nil {
			return err
		}
	} else {
		containerID = c.ID
	}

	// TODO: check status
	if err := d.cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
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
	c, found, err := d.containerExists(ctx, serverID)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("container not found")
	}
	containerID := c.ID
	if containerID == "" {
		return errors.New("container is not running")
	}
	// TODO: check container status
	return d.cli.ContainerStop(ctx, containerID, container.StopOptions{})
}

func (d *DockerContainerManagementService) containerExists(
	ctx context.Context,
	serverID valueobjects.MinecraftServerID,
) (*types.Container, bool, error) {
	// TODO: filter by name
	containers, err := d.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, false, err
	}
	// for some reason, actual name starts with a '/'
	containerName := fmt.Sprintf("/%s", getContainerName(serverID))
	for _, c := range containers {
		if slices.Contains(c.Names, containerName) {
			return &c, true, nil
		}
	}
	return nil, false, nil
}

func getContainerName(id valueobjects.MinecraftServerID) string {
	return fmt.Sprintf("warden-%s", id.Value().String())
}

func getVolumeName(hostpath string, id valueobjects.MinecraftServerID) string {
	return fmt.Sprintf(
		"%s/papermc%s",
		hostpath,
		strings.ReplaceAll(id.Value().String(), "-", ""),
	)
}
