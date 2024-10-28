package main

import (
	"log"
	"net/http"
	"std/github.com/dodo/Tickit-backend/config"
	"std/github.com/dodo/Tickit-backend/repository"
	"std/github.com/dodo/Tickit-backend/routes"
	"std/github.com/dodo/Tickit-backend/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	db, _ := config.ConnectDB()
	if db == nil {
		log.Fatal("데이터베이스 연결에 실패했습니다")
	}

	ticketRepo := repository.NewTicketRepository(db)
	ticketUseCase := usecase.NewTicketUseCase(ticketRepo)

	handlers := routes.HandlerContainer{
		TicketUsecase: ticketUseCase,
	}

	router := routes.SetupRouter(handlers)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	router.Run(":7000")
}
