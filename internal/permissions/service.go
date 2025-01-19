package permissions

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/auth"
	"github.com/maximekuhn/warden/internal/transaction"
)

type PermissionsService struct {
	backend Backend
}

func NewPermissionsService(b Backend) *PermissionsService {
	return &PermissionsService{
		backend: b,
	}
}

func (s *PermissionsService) Create(
	ctx context.Context,
	uow transaction.UnitOfWork,
	userID uuid.UUID,
	plan Plan,
) error {
	user := NewUser(userID, plan, make(map[uuid.UUID]Role))
	return s.backend.Save(ctx, uow, *user)
}

// HasMinecraftServerPermission returns if the logged user has permission
// to perform the provided action in the specific Minecraft server.
//
// If a non-nil error is returned, the result should be discarded.
func (s *PermissionsService) HasMinecraftServerPermission(
	ctx context.Context,
	uow transaction.UnitOfWork,
	loggedUser *auth.User,
	serverID uuid.UUID,
	action Action,
) (bool, error) {
	user, found, err := s.backend.GetById(ctx, uow, loggedUser.ID)
	if err != nil {
		return false, err
	}
	if !found {
		return false, ErrUserNotFound
	}

	role, found := user.Roles[serverID]
	if !found {
		return false, ErrUserNotInServer
	}

	actions, found := roleToActions[role]
	if !found {
		return false, fmt.Errorf("no actions found for role '%s'", role)
	}

	for _, act := range actions {
		if action == act {
			// user is allowed to perform action in this server
			return true, nil
		}
	}
	return false, nil
}

// HasPolicy returns if the logged user has the correct plan for the provided policy.
// If a non-nil error is returned, the result should be discarded.
func (s *PermissionsService) HasPolicy(
	ctx context.Context,
	uow transaction.UnitOfWork,
	loggedUser *auth.User,
	policy Policy,
) (bool, error) {
	user, found, err := s.backend.GetById(ctx, uow, loggedUser.ID)
	if err != nil {
		return false, err
	}
	if !found {
		return false, ErrUserNotFound
	}

	policies, found := planToPolicies[user.Plan]
	if !found {
		return false, fmt.Errorf("no policies found for plan '%s'", user.Plan)
	}

	for _, p := range policies {
		if p == policy {
			// user has plan for this policy
			return true, nil
		}
	}
	return false, nil
}
