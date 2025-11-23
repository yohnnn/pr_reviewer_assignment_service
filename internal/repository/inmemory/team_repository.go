package inmemory

import (
	"context"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
)

type TestTeamRepo struct {
	Teams map[string]*models.Team
}

func NewTestTeamRepo() *TestTeamRepo {
	return &TestTeamRepo{
		Teams: make(map[string]*models.Team),
	}
}

func (s *TestTeamRepo) CreateTeam(ctx context.Context, team *models.Team) error {
	if _, exists := s.Teams[team.Name]; exists {
		return models.ErrAlreadyExists
	}

	s.Teams[team.Name] = team
	return nil
}

func (s *TestTeamRepo) GetTeam(ctx context.Context, name string) (*models.Team, error) {
	t, ok := s.Teams[name]
	if !ok {
		return nil, models.ErrNotFound
	}
	return t, nil
}
