package inmemory

import (
	"context"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
)

type TestStatsRepo struct {
	MockData []*models.ReviewerStat
	MockErr  error
}

func NewTestStatsRepo() *TestStatsRepo {
	return &TestStatsRepo{
		MockData: make([]*models.ReviewerStat, 0),
	}
}

func (r *TestStatsRepo) GetTopReviewers(ctx context.Context) ([]*models.ReviewerStat, error) {
	if r.MockErr != nil {
		return nil, r.MockErr
	}
	if r.MockData == nil {
		return []*models.ReviewerStat{}, nil
	}
	return r.MockData, nil
}
