package services

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/repository/inmemory"
)

func TestUserService_Simple(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := inmemory.NewTestUserRepo()

	repo.Users["u1"] = &models.User{ID: "u1", Name: "Vasya", IsActive: true}

	service := NewUserService(repo, logger)
	ctx := context.Background()

	t.Run("Deactivate existing user", func(t *testing.T) {

		user, err := service.SetUserIsActive(ctx, "u1", false)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if user.IsActive != false {
			t.Errorf("Expected IsActive to be false, got true")
		}

		savedUser, _ := repo.GetUser(ctx, "u1")
		if savedUser.IsActive != false {
			t.Errorf("User in DB should be inactive")
		}
	})

	t.Run("Update non-existing user", func(t *testing.T) {
		_, err := service.SetUserIsActive(ctx, "ghost_user", true)

		if err == nil {
			t.Error("Expected error for non-existing user, got nil")
		}
		if !errors.Is(err, models.ErrNotFound) {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	t.Run("Mass deactivate", func(t *testing.T) {

		repo.Users["u2"] = &models.User{ID: "u2", IsActive: true}
		repo.Users["u3"] = &models.User{ID: "u3", IsActive: true}

		err := service.DeactivateUsers(ctx, []string{"u2", "u3"})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		u2, _ := repo.GetUser(ctx, "u2")
		u3, _ := repo.GetUser(ctx, "u3")

		if u2.IsActive {
			t.Error("u2 must be inactive")
		}
		if u3.IsActive {
			t.Error("u3 must be inactive")
		}
	})
}
