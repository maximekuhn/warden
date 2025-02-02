package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/permissions"
)

type SqlitePermissionsBackend struct {
	db *sql.DB
}

func NewSqlitePermissionsBackend(db *sql.DB) *SqlitePermissionsBackend {
	return &SqlitePermissionsBackend{db: db}
}

func (s *SqlitePermissionsBackend) Save(ctx context.Context, uow transaction.UnitOfWork, user permissions.User) error {
	// TODO: save role for each server in user_role_server table
	// it should be done in a transaction
	suow := castUnitOfWorkOrPanic(uow)

	// save plan and user id in user_plan table
	query := `
    INSERT INTO user_policy_plan (user_id, user_plan) VALUES (?, ?)
    `
	_, err := suow.ExecContext(ctx, query, user.ID, planToDb(user.Plan))
	return err

}

func (s *SqlitePermissionsBackend) AddRole(
	ctx context.Context,
	uow transaction.UnitOfWork,
	userID uuid.UUID,
	serverID valueobjects.MinecraftServerID,
	role permissions.Role,
) error {
	suow := castUnitOfWorkOrPanic(uow)

	query := `
    INSERT INTO user_role_server (user_id, server_id, user_role)
    VALUES (?, ?, ?)
    `

	_, err := suow.ExecContext(ctx, query, userID, serverID.Value(), roleToDb(role))
	return err
}

func (s *SqlitePermissionsBackend) GetById(ctx context.Context, uow transaction.UnitOfWork, userID uuid.UUID) (*permissions.User, bool, error) {
	suow := castUnitOfWorkOrPanic(uow)

	query := `
    SELECT user_plan FROM user_policy_plan WHERE user_id = ?
    `
	var dbPlan int
	err := suow.QueryRowContext(ctx, query, userID).Scan(&dbPlan)
	if err != nil {
		return nil, false, err
	}

	plan, err := planFromDb(dbPlan)
	if err != nil {
		return nil, false, err
	}

	// note: we could fetch all data in the same query, but no need to optimize this yet
	queryRoles := `
    SELECT server_id, user_role
    FROM user_role_server
    `
	rows, err := suow.QueryContext(ctx, queryRoles)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	roles := make(map[uuid.UUID]permissions.Role)
	for rows.Next() {
		var serverID uuid.UUID
		var dbRole int

		if err := rows.Scan(&serverID, &dbRole); err != nil {
			return nil, false, err
		}

		role, err := dbRoleToRole(dbRole)
		if err != nil {
			return nil, false, err
		}

		roles[serverID] = role
	}

	user := permissions.NewUser(userID, plan, roles)
	return user, true, nil
}

func planToDb(plan permissions.Plan) int {
	switch plan {
	case permissions.PlanFree:
		return 1
	case permissions.PlanPro:
		return 2
	default:
		return -1
	}
}

func planFromDb(dbPlan int) (permissions.Plan, error) {
	switch dbPlan {
	case 1:
		return permissions.PlanFree, nil
	case 2:
		return permissions.PlanPro, nil
	default:
		return permissions.PlanFree, fmt.Errorf("corrupted plan: %d", dbPlan)
	}
}

func roleToDb(role permissions.Role) int {
	switch role {
	case permissions.RoleViewer:
		return 1
	case permissions.RoleAdmin:
		return 2
	default:
		return -1
	}
}

func dbRoleToRole(dbRole int) (permissions.Role, error) {
	switch dbRole {
	case 1:
		return permissions.RoleViewer, nil
	case 2:
		return permissions.RoleAdmin, nil
	default:
		return permissions.RoleViewer, fmt.Errorf("corrupted role: %d", dbRole)
	}
}
