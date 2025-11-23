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

func TestStatsService_Simple(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := inmemory.NewTestStatsRepo()
	ctx := context.Background()

	t.Run("Success: Returns top reviewers", func(t *testing.T) {
		repo.MockData = []*models.ReviewerStat{
			{UserID: "u1", Username: "Alice", ReviewCount: 10},
			{UserID: "u2", Username: "Bob", ReviewCount: 5},
		}

		service := NewStatsService(repo, logger)

		got, err := service.GetTopReviewers(ctx)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(got) != 2 {
			t.Errorf("Expected 2 items, got %d", len(got))
		}
		if got[0].Username != "Alice" {
			t.Errorf("Wrong data in first item")
		}
	})

	t.Run("Fail: Database error", func(t *testing.T) {
		repo.MockData = nil
		repo.MockErr = errors.New("db connection failed")

		service := NewStatsService(repo, logger)

		got, err := service.GetTopReviewers(ctx)

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if got != nil {
			t.Error("Expected nil result on error")
		}
	})
}
