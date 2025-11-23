package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/app/handlers/requests"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/models"
	"github.com/yohnnn/pr_reviewer_assignment_service/internal/services"
)

type PullRequestHandler struct {
	pullRequestService services.PullRequestServiceInterface
	log                *slog.Logger
}

func NewPullRequestHandler(
	pullRequestService services.PullRequestServiceInterface,
	log *slog.Logger,
) *PullRequestHandler {
	return &PullRequestHandler{
		pullRequestService: pullRequestService,
		log:                log,
	}
}

// POST /pullRequest/create
func (h *PullRequestHandler) CreatePullRequest(c *gin.Context) {
	var req requests.CreatePRReq
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WarnContext(ctx, "failed to validate or decode request", "err", err)
		writeErrorResponse(c, http.StatusBadRequest, ErrCodeInvalidFormat, "Invalid request payload")
		return
	}

	pr, err := h.pullRequestService.CreatePullRequest(ctx, req.ID, req.Name, req.AuthorID)
	if err != nil {
		if errors.Is(err, models.ErrAlreadyExists) {
			writeErrorResponse(c, http.StatusConflict, ErrCodePRExists, "PR id already exists")
		} else if errors.Is(err, models.ErrNotFound) {
			writeErrorResponse(c, http.StatusNotFound, ErrCodeNotFound, "Author not found")
		} else {
			writeErrorResponse(c, http.StatusInternalServerError, ErrCodeInternal, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"pr": pr})
}

// POST /pullRequest/merge
func (h *PullRequestHandler) MergePullRequest(c *gin.Context) {
	var req requests.MergePRReq
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WarnContext(ctx, "failed to validate or decode request", "err", err)
		writeErrorResponse(c, http.StatusBadRequest, ErrCodeInvalidFormat, "Invalid request payload")
		return
	}

	pr, err := h.pullRequestService.MergePullRequest(ctx, req.ID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			writeErrorResponse(c, http.StatusNotFound, ErrCodeNotFound, "PR not found")
		} else {
			writeErrorResponse(c, http.StatusInternalServerError, ErrCodeInternal, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"pr": pr})
}

// POST /pullRequest/reassign
func (h *PullRequestHandler) ReassignPR(c *gin.Context) {
	var req requests.ReassignPRReq
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WarnContext(ctx, "failed to validate request", "err", err)
		writeErrorResponse(c, http.StatusBadRequest, ErrCodeInvalidFormat, "Invalid request payload")
		return
	}

	pr, newReviewerID, err := h.pullRequestService.ReassignReviewer(ctx, req.ID, req.OldUserID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			writeErrorResponse(c, http.StatusNotFound, ErrCodeNotFound, "PR or User not found")
		} else if errors.Is(err, models.ErrPRMerged) {
			writeErrorResponse(c, http.StatusConflict, ErrCodePRMerged, "cannot reassign on merged PR")
		} else if errors.Is(err, models.ErrNotAssigned) {
			writeErrorResponse(c, http.StatusConflict, ErrCodeNotAssigned, "reviewer is not assigned to this PR")
		} else if errors.Is(err, models.ErrNoCandidates) {
			writeErrorResponse(c, http.StatusConflict, ErrCodeNoCandidate, "no active replacement candidate in team")
		} else {
			writeErrorResponse(c, http.StatusInternalServerError, ErrCodeInternal, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr":          pr,
		"replaced_by": newReviewerID,
	})
}
