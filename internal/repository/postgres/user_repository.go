package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUser(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, name, team_name, is_active FROM users WHERE id = $1`
	var u models.User
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err) // Можно возвращать models.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) SetUserIsActive(ctx context.Context, id string, isActive bool) error {
	query := `UPDATE users SET is_active = $1 WHERE id = $2`
	res, err := r.db.Exec(ctx, query, isActive, id)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *UserRepository) GetActiveUsersByTeam(ctx context.Context, teamName string) ([]models.User, error) {
	query := `
        SELECT id, name, team_name, is_active 
        FROM users 
        WHERE team_name = $1 AND is_active = true
    `
	rows, err := r.db.Query(ctx, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to query active users for team %s: %w", teamName, err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.User, error) {
		var u models.User
		err := row.Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive)
		return u, err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to collect active users: %w", err)
	}

	return users, nil
}

// func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
// 	query := `
//         INSERT INTO users (id, name, team_name, is_active)
//         VALUES ($1, $2, $3, $4)
//         ON CONFLICT (id) DO UPDATE
//         SET name = EXCLUDED.name,
//             team_name = EXCLUDED.team_name,
//             is_active = EXCLUDED.is_active
//     	`
// 	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.TeamName, user.IsActive)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
