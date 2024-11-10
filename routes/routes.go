package routes

import (
	"net/http"

	"github.com/doyeon0307/tickit-backend/config"
	"github.com/doyeon0307/tickit-backend/domain"
	"github.com/doyeon0307/tickit-backend/handler"
	"github.com/doyeon0307/tickit-backend/service"

	"github.com/gin-gonic/gin"
)

type HandlerContainer struct {
	TicketUsecase   domain.TicketUsecase
	ScheduleUsecase domain.ScheduleUsecase
	UserUsecase     domain.UserUsecase
	S3Config        config.S3Config
}

func SetupRouter(handlers HandlerContainer) *gin.Engine {
	router := gin.Default()

	config.SetUpSwagger(router)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	v1 := router.Group("/api")
	{
		v1.GET("/health", healthCheck)

		handler.NewUserHandler(v1, handlers.UserUsecase)

		authorized := v1.Group("")
		authorized.Use(service.AuthMiddleware())
		{
			handler.NewTicketHandler(authorized, handlers.TicketUsecase)
			handler.NewScheduleHandler(authorized, handlers.ScheduleUsecase)
			handler.NewS3Handler(authorized, &handlers.S3Config)
		}
	}

	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
