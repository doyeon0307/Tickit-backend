package main

import (
	"context"
	"log"
	"net/http"
	"std/github.com/dodo/Tickit-backend/repository"
	"std/github.com/dodo/Tickit-backend/repository/mongodb"
	"std/github.com/dodo/Tickit-backend/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	client, err := mongodb.NewMongoDBClient("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("tickit")

	ticketRepo := repository.NewTicketRepository(db)
	ticketUseCase := usecase.NewTicketUseCase(ticketRepo)

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})
	router.Run(":7000")
}
