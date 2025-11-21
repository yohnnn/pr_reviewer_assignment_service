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
	if err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", models.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) SetUserIsActive(ctx context.Context, id string, isActive bool) (*models.User, error) {
	query := `
        UPDATE users
        SET is_active = $2
        WHERE id = $1
        RETURNING id, name, team_name, is_active
    `
	var u models.User
	if err := r.db.QueryRow(ctx, query, id, isActive).Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &u, nil
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
