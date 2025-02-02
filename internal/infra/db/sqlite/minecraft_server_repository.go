package sqlite

import (
	"context"
	"database/sql"

	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/transaction"
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
		ms.CreatedAt,
		ms.UpdatedAt,
	)
	return err
}
