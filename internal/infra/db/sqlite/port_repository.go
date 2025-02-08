package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
)

type SqlitePortRepository struct {
	db *sql.DB
}

func NewSqlitePortRepository(db *sql.DB) *SqlitePortRepository {
	return &SqlitePortRepository{
		db: db,
	}
}

func (s *SqlitePortRepository) Save(
	ctx context.Context,
	uow transaction.UnitOfWork,
	port uint16,
	serverID valueobjects.MinecraftServerID,
) error {
	suow := castUnitOfWorkOrPanic(uow)

	query := `
    INSERT INTO minecraft_server_port (server_id, port) VALUES (?, ?)
    `
	_, err := suow.ExecContext(ctx, query, serverID.Value(), port)
	return err
}

func (s *SqlitePortRepository) GetByServerID(
	ctx context.Context,
	uow transaction.UnitOfWork,
	serverID valueobjects.MinecraftServerID,
) (uint16, bool, error) {
	suow := castUnitOfWorkOrPanic(uow)

	query := `
    SELECT port FROM minecraft_server_port WHERE server_id = ?
    `
	var port uint16
	err := suow.QueryRowContext(ctx, query, serverID.Value()).Scan(&port)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return port, true, nil
}
