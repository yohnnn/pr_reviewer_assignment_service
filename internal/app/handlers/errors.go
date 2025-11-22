package handlers

import (
	"github.com/gin-gonic/gin"
)

const (
	ErrCodeInvalidFormat = "INVALID_FORMAT"
	ErrCodeInternal      = "INTERNAL_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeTeamExists    = "TEAM_EXISTS"
	ErrCodePRExists      = "PR_EXISTS"
	ErrCodeNoCandidates  = "NO_CANDIDATES"
	ErrCodePRMerged      = "PR_MERGED"
	ErrCodeNotAssigned   = "NOT_ASSIGNED"
	ErrCodeNoCandidate   = "NO_CANDIDATE"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func writeErrorResponse(c *gin.Context, status int, code string, msg string) {
	var errResp ErrorResponse
	errResp.Error.Code = code
	errResp.Error.Message = msg
	c.JSON(status, errResp)
}
