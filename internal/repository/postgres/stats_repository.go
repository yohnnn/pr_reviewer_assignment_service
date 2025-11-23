package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
)

type StatsRepository struct {
	db *pgxpool.Pool
}

func NewStatsRepository(db *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetTopReviewers(ctx context.Context) ([]*models.ReviewerStat, error) {

	query := `
		SELECT u.id, u.name, COUNT(prr.pr_id) as review_count
		FROM users u
		LEFT JOIN pr_reviewers prr ON u.id = prr.reviewer_id
		GROUP BY u.id, u.name
		ORDER BY review_count DESC
		LIMIT 10
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query stats: %w", err)
	}
	defer rows.Close()

	stats, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.ReviewerStat, error) {
		var s models.ReviewerStat
		err := row.Scan(&s.UserID, &s.Username, &s.ReviewCount)
		return &s, err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to collect stats: %w", err)
	}

	return stats, nil
}
