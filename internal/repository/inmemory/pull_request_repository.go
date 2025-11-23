package inmemory

import (
	"context"
	"time"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
)

type TestPRRepo struct {
	Prs map[string]*models.PullRequest
}

func NewTestPRRepo() *TestPRRepo {
	return &TestPRRepo{
		Prs: make(map[string]*models.PullRequest),
	}
}

func (r *TestPRRepo) CreatePR(ctx context.Context, pr *models.PullRequest) error {
	if _, ok := r.Prs[pr.ID]; ok {
		return models.ErrAlreadyExists
	}
	r.Prs[pr.ID] = pr
	return nil
}

func (r *TestPRRepo) MergePR(ctx context.Context, id string) (*models.PullRequest, error) {
	pr, ok := r.Prs[id]
	if !ok {
		return nil, models.ErrNotFound
	}
	pr.Status = "MERGED"
	now := time.Now()
	pr.MergedAt = &now
	return pr, nil
}

func (r *TestPRRepo) GetByID(ctx context.Context, id string) (*models.PullRequest, error) {
	pr, ok := r.Prs[id]
	if !ok {
		return nil, models.ErrNotFound
	}
	return pr, nil
}

func (r *TestPRRepo) ReassignReviewer(ctx context.Context, id, oldID, newID string) error {
	pr, ok := r.Prs[id]
	if !ok {
		return nil
	}
	for i, rev := range pr.Reviewers {
		if rev == oldID {
			pr.Reviewers[i] = newID
			return nil
		}
	}
	return nil
}

func (r *TestPRRepo) GetByReviewerID(_ context.Context, _ string) ([]*models.PullRequestShort, error) {
	return nil, nil
}
