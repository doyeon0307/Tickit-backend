package domain

import "github.com/doyeon0307/tickit-backend/models"

type TicketUsecase interface {
	GetTicketPreviews(userId string) ([]*models.TicketPreview, error)
	GetTicketByID(userId, id string) (*models.Ticket, error)
	CreateTicket(userId string, ticket *models.Ticket) (string, error)
	UpdateTicket(userId, id string, ticket *models.Ticket) error
	DeleteTicket(userId, id string) error
}
