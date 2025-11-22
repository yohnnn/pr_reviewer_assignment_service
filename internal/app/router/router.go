package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/app/handlers"
)

type Router struct {
	userHandler *handlers.UserHandler
	teamHandler *handlers.TeamHandler
	prHandler   *handlers.PullRequestHandler
}

func NewRouter(
	userHandler *handlers.UserHandler,
	teamHandler *handlers.TeamHandler,
	prHandler *handlers.PullRequestHandler,
) *Router {
	return &Router{
		userHandler: userHandler,
		teamHandler: teamHandler,
		prHandler:   prHandler,
	}
}

func (r *Router) InitRoutes() *gin.Engine {

	router := gin.Default()

	api := router.Group("/")
	{
		users := api.Group("/users")
		{
			users.POST("/setIsActive", r.userHandler.SetActive)

			users.GET("/getReview", r.userHandler.GetReview)
		}

		teams := api.Group("/team")
		{
			teams.POST("/add", r.teamHandler.AddTeam)

			teams.GET("/get", r.teamHandler.GetTeam)
		}

		prs := api.Group("/pullRequest")
		{
			prs.POST("/create", r.prHandler.CreatePullRequest)

			prs.POST("/reassign", r.prHandler.ReassignPR)

			prs.POST("/merge", r.prHandler.MergePullRequest)
		}
	}

	return router
}
