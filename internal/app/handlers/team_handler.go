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
	TeamService *services.TeamService
	log         *slog.Logger
}

func NewTeamHandler(
	TeamService *services.TeamService,
	log *slog.Logger,
) *TeamHandler {
	return &TeamHandler{
		TeamService: TeamService,
		log:         log,
	}
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	ctx := c.Request.Context()

	if teamName == "" {
		h.log.WarnContext(ctx, "missing team_name query param", "team_name", teamName)
		writeErrorResponse(c, http.StatusBadRequest, ErrCodeInvalidFormat, "team_name query parameter is required")
		return
	}

	team, err := h.TeamService.GetTeam(ctx, teamName)
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
	createdTeam, err := h.TeamService.CreateTeam(ctx, team)
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
