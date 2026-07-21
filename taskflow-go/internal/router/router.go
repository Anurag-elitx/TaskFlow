package router

import (
	"taskflow/internal/handlers"
	"taskflow/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	teamHandler *handlers.TeamHandler,
	taskHandler *handlers.TaskHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetMe)
				users.PATCH("/me", userHandler.UpdateMe)
			}

			teams := protected.Group("/teams")
			{
				teams.POST("", teamHandler.CreateTeam)
				teams.GET("", teamHandler.GetTeams)
				teams.POST("/:id/members", teamHandler.AddMember)
				teams.POST("/:teamId/tasks", taskHandler.CreateTask)
				teams.GET("/:teamId/tasks", taskHandler.GetTasks)
			}

			tasks := protected.Group("/tasks")
			{
				tasks.GET("/:id", taskHandler.GetTask)
				tasks.PATCH("/:id", taskHandler.UpdateTask)
				tasks.DELETE("/:id", taskHandler.DeleteTask)
			}
		}
	}

	return r
}
