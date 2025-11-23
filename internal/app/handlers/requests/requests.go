package requests

type CreateTeamReq struct {
	Name    string          `json:"team_name" binding:"required"`
	Members []TeamMemberReq `json:"members"   binding:"required,dive"`
}

type TeamMemberReq struct {
	UserID   string `json:"user_id"   binding:"required"`
	UserName string `json:"username"  binding:"required"`
	IsActive bool   `json:"is_active"`
}

type SetActiveReq struct {
	UserID   string `json:"user_id"   binding:"required"`
	IsActive bool   `json:"is_active"`
}

type CreatePRReq struct {
	ID       string `json:"pull_request_id"   binding:"required"`
	Name     string `json:"pull_request_name" binding:"required"`
	AuthorID string `json:"author_id"         binding:"required"`
}

type MergePRReq struct {
	ID string `json:"pull_request_id" binding:"required"`
}

type ReassignPRReq struct {
	ID        string `json:"pull_request_id" binding:"required"`
	OldUserID string `json:"old_user_id"     binding:"required"`
}

type DeactivateRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1"`
}
