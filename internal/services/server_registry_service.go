package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/entities"
	"github.com/maximekuhn/warden/internal/repositories"
	"github.com/maximekuhn/warden/internal/valueobjects"
)

// ServerRegistryService is a service to handle Minecraft server creation,
// deletion and list fetch.
type ServerRegistryService struct {
	msRepository repositories.MinecraftServerRepository
}

func NewServerRegistryService(r repositories.MinecraftServerRepository) *ServerRegistryService {
	return &ServerRegistryService{msRepository: r}
}

func (s *ServerRegistryService) Create(
	ctx context.Context,
	loggedUser *auth.User,
	name valueobjects.MinecraftServerName,
) error {
	now := time.Now()
	server := entities.NewMinecraftServer(
		uuid.New(),
		loggedUser.ID,
		make([]uuid.UUID, 0),
		name,
		valueobjects.MinecraftServerStatusStopped,
		now,
		now,
	)
	return s.msRepository.Save(ctx, *server)
}

func (s *ServerRegistryService) GetAll(
	ctx context.Context,
	loggedUser *auth.User,
) ([]entities.MinecraftServer, error) {
	return s.msRepository.GetAllForUser(ctx, loggedUser.ID)
}

func (s *ServerRegistryService) Get(
	ctx context.Context,
	serverID uuid.UUID,
) (*entities.MinecraftServer, error) {
	return s.msRepository.GetById(ctx, serverID)
}
