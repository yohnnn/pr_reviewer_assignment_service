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
	if isActive {
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

	user, err := s.repo.GetUser(ctx, userID)

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			s.log.WarnContext(ctx, "user not found", "user_id", userID)
			return nil, err
		} else {
			s.log.ErrorContext(ctx, "failed to get user", "err", err)
			return nil, err
		}
	}

	if err := s.repo.DeactivateUsers(ctx, []string{userID}); err != nil {
		s.log.ErrorContext(ctx, "failed to deactivate user", "err", err)
		return nil, err
	}

	user.IsActive = false

	return user, nil
}

func (s *UserService) DeactivateUsers(ctx context.Context, userIDs []string) error {
	s.log.InfoContext(ctx, "deactivating users", "user_ids", userIDs)
	if err := s.repo.DeactivateUsers(ctx, userIDs); err != nil {
		s.log.ErrorContext(ctx, "failed to deactivate users", "err", err)
		return err
	}
	return nil
}
