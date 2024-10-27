package domain

import "std/github.com/dodo/Tickit-backend/models"

type TicketUsecase interface {
	GetTicketPreviews() ([]*models.TicketPreview, error)
	GetTicketByID(id string) (*models.Ticket, error)
	CreateTicket(ticket *models.Ticket) error
	UpdateTicket(id string, ticket *models.Ticket) error
	DeleteTicket(id string) error
}
