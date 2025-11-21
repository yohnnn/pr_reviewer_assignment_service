package repository

import (
	"context"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
)

type TeamRepositoryInterface interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeam(ctx context.Context, name string) (*models.Team, error)
}

type UserRepositoryInterface interface {
	GetUser(ctx context.Context, id string) (*models.User, error)
	SetUserIsActive(ctx context.Context, id string, isActive bool) (*models.User, error)
	GetActiveUsersByTeam(ctx context.Context, teamName string) ([]models.User, error)
}

type PullRequestRepositoryInterface interface {
	CreatePR(ctx context.Context, pr *models.PullRequest) error
	MergePR(ctx context.Context, id string) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, id, oldUserID, newUserID string) error
	GetByReviewerID(ctx context.Context, userID string) ([]*models.PullRequestShort, error)
	GetByID(ctx context.Context, id string) (*models.PullRequest, error)
}
