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

func TestTeamService_Simple(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := inmemory.NewTestTeamRepo()
	service := NewTeamService(repo, logger)
	ctx := context.Background()

	t.Run("Create new team", func(t *testing.T) {
		newTeam := &models.Team{
			Name: "backend",
			Members: []models.TeamMember{
				{UserID: "u1", UserName: "Alice", IsActive: true},
			},
		}

		created, err := service.CreateTeam(ctx, newTeam)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if created.Name != "backend" {
			t.Errorf("Expected team name 'backend', got %s", created.Name)
		}

		exist, _ := repo.GetTeam(ctx, "backend")
		if exist == nil {
			t.Error("Team was not saved to DB")
		}
	})

	t.Run("Create duplicate team", func(t *testing.T) {

		dupTeam := &models.Team{Name: "backend"}

		_, err := service.CreateTeam(ctx, dupTeam)

		if err == nil {
			t.Error("Expected error for duplicate team, got nil")
		}
		if !errors.Is(err, models.ErrAlreadyExists) {
			t.Errorf("Expected ErrAlreadyExists, got %v", err)
		}
	})

	t.Run("Get missing team", func(t *testing.T) {
		_, err := service.GetTeam(ctx, "frontend")

		if err == nil {
			t.Error("Expected error for missing team, got nil")
		}
		if !errors.Is(err, models.ErrNotFound) {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})
}
