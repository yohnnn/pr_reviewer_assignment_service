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

func (r *UserRepository) DeactivateUsers(ctx context.Context, userIDs []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	queryDeactivate := `UPDATE users SET is_active = false WHERE id = ANY($1)`
	_, err = tx.Exec(ctx, queryDeactivate, userIDs)
	if err != nil {
		return fmt.Errorf("failed to deactivate users: %w", err)
	}

	query := `
		SELECT prr.pr_id, prr.reviewer_id, pr.author_id, u.team_name
		FROM pr_reviewers prr
		JOIN pull_requests pr ON prr.pr_id = pr.id
		JOIN users u ON pr.author_id = u.id
		WHERE prr.reviewer_id = ANY($1) 
		  AND pr.status = 'OPEN'
	`

	rows, err := tx.Query(ctx, query, userIDs)
	if err != nil {
		return fmt.Errorf("failed to select reviews to update: %w", err)
	}

	reviews, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.ReviewToUpdate, error) {
		var item models.ReviewToUpdate
		err := row.Scan(&item.PRID, &item.OldReviewerID, &item.AuthorID, &item.TeamName)
		return item, err
	})
	if err != nil {
		return fmt.Errorf("failed to collect reviews to update: %w", err)
	}

	queryDeleteReviewer := "DELETE FROM pr_reviewers WHERE pr_id=$1 AND reviewer_id=$2"

	queryUpdateReviewer := "UPDATE pr_reviewers SET reviewer_id=$1 WHERE pr_id=$2 AND reviewer_id=$3"

	for _, rev := range reviews {
		var newReviewerID string

		queryCandidate := `
			SELECT id FROM users
			WHERE team_name = $1
			  AND is_active = true
			  AND id != $2
			  AND id NOT IN (
			      SELECT reviewer_id FROM pr_reviewers WHERE pr_id = $3
			  )
			LIMIT 1
		`

		err := tx.QueryRow(ctx, queryCandidate, rev.TeamName, rev.AuthorID, rev.PRID).Scan(&newReviewerID)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				if _, err := tx.Exec(ctx, queryDeleteReviewer, rev.PRID, rev.OldReviewerID); err != nil {
					return fmt.Errorf("failed to delete reviewer: %w", err)
				}
			} else {
				return fmt.Errorf("failed to find candidate: %w", err)
			}
		} else {
			if _, err := tx.Exec(ctx, queryUpdateReviewer, newReviewerID, rev.PRID, rev.OldReviewerID); err != nil {
				return fmt.Errorf("failed to update reviewer: %w", err)
			}
		}
	}

	return tx.Commit(ctx)
}
