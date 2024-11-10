package domain

import (
	"github.com/doyeon0307/tickit-backend/dto"
	"github.com/doyeon0307/tickit-backend/models"
)

type UserUsecase interface {
	GetProfile(id string) (*dto.KakaoProfile, error)
	CreateUser(idToken string, accessToken string) (string, error)
	DeleteUser(id string) error
	GetUserByOAuthId(oauthId string) (*models.User, error)
}
