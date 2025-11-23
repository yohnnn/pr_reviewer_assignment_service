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

func TestPullRequestService_Simple(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	userRepo := inmemory.NewTestUserRepo()
	userRepo.Users["u1"] = &models.User{ID: "u1", TeamName: "backend", IsActive: true}
	userRepo.Users["u2"] = &models.User{ID: "u2", TeamName: "backend", IsActive: true}
	userRepo.Users["u3"] = &models.User{ID: "u3", TeamName: "backend", IsActive: true}
	userRepo.Users["u4"] = &models.User{ID: "u4", TeamName: "backend", IsActive: true}

	prRepo := inmemory.NewTestPRRepo()

	service := NewPullRequestService(prRepo, userRepo, logger)
	ctx := context.Background()

	t.Run("Create PR", func(t *testing.T) {
		pr, err := service.CreatePullRequest(ctx, "pr-1", "Fix bug", "u1")

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if pr.Status != "OPEN" {
			t.Error("PR should be OPEN")
		}
		if len(pr.Reviewers) != 2 {
			t.Errorf("Expected 2 reviewers, got %d", len(pr.Reviewers))
		}
		for _, r := range pr.Reviewers {
			if r == "u1" {
				t.Error("Author cannot be a reviewer")
			}
		}
	})

	t.Run("Reassign Reviewer", func(t *testing.T) {
		pr, _ := prRepo.GetByID(ctx, "pr-1")
		oldRev := pr.Reviewers[0]

		prUpdated, newRev, err := service.ReassignReviewer(ctx, "pr-1", oldRev)

		if err != nil {
			t.Fatalf("Reassign failed: %v", err)
		}
		if newRev == oldRev {
			t.Error("Reviewer ID did not change")
		}

		found := false
		for _, r := range prUpdated.Reviewers {
			if r == newRev {
				found = true
			}
		}
		if !found {
			t.Error("New reviewer not found in list")
		}
	})

	t.Run("Merge PR", func(t *testing.T) {
		pr, err := service.MergePullRequest(ctx, "pr-1")

		if err != nil {
			t.Fatalf("Merge failed: %v", err)
		}
		if pr.Status != "MERGED" {
			t.Error("Status should be MERGED")
		}

		pr2, err2 := service.MergePullRequest(ctx, "pr-1")
		if err2 != nil {
			t.Error("Idempotent merge failed")
		}
		if pr2.Status != "MERGED" {
			t.Error("Status should stay MERGED")
		}
	})

	t.Run("Reassign on MERGED PR", func(t *testing.T) {
		pr, _ := prRepo.GetByID(ctx, "pr-1")
		reviewer := pr.Reviewers[0]

		_, _, err := service.ReassignReviewer(ctx, "pr-1", reviewer)

		if !errors.Is(err, models.ErrPRMerged) {
			t.Errorf("Expected ErrPRMerged, got %v", err)
		}
	})
}
