package domain

import (
	"time"

	"github.com/doyeon0307/tickit-backend/dto"
	"github.com/doyeon0307/tickit-backend/models"
)

type UserUsecase interface {
	GetProfile(id string) (*dto.KakaoProfile, error)
	CreateUser(idToken string, accessToken string) (string, error)
	DeleteUser(id string) error
	GetUserByOAuthId(oauthId string) (*models.User, error)
	SaveRefreshToken(userId string, refreshToken string, expiryTime time.Time) error
	ValidateStoredRefreshToken(userId string, refreshToken string) (bool, error)
	WithdrawUser(userId string) error
	Logout(userId string) error
}
