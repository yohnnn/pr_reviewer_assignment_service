package services

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"slices"
	"time"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/repository"
)

type PullRequestService struct {
	prRepo   repository.PullRequestRepositoryInterface
	userRepo repository.UserRepositoryInterface
	log      *slog.Logger
}

func NewPullRequestService(
	prRepo repository.PullRequestRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	log *slog.Logger,
) *PullRequestService {
	return &PullRequestService{
		prRepo:   prRepo,
		userRepo: userRepo,
		log:      log,
	}
}

func (s *PullRequestService) CreatePullRequest(ctx context.Context, id, name, authorID string) (*models.PullRequest, error) {
	s.log.InfoContext(ctx, "creating PR", "PR_id", id)
	user, err := s.userRepo.GetUser(ctx, authorID)

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			s.log.WarnContext(ctx, "user not found", "user_id", authorID)
		} else {
			s.log.ErrorContext(ctx, "failed to get user", "err", err)
		}
		return nil, err
	}

	activeUsers, err := s.userRepo.GetActiveUsersByTeam(ctx, user.TeamName)

	if err != nil {
		s.log.ErrorContext(ctx, "failed to get active users", "err", err)
		return nil, err
	}

	candidates := make([]string, 0, len(activeUsers))
	for _, user := range activeUsers {
		if user.ID != authorID {
			candidates = append(candidates, user.ID)
		}
	}

	var reviewers []string
	if len(candidates) <= 2 {
		reviewers = candidates
	} else {
		idx1 := rand.Intn(len(candidates))
		idx2 := rand.Intn(len(candidates))
		for idx1 == idx2 {
			idx2 = rand.Intn(len(candidates))
		}
		reviewers = []string{candidates[idx1], candidates[idx2]}
	}

	pr := &models.PullRequest{
		ID:        id,
		Name:      name,
		AuthorID:  authorID,
		Status:    "OPEN",
		CreatedAt: time.Now(),
		Reviewers: reviewers,
	}

	if err := s.prRepo.CreatePR(ctx, pr); err != nil {
		if errors.Is(err, models.ErrAlreadyExists) {
			s.log.WarnContext(ctx, "pr already exists", "pr_id", id)
		} else {
			s.log.ErrorContext(ctx, "failed to create pr", "err", err)
		}
		return nil, err
	}
	return pr, nil
}

func (s *PullRequestService) MergePullRequest(ctx context.Context, id string) (*models.PullRequest, error) {
	s.log.InfoContext(ctx, "merging PR", "PR_id", id)
	pr, err := s.prRepo.MergePR(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			s.log.WarnContext(ctx, "pr not found", "pr_id", id)
		} else {
			s.log.ErrorContext(ctx, "failed to merge pr", "err", err)
		}
		return nil, err
	}
	return pr, nil
}

func (s *PullRequestService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (string, *models.PullRequest, error) {
	s.log.InfoContext(ctx, "reassigning reviewer", "pr_id", prID, "old_reviewer", oldUserID)

	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			s.log.WarnContext(ctx, "pr not found", "pr_id", prID)
		} else {
			s.log.ErrorContext(ctx, "failed to get pr", "err", err)
		}
		return "", nil, err
	}

	isReviewer := slices.Contains(pr.Reviewers, oldUserID)

	if !isReviewer {
		s.log.WarnContext(ctx, "user is not a reviewer", "pr_id", prID, "user_id", oldUserID)
		return "", nil, models.ErrNotFound
	}

	author, err := s.userRepo.GetUser(ctx, pr.AuthorID)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to get author", "err", err)
		return "", nil, err
	}

	activeUsers, err := s.userRepo.GetActiveUsersByTeam(ctx, author.TeamName)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to get active users", "err", err)
		return "", nil, err
	}

	currentReviewersMap := make(map[string]bool)
	for _, r := range pr.Reviewers {
		currentReviewersMap[r] = true
	}

	candidates := make([]string, 0, len(activeUsers))
	for _, u := range activeUsers {
		if u.ID == pr.AuthorID {
			continue
		}
		if u.ID == oldUserID {
			continue
		}
		if currentReviewersMap[u.ID] {
			continue
		}
		candidates = append(candidates, u.ID)
	}

	if len(candidates) == 0 {
		errNoCandidates := errors.New("no other candidates available")
		s.log.WarnContext(ctx, "no candidates for reassign", "pr_id", prID)
		return "", nil, errNoCandidates
	}

	newReviewerID := candidates[rand.Intn(len(candidates))]

	if err := s.prRepo.ReassignReviewer(ctx, prID, oldUserID, newReviewerID); err != nil {
		s.log.ErrorContext(ctx, "failed to reassign in db", "err", err)
		return "", nil, err
	}

	for i, r := range pr.Reviewers {
		if r == oldUserID {
			pr.Reviewers[i] = newReviewerID
			break
		}
	}

	return newReviewerID, pr, nil
}

func (s *PullRequestService) GetReviewerPRs(ctx context.Context, userID string) ([]*models.PullRequestShort, error) {
	prs, err := s.prRepo.GetByReviewerID(ctx, userID)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to get reviewer prs", "user_id", userID, "err", err)
		return nil, err
	}

	if prs == nil {
		return []*models.PullRequestShort{}, nil
	}

	return prs, nil
}
