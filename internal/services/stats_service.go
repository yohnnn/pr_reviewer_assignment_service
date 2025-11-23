package services

import (
	"context"
	"log/slog"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/repository"
)

type StatsService struct {
	repo repository.StatsRepositoryInterface
	log  *slog.Logger
}

func NewStatsService(repo repository.StatsRepositoryInterface, log *slog.Logger) *StatsService {
	return &StatsService{repo: repo, log: log}
}

func (s *StatsService) GetTopReviewers(ctx context.Context) ([]*models.ReviewerStat, error) {
	stats, err := s.repo.GetTopReviewers(ctx)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to get reviewer stats", "error", err)
		return nil, err
	}
	return stats, err
}
