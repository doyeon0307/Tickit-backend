package domain

import "github.com/doyeon0307/tickit-backend/models"

type TicketUsecase interface {
	GetTicketPreviews() ([]*models.TicketPreview, error)
	GetTicketByID(id string) (*models.Ticket, error)
	CreateTicket(ticket *models.Ticket) (string, error)
	UpdateTicket(id string, ticket *models.Ticket) error
	DeleteTicket(id string) error
}
