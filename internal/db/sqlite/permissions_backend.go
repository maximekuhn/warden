package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/permissions"
	"github.com/maximekuhn/warden/internal/transaction"
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
func (s *SqlitePermissionsBackend) GetById(ctx context.Context, uow transaction.UnitOfWork, userID uuid.UUID) (*permissions.User, bool, error) {
	// TODO: retrieve role from user_role_server table
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

	user := permissions.NewUser(userID, plan, make(map[uuid.UUID]permissions.Role))
	return user, false, nil
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
