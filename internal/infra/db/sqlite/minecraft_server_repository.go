package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type SqliteMinecraftServerRepository struct {
	db *sql.DB
}

func NewSqliteMinecraftServerRepository(db *sql.DB) *SqliteMinecraftServerRepository {
	return &SqliteMinecraftServerRepository{
		db: db,
	}
}
func (s *SqliteMinecraftServerRepository) Save(
	ctx context.Context,
	uow transaction.UnitOfWork,
	ms entities.MinecraftServer,
) error {
	suow := castUnitOfWorkOrPanic(uow)

	query := `
    INSERT INTO minecraft_server
    (id, owner_id, name, status, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?, ?)
    `

	_, err := suow.ExecContext(
		ctx,
		query,
		ms.ID.Value(),
		ms.OwnerID,
		ms.Name.Value(),
		ms.Status,
		ms.CreatedAt,
		ms.UpdatedAt,
	)
	return err
}

func msStatusToSqlite(s valueobjects.MinecraftServerStatus) int {
	switch s {
	case valueobjects.MinecraftServerStatusRunning:
		return 1
	case valueobjects.MinecraftServerStatusStopped:
		return 2
	default:
		panic("unreachable")
	}
}

func sqliteStatusToMsStatus(s int) (valueobjects.MinecraftServerStatus, error) {
	switch s {
	case 1:
		return valueobjects.MinecraftServerStatusRunning, nil
	case 2:
		return valueobjects.MinecraftServerStatusStopped, nil
	default:
		return "", fmt.Errorf("corrupted status (%d)", s)
	}
}
