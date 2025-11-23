package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/yohnnn/pr_reviewer_assignment_service/internal/app/handlers"
)

type Router struct {
	userHandler  *handlers.UserHandler
	teamHandler  *handlers.TeamHandler
	prHandler    *handlers.PullRequestHandler
	statsHandler *handlers.StatsHandler
}

func NewRouter(
	userHandler *handlers.UserHandler,
	teamHandler *handlers.TeamHandler,
	prHandler *handlers.PullRequestHandler,
	statsHandler *handlers.StatsHandler,
) *Router {
	return &Router{
		userHandler:  userHandler,
		teamHandler:  teamHandler,
		prHandler:    prHandler,
		statsHandler: statsHandler,
	}
}

func (r *Router) InitRoutes() *gin.Engine {

	router := gin.Default()

	router.Use(cors.Default())

	api := router.Group("/")
	{
		users := api.Group("/users")
		{
			users.POST("/setIsActive", r.userHandler.SetActive)

			users.GET("/getReview", r.userHandler.GetReview)

			users.POST("/deactivate", r.userHandler.DeactivateUsers)
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

		router.GET("/stats", r.statsHandler.GetStats)
	}

	return router
}
