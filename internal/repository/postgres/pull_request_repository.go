package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
)

type PullRequestRepository struct {
	db *pgxpool.Pool
}

func NewPullRequestRepository(db *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

// CreatePR соответствует интерфейсу
func (r *PullRequestRepository) CreatePR(ctx context.Context, pr *models.PullRequest) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	queryPR := `
        INSERT INTO pull_requests (id, name, author_id, status, created_at, merged_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err = tx.Exec(ctx, queryPR, pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt)
	if err != nil {
		if IsUnique(err) {
			return fmt.Errorf("pr exists: %w", models.ErrAlreadyExists)
		}
		return fmt.Errorf("failed to insert pr: %w", err)
	}

	if len(pr.Reviewers) > 0 {
		queryReviewers := `
            INSERT INTO pr_reviewers (pr_id, reviewer_id)
            VALUES ($1, $2)
        `
		for _, reviewerID := range pr.Reviewers {
			if _, err := tx.Exec(ctx, queryReviewers, pr.ID, reviewerID); err != nil {
				return fmt.Errorf("failed to insert reviewer %s: %w", reviewerID, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}

// MergePR соответствует интерфейсу
func (r *PullRequestRepository) MergePR(ctx context.Context, id string) error {
	query := `
        UPDATE pull_requests 
        SET status = 'MERGED', merged_at = $2
        WHERE id = $1
    `
	res, err := r.db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to merge pr: %w", err)
	}
	if res.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (r *PullRequestRepository) ReassignReviewer(ctx context.Context, prID, oldUserID, newUserID string) error {
	query := `
        UPDATE pr_reviewers
        SET reviewer_id = $3
        WHERE pr_id = $1 AND reviewer_id = $2
    `
	res, err := r.db.Exec(ctx, query, prID, oldUserID, newUserID)
	if err != nil {
		if IsUnique(err) {
			return fmt.Errorf("reviewer already assigned: %w", models.ErrAlreadyExists)
		}
		return fmt.Errorf("failed to reassign: %w", err)
	}
	if res.RowsAffected() == 0 {
		return models.ErrNotAssigned
	}
	return nil
}

// GetByReviewerID возвращает []*models.PullRequestShort
func (r *PullRequestRepository) GetByReviewerID(ctx context.Context, userID string) ([]*models.PullRequestShort, error) {
	query := `
        SELECT pr.id, pr.name, pr.author_id, pr.status
        FROM pull_requests pr
        INNER JOIN pr_reviewers prr ON pr.id = prr.pr_id
        WHERE prr.reviewer_id = $1
        ORDER BY pr.created_at DESC
    `
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query prs: %w", err)
	}
	defer rows.Close()

	// Используем CollectRows с возвратом указателя
	result, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.PullRequestShort, error) {
		var pr models.PullRequestShort
		err := row.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status)
		return &pr, err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to collect prs: %w", err)
	}

	return result, nil
}

// GetByID нужен для сервиса (хотя в интерфейсе выше его нет, но мы договорились его добавить)
func (r *PullRequestRepository) GetByID(ctx context.Context, id string) (*models.PullRequest, error) {
	query := `
        SELECT id, name, author_id, status, created_at, merged_at 
        FROM pull_requests 
        WHERE id = $1
    `
	var pr models.PullRequest
	err := r.db.QueryRow(ctx, query, id).Scan(
		&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get pr: %w", err)
	}

	queryRev := `SELECT reviewer_id FROM pr_reviewers WHERE pr_id = $1`
	rows, err := r.db.Query(ctx, queryRev, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewers: %w", err)
	}
	defer rows.Close()

	reviewers, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("failed to collect reviewers: %w", err)
	}
	pr.Reviewers = reviewers

	return &pr, nil
}
