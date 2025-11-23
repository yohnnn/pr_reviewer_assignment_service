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

type TeamHandler struct {
	teamService services.TeamServiceInterface
	log         *slog.Logger
}

func NewTeamHandler(
	teamService services.TeamServiceInterface,
	log *slog.Logger,
) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
		log:         log,
	}
}

// GET /team/get
func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	ctx := c.Request.Context()

	if teamName == "" {
		h.log.WarnContext(ctx, "missing team_name query param", "team_name", teamName)
		writeErrorResponse(c, http.StatusBadRequest, ErrCodeInvalidFormat, "team_name query parameter is required")
		return
	}

	team, err := h.teamService.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			writeErrorResponse(c, http.StatusNotFound, ErrCodeNotFound, "Team not found")
		} else {
			writeErrorResponse(c, http.StatusInternalServerError, ErrCodeInternal, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, team)
}

// POST /team/get
func (h *TeamHandler) AddTeam(c *gin.Context) {
	var req requests.CreateTeamReq
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WarnContext(ctx, "failed to validate or decode request", "err", err)
		writeErrorResponse(c, http.StatusBadRequest, ErrCodeInvalidFormat, "Invalid request payload")
		return
	}

	members := make([]models.TeamMember, len(req.Members))

	for i, m := range req.Members {
		members[i] = models.TeamMember{
			UserID:   m.UserID,
			UserName: m.UserName,
			IsActive: m.IsActive,
		}
	}

	team := &models.Team{
		Name:    req.Name,
		Members: members,
	}
	createdTeam, err := h.teamService.CreateTeam(ctx, team)
	if err != nil {
		if errors.Is(err, models.ErrAlreadyExists) {
			writeErrorResponse(c, http.StatusBadRequest, ErrCodeTeamExists, "team_name already exists")
		} else {
			writeErrorResponse(c, http.StatusInternalServerError, ErrCodeInternal, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"team": createdTeam})

}
