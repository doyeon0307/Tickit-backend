package domain

import (
	"context"

	"github.com/doyeon0307/tickit-backend/models"
)

type TicketRepository interface {
	GetPreviews(ctx context.Context) ([]*models.TicketPreview, error)
	GetById(ctx context.Context, id string) (*models.Ticket, error)
	Create(ctx context.Context, ticket *models.Ticket) (string, error)
	Update(stx context.Context, id string, ticket *models.Ticket) error
	Delete(ctx context.Context, id string) error
}
