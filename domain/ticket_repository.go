package domain

import (
	"context"

	"github.com/doyeon0307/tickit-backend/models"
)

type TicketRepository interface {
	GetPreviews(ctx context.Context, userId string) ([]*models.Ticket, error)
	GetById(ctx context.Context, userId, id string) (*models.Ticket, error)
	Create(ctx context.Context, userId string, ticket *models.Ticket) (string, error)
	Update(stx context.Context, userId, id string, ticket *models.Ticket) error
	Delete(ctx context.Context, id string) error
}
