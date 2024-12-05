package domain

import (
	"context"
	"time"

	"github.com/doyeon0307/tickit-backend/models"
)

type UserRepository interface {
	GetById(ctx context.Context, id string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (string, error)
	Delete(ctx context.Context, id string) error
	GetByOAuthId(ctx context.Context, oauthId string) (*models.User, error)
	SaveRefreshToken(ctx context.Context, userId string, refreshToken string, expiryTime time.Time) error
	GetRefreshToken(ctx context.Context, userId string) (string, error)
}
