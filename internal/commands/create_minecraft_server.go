package commands

import (
	"context"
	"errors"

	"github.com/maximekuhn/warden/internal/transaction"
	"github.com/maximekuhn/warden/internal/valueobjects"
)

type CreateMinecraftServerCommand struct {
	Name valueobjects.MinecraftServerName
}

type CreateMinecraftServerCommandHandler struct {
	uowProvider transaction.UnitOfWorkProvider
}

func NewCreateMinecraftServerCommandHandler(uowProvider transaction.UnitOfWorkProvider) *CreateMinecraftServerCommandHandler {
	return &CreateMinecraftServerCommandHandler{uowProvider: uowProvider}
}

func (h *CreateMinecraftServerCommandHandler) Handle(
	ctx context.Context,
	cmd CreateMinecraftServerCommand,
) error {
	return errors.New("not yet implemented")

}
