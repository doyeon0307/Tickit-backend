package main

import (
	"log"

	"github.com/doyeon0307/tickit-backend/config"
	"github.com/doyeon0307/tickit-backend/repository"
	"github.com/doyeon0307/tickit-backend/routes"
	"github.com/doyeon0307/tickit-backend/usecase"
)

// @title Tickit!
// @version 1.0
// @description 소중한 기억을 나만의 티켓북에 기록하세요
// @host localhost:7000
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

	router.Run(":7000")
}
