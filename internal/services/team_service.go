package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/repository"
)

type TeamService struct {
	repo repository.TeamRepositoryInterface
	log  *slog.Logger
}

func NewTeamService(repo repository.TeamRepositoryInterface, log *slog.Logger) *TeamService {
	return &TeamService{
		repo: repo,
		log:  log,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	s.log.InfoContext(ctx, "creating team", "team_name", team.Name)

	if err := s.repo.CreateTeam(ctx, team); err != nil {
		if errors.Is(err, models.ErrAlreadyExists) {
			s.log.WarnContext(ctx, "team already exist", "team_name", team.Name)
		} else {
			s.log.ErrorContext(ctx, "failed to create team", "err", err)
		}
		return nil, err
	}
	return team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*models.Team, error) {

	team, err := s.repo.GetTeam(ctx, name)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			s.log.WarnContext(ctx, "team not found", "team", name)
		} else {
			s.log.ErrorContext(ctx, "failed to get team", "err", err)
		}
		return nil, err
	}

	return team, nil
}
