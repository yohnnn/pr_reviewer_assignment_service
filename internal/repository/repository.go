package repository

import (
	"context"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team, members []models.User) error
	GetTeam(ctx context.Context, name string) (*models.Team, error)
}

type UserRepository interface {
	// CreateUser(ctx context.Context, user *models.User) error
	SetUserIsActive(ctx context.Context, id string, isActive bool) error
	GetActiveUsersByTeam(ctx context.Context, teamName string) ([]models.User, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr *models.PullRequest) error
	Merge(ctx context.Context, id string) error
	ReassignReviewer(ctx context.Context, id, oldUserID, newUserID string) error
	GetByReviewerID(ctx context.Context, userID string) ([]*models.PullRequest, error)
}
