package services

import (
	"context"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
)

type TeamServiceInterface interface {
	CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error)
	GetTeam(ctx context.Context, name string) (*models.Team, error)
}

type UserServiceInterface interface {
	SetUserIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
}

type PullRequestServiceInterface interface {
	CreatePullRequest(ctx context.Context, id, name, authorID string) (*models.PullRequest, error)
	MergePullRequest(ctx context.Context, id string) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (newReviewerID string, pr *models.PullRequest, err error)
	GetReviewerPRs(ctx context.Context, userID string) ([]*models.PullRequestShort, error)
}
