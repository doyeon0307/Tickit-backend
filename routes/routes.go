package routes

import (
	"net/http"
	"std/github.com/dodo/Tickit-backend/domain"
	"std/github.com/dodo/Tickit-backend/handler"

	"github.com/gin-gonic/gin"
)

type HandlerContainer struct {
	TicketUsecase domain.TicketUsecase
}

func SetupRouter(handlers HandlerContainer) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthCheck)

		handler.NewTicketHandler(v1, handlers.TicketUsecase)
	}

	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
