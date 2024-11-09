package domain

import (
	"context"

	"github.com/doyeon0307/tickit-backend/models"
)

type ScheduleRepository interface {
	GetPreviewsForTicket(ctx context.Context, userId, date string) ([]*models.Schedule, error)
	GetPreviewsForCalendar(ctx context.Context, userId, startDate, endDate string) ([]*models.Schedule, error)
	GetById(ctx context.Context, userId, id string) (*models.Schedule, error)
	Create(ctx context.Context, schedule *models.Schedule) (string, error)
	Update(ctx context.Context, userId, id string, schedule *models.Schedule) error
	Delete(ctx context.Context, userId, id string) error
}
