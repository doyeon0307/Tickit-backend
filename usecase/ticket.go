package usecase

import (
	"context"

	"github.com/doyeon0307/tickit-backend/domain"
	"github.com/doyeon0307/tickit-backend/models"
)

type ticketUsecase struct {
	ticketRepo domain.TicketRepository
}

func NewTicketUseCase(repo domain.TicketRepository) domain.TicketUsecase {
	return &ticketUsecase{
		ticketRepo: repo,
	}
}

func (u ticketUsecase) GetTicketPreviews(userId string) ([]*models.TicketPreview, error) {
	return u.ticketRepo.GetPreviews(context.Background(), userId)
}

func (u ticketUsecase) GetTicketByID(userId, id string) (*models.Ticket, error) {
	return u.ticketRepo.GetById(context.Background(), userId, id)
}

func (u ticketUsecase) CreateTicket(userId string, ticket *models.Ticket) (string, error) {
	return u.ticketRepo.Create(context.Background(), userId, ticket)
}

func (u ticketUsecase) UpdateTicket(userId, id string, ticket *models.Ticket) error {
	return u.ticketRepo.Update(context.Background(), userId, id, ticket)
}

func (u ticketUsecase) DeleteTicket(userId, id string) error {
	return u.ticketRepo.Delete(context.Background(), userId, id)
}
