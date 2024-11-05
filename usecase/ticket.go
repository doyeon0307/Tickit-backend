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

func (u ticketUsecase) GetTicketPreviews() ([]*models.TicketPreview, error) {
	return u.ticketRepo.GetPreviews(context.Background())
}

func (u ticketUsecase) GetTicketByID(id string) (*models.Ticket, error) {
	return u.ticketRepo.GetById(context.Background(), id)
}

func (u ticketUsecase) CreateTicket(ticket *models.Ticket) (string, error) {
	return u.ticketRepo.Create(context.Background(), ticket)
}

func (u ticketUsecase) UpdateTicket(id string, ticket *models.Ticket) error {
	return u.ticketRepo.Update(context.Background(), id, ticket)
}

func (u ticketUsecase) DeleteTicket(id string) error {
	return u.ticketRepo.Delete(context.Background(), id)
}
