package inmemory

import (
	"context"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
)

type TestUserRepo struct {
	Users map[string]*models.User
}

func NewTestUserRepo() *TestUserRepo {
	return &TestUserRepo{
		Users: make(map[string]*models.User),
	}
}

func (s *TestUserRepo) GetUser(ctx context.Context, id string) (*models.User, error) {
	u, ok := s.Users[id]
	if !ok {
		return nil, models.ErrNotFound
	}

	return u, nil
}

func (s *TestUserRepo) SetUserIsActive(ctx context.Context, id string, isActive bool) (*models.User, error) {
	u, ok := s.Users[id]
	if !ok {
		return nil, models.ErrNotFound
	}
	u.IsActive = isActive
	return u, nil
}

func (s *TestUserRepo) DeactivateUsers(ctx context.Context, userIDs []string) error {
	for _, id := range userIDs {
		if u, ok := s.Users[id]; ok {
			u.IsActive = false
		}
	}
	return nil
}

func (r *TestUserRepo) GetActiveUsersByTeam(ctx context.Context, teamName string) ([]models.User, error) {
	var result []models.User
	for _, u := range r.Users {
		if u.TeamName == teamName && u.IsActive {
			result = append(result, *u)
		}
	}
	return result, nil
}
