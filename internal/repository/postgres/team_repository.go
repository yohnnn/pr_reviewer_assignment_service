package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
)

type TeamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team *models.Team) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	queryTeam := `INSERT INTO teams (name) VALUES ($1)`
	if _, err = tx.Exec(ctx, queryTeam, team.Name); err != nil {
		if IsUnique(err) {
			return fmt.Errorf("team exists: %w", models.ErrAlreadyExists)
		}
		return fmt.Errorf("failed to insert team %s: %w", team.Name, err)
	}

	queryUser := `
        INSERT INTO users (id, name, team_name, is_active) 
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE 
        SET name = EXCLUDED.name, 
            team_name = EXCLUDED.team_name, 
            is_active = EXCLUDED.is_active
    `

	for _, member := range team.Members {
		if _, err = tx.Exec(ctx, queryUser, member.UserID, member.UserName, team.Name, member.IsActive); err != nil {
			return fmt.Errorf("failed to upsert member %s: %w", member.UserID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *TeamRepository) GetTeam(ctx context.Context, name string) (*models.Team, error) {
	var found bool
	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1)`
	if err := r.db.QueryRow(ctx, query, name).Scan(&found); err != nil {
		return nil, fmt.Errorf("failed to check team existence: %w", err)
	}

	if !found {
		return nil, models.ErrNotFound
	}

	query = `
        SELECT id, name, is_active 
        FROM users 
        WHERE team_name = $1
    `
	rows, err := r.db.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to query members for team %s: %w", name, err)
	}
	defer rows.Close()

	members, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.TeamMember, error) {
		var u models.TeamMember
		err := row.Scan(&u.UserID, &u.UserName, &u.IsActive)
		return u, err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to collect members: %w", err)
	}

	return &models.Team{
		Name:    name,
		Members: members,
	}, nil
}
