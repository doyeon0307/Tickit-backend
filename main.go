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
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	s3Config, err := config.NewS3Config(
		config.AWS_ACCESS_KEY,
		config.AWS_SCRET_KEY,
		"us-east-1",
		"tickit-s3-bucket",
	)
	if err != nil {
		log.Fatal("S3 연결에 실패했습니다")
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("데이터베이스 연결에 실패했습니다")
	}

	ticketRepo := repository.NewTicketRepository(db)
	ticketUsecase := usecase.NewTicketUseCase(ticketRepo)

	scheduleRepo := repository.NewScheduleRepository(db)
	scheduleUsecase := usecase.NewScheduleUsecase(scheduleRepo)

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	handlers := routes.HandlerContainer{
		TicketUsecase:   ticketUsecase,
		ScheduleUsecase: scheduleUsecase,
		UserUsecase:     userUsecase,
		S3Config:        *s3Config,
	}

	router := routes.SetupRouter(handlers)

	router.Run(":7000")
}
