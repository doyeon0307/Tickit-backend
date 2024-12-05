package service

import (
	"os"
	"time"

	"github.com/doyeon0307/tickit-backend/common"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

var secretKey = os.Getenv("JWT_SECRET_KEY")

const (
	accessTokenDuration  = 30 * time.Minute
	refreshTokenDuration = 30 * 24 * time.Hour
)

func GenerateAccessToken(userId string) (string, error) {
	claims := Claims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", &common.AppError{
			Code:    common.ErrServer,
			Message: "토큰 생성에 실패했습니다",
			Err:     err,
		}
	}

	return accessToken, nil
}

func GenerateRefreshToken(userId string) (string, time.Time, error) {
	expiryTime := time.Now().Add(refreshTokenDuration)
	claims := Claims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", time.Time{}, &common.AppError{
			Code:    common.ErrServer,
			Message: "토큰 생성에 실패했습니다",
			Err:     err,
		}
	}

	return refreshToken, expiryTime, nil
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", &common.AppError{
				Code:    common.ErrUnauthorized,
				Message: "유효하지 않은 토큰입니다",
				Err:     err,
			}
		}
		ve, ok := err.(*jwt.ValidationError)
		if ok && ve.Errors == jwt.ValidationErrorExpired {
			return "", &common.AppError{
				Code:    common.ErrUnauthorized,
				Message: "토큰이 만료되었습니다",
				Err:     err,
			}
		}
		return "", &common.AppError{
			Code:    common.ErrUnauthorized,
			Message: "토큰 검증에 실패했습니다",
			Err:     err,
		}
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserId, nil
	}

	return "", &common.AppError{
		Code:    common.ErrUnauthorized,
		Message: "토큰 검증에 실패했습니다",
	}
}
