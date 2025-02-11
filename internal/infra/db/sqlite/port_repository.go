package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/maximekuhn/warden/internal/domain/services"
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
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return services.ErrNoPortAvailable
		}
	}
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

func (s *SqlitePortRepository) GetAll(ctx context.Context, uow transaction.UnitOfWork) ([]uint16, error) {
	suow := castUnitOfWorkOrPanic(uow)

	query := `
    SELECT port FROM minecraft_server_port
    `

	rows, err := suow.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ports := make([]uint16, 0)
	for rows.Next() {
		var port uint16
		if err := rows.Scan(&port); err != nil {
			return ports, err
		}
		ports = append(ports, port)
	}

	return ports, nil
}
