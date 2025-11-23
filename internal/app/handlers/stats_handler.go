package handlers

import (
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/services"
)

type StatsHandler struct {
	statsService services.StatsServiceInterface
	log          *slog.Logger
}

func NewStatsHandler(
	statsService services.StatsServiceInterface,
	log *slog.Logger,
) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
		log:          log,
	}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	stats, err := h.statsService.GetTopReviewers(c.Request.Context())
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, ErrCodeStatsError, "Failed to get stats")
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}
