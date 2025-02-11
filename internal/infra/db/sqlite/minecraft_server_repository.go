package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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
		msStatusToSqlite(ms.Status),
		ms.CreatedAt,
		ms.UpdatedAt,
	)
	if err != nil {
	}
	return err
}

func (s *SqliteMinecraftServerRepository) GetAllForUser(
	ctx context.Context,
	uow transaction.UnitOfWork,
	userID uuid.UUID,
) ([]entities.MinecraftServer, error) {
	// TODO: handle when user has different role than owner
	suow := castUnitOfWorkOrPanic(uow)

	query := `
    SELECT id, owner_id, name, status, created_at, updated_at
    FROM minecraft_server
    WHERE owner_id = ?
    `

	rows, err := suow.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	servers := make([]entities.MinecraftServer, 0)
	for rows.Next() {
		server, err := convertMinecraftServerRows(rows)
		if err != nil {
			return servers, nil
		}
		servers = append(servers, *server)
	}

	return servers, nil
}

func (s *SqliteMinecraftServerRepository) Update(
	ctx context.Context,
	uow transaction.UnitOfWork,
	old, new entities.MinecraftServer,
) error {
	suow := castUnitOfWorkOrPanic(uow)

	if old.ID != new.ID {
		return errors.New("old and new have different IDs.. wtf!!")
	}

	// TODO: check other fields (Updated at !!), not only status
	// TODO: remove this check, it's not the db scope
	if old.Status == new.Status {
		return errors.New("nothing to update (only checked status field)")
	}
	query := `
    UPDATE minecraft_server SET status = ?, updated_at = ? WHERE id = ?
    `
	res, err := suow.ExecContext(
		ctx,
		query,
		msStatusToSqlite(new.Status),
		new.UpdatedAt,
		new.ID.Value(),
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return errors.New("sqlite: update did not update any row")
	}
	return nil
}

func (s *SqliteMinecraftServerRepository) GetByID(
	ctx context.Context,
	uow transaction.UnitOfWork,
	serverID valueobjects.MinecraftServerID,
) (*entities.MinecraftServer, bool, error) {
	// TODO: get members
	suow := castUnitOfWorkOrPanic(uow)

	query := `
    SELECT id, owner_id, name, status, created_at, updated_at
    FROM minecraft_server
    WHERE id = ?
    `
	row := suow.QueryRowContext(ctx, query, serverID.Value())
	srv, err := convertMinecraftServerRow(row)
	return srv, true, err
}

func convertMinecraftServerRows(rows *sql.Rows) (*entities.MinecraftServer, error) {
	var id uuid.UUID
	var ownerID uuid.UUID
	var name string
	var status int
	var createdAt time.Time
	var updatedAt time.Time

	if err := rows.Scan(&id, &ownerID, &name, &status, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	serverID, err := valueobjects.NewMinecraftServerID(id)
	if err != nil {
		return nil, err
	}

	serverName, err := valueobjects.NewMinecraftServerName(name)
	if err != nil {
		return nil, err
	}

	serverStatus, err := sqliteStatusToMsStatus(status)
	if err != nil {
		return nil, err
	}
	return entities.NewMinecraftServer(
		serverID,
		ownerID,
		make([]uuid.UUID, 0),
		serverName,
		serverStatus,
		createdAt,
		updatedAt,
	), nil
}

func convertMinecraftServerRow(row *sql.Row) (*entities.MinecraftServer, error) {
	var id uuid.UUID
	var ownerID uuid.UUID
	var name string
	var status int
	var createdAt time.Time
	var updatedAt time.Time

	if err := row.Scan(&id, &ownerID, &name, &status, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	serverID, err := valueobjects.NewMinecraftServerID(id)
	if err != nil {
		return nil, err
	}

	serverName, err := valueobjects.NewMinecraftServerName(name)
	if err != nil {
		return nil, err
	}

	serverStatus, err := sqliteStatusToMsStatus(status)
	if err != nil {
		return nil, err
	}
	return entities.NewMinecraftServer(
		serverID,
		ownerID,
		make([]uuid.UUID, 0),
		serverName,
		serverStatus,
		createdAt,
		updatedAt,
	), nil
}

func msStatusToSqlite(s valueobjects.MinecraftServerStatus) int {
	switch s {
	case valueobjects.MinecraftServerStatusRunning:
		return 1
	case valueobjects.MinecraftServerStatusStopped:
		return 2
	case valueobjects.MinecraftServerStatusStarting:
		return 3
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
	case 3:
		return valueobjects.MinecraftServerStatusStarting, nil
	default:
		return "", fmt.Errorf("corrupted status (%d)", s)
	}
}
