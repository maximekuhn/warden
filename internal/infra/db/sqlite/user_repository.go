package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/warden/internal/domain/entities"
	"github.com/maximekuhn/warden/internal/domain/transaction"
	"github.com/maximekuhn/warden/internal/domain/valueobjects"
	"github.com/maximekuhn/warden/internal/permissions"
)

type SqliteUserRepository struct {
	db *sql.DB
}

func NewSqliteUserRepository(db *sql.DB) *SqliteUserRepository {
	return &SqliteUserRepository{
		db: db,
	}
}

func (s *SqliteUserRepository) GetAll(ctx context.Context, uow transaction.UnitOfWork, limit, offset uint) ([]entities.User, error) {
	suow := castUnitOfWorkOrPanic(uow)

	query := `
    SELECT auth.user_id, auth.email, auth.created_at, upa.user_plan
    FROM auth
    LEFT JOIN user_policy_plan upa ON auth.user_id = upa.user_id
    ORDER BY auth.email
    LIMIT ?
    OFFSET ?
    `

	rows, err := suow.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]entities.User, 0)
	for rows.Next() {
		var userId uuid.UUID
		var emailStr string
		var createdAt time.Time
		var userPlan int

		if err := rows.Scan(&userId, &emailStr, &createdAt, &userPlan); err != nil {
			return users, err
		}

		email, err := valueobjects.NewEmail(emailStr)
		if err != nil {
			return users, err
		}

		plan, err := planFromDb(userPlan)
		if err != nil {
			return users, err
		}

		users = append(
			users,
			*entities.NewUser(userId, email, plan, createdAt),
		)
	}
	return users, nil
}

func (s *SqliteUserRepository) Count(ctx context.Context, uow transaction.UnitOfWork) (uint, error) {
	suow := castUnitOfWorkOrPanic(uow)
	query := `SELECT COUNT(user_id) FROM auth`
	var count uint
	err := suow.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *SqliteUserRepository) UpdatePlan(
	ctx context.Context,
	uow transaction.UnitOfWork,
	userID uuid.UUID,
	newPlan permissions.Plan,
) error {
	suow := castUnitOfWorkOrPanic(uow)
	query := `UPDATE user_policy_plan SET user_plan = ? WHERE user_id = ?`
	res, err := suow.ExecContext(ctx, query, planToDb(newPlan), userID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("expected to affect 1 row, affect %d", rowsAffected)
	}
	return nil
}
