package models

import "time"

type OAuthType string

const (
	OAuthKakao  OAuthType = "KAKAO"
	OAuthGoogle OAuthType = "GOOGLE"
)

type User struct {
	Id      string `json:"id" bson:"_id,omitempty"`
	OAuthId string `json:"oauthId" bson:"oauthId"`
	// OAuthType OAuthType          `json:"oauthType" bson:"oauthType"`
	Name string `json:"name" bson:"name"`
	// Email       string    `json:"email" bson:"email"`
	// CreatedAt time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	RefreshToken string    `json:"refreshToken" bson:"refreshToken"`
	TokenExpiry  time.Time `json:"tokenExpiry" bson:"tokenExpiry"`
}
