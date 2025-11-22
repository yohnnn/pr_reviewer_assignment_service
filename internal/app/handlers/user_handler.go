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

type UserHandler struct {
	userService *services.UserService
	prService   *services.PullRequestService
	log         *slog.Logger
}

func NewUserHandler(
	userService *services.UserService,
	prService *services.PullRequestService,
	log *slog.Logger,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		prService:   prService,
		log:         log,
	}
}

func (h *UserHandler) SetActive(c *gin.Context) {
	var req requests.SetActiveReq
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WarnContext(ctx, "failed to validate or decode request", "err", err)
		writeErrorResponse(c, http.StatusBadRequest, ErrCodeInvalidFormat, "Invalid request payload")
		return
	}
	user, err := h.userService.SetUserIsActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			writeErrorResponse(c, http.StatusNotFound, ErrCodeNotFound, "User not found")
		} else {
			writeErrorResponse(c, http.StatusInternalServerError, ErrCodeInternal, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")
	ctx := c.Request.Context()

	if userID == "" {
		h.log.WarnContext(ctx, "missing user_id query param")
		writeErrorResponse(c, http.StatusBadRequest, ErrCodeInvalidFormat, "user_id query parameter is required")
		return
	}

	prShort, err := h.prService.GetReviewerPRs(ctx, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			writeErrorResponse(c, http.StatusNotFound, ErrCodeNotFound, "User not found")
		} else {
			writeErrorResponse(c, http.StatusInternalServerError, ErrCodeInternal, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"pull_requests": prShort,
	})
}
