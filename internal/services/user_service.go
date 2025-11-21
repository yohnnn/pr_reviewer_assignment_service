package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/repository"
)

type UserService struct {
	repo repository.UserRepositoryInterface
	log  *slog.Logger
}

func NewUserService(repo repository.UserRepositoryInterface, log *slog.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
	}
}

func (s *UserService) SetUserIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	s.log.InfoContext(ctx, "setting user active status", "user_id", userID, "is_active", isActive)

	user, err := s.repo.SetUserIsActive(ctx, userID, isActive)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			s.log.WarnContext(ctx, "user not found for update", "user_id", userID)
		} else {
			s.log.ErrorContext(ctx, "failed to update user", "err", err)
		}
		return nil, err
	}

	return user, nil
}
