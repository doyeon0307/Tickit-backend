package repository

import (
	"context"
	"time"

	"github.com/doyeon0307/tickit-backend/common"
	"github.com/doyeon0307/tickit-backend/domain"
	"github.com/doyeon0307/tickit-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (m *userRepository) GetById(ctx context.Context, id string) (*models.User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "잘못된 아이디가 추출되었습니다. 토큰을 확인해주세요.",
			Err:     err,
		}
	}
	var profile models.User
	err = m.collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &common.AppError{
				Code:    common.ErrNotFound,
				Message: "사용자가 존재하지 않습니다. 토큰을 확인해주세요.",
				Err:     err,
			}
		}
		return nil, &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}

	return &profile, nil
}

func (m *userRepository) Create(ctx context.Context, user *models.User) (string, error) {
	// oauthId가 이미 존재하는지 확인
	exists, err := m.collection.CountDocuments(ctx, bson.M{"oauthId": user.OAuthId})
	if err != nil {
		return "", &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}

	if exists > 0 {
		return "", &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "이미 존재하는 사용자입니다. 로그인을 시도해주세요.",
			Err:     err,
		}
	}

	result, err := m.collection.InsertOne(ctx, user)
	if err != nil {
		return "", &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}

	user.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return user.Id, nil
}

func (m *userRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "잘못된 아이디가 추출되었습니다. 토큰을 확인해주세요.",
			Err:     err,
		}
	}

	result, err := m.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}

	if result.DeletedCount == 0 {
		return &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "존재하지 않는 사용자는 삭제할 수 없습니다",
			Err:     err,
		}
	}

	return nil
}

func (m *userRepository) GetByOAuthId(ctx context.Context, oauthId string) (*models.User, error) {
	var user models.User
	err := m.collection.FindOne(ctx, bson.M{"oauthId": oauthId}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &common.AppError{
				Code:    common.ErrNotFound,
				Message: "사용자를 찾을 수 없습니다",
				Err:     err,
			}
		}
		return nil, &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}
	return &user, nil
}

func (m *userRepository) SaveRefreshToken(ctx context.Context, userId string, refreshToken string, expiryTime time.Time) error {
	filter := bson.M{"_id": userId}
	update := bson.M{
		"$set": bson.M{
			"refreshToken": refreshToken,
			"tokenExpiry":  expiryTime,
		},
	}

	_, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &common.AppError{
			Code:    common.ErrServer,
			Message: "Refresh Token 저장에 실패했습니다",
			Err:     err,
		}
	}

	return nil
}

func (m *userRepository) GetRefreshToken(ctx context.Context, userId string) (string, error) {
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return "", &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "잘못된 아이디가 추출되었습니다. 토큰을 확인해주세요.",
			Err:     err,
		}
	}

	filter := bson.M{"_id": objId}
	var user models.User

	err = m.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", &common.AppError{
				Code:    common.ErrNotFound,
				Message: "사용자가 존재하지 않습니다. 토큰을 확인해주세요.",
				Err:     err,
			}
		}
		return "", &common.AppError{
			Code:    common.ErrServer,
			Message: "사용자 조회에 실패했습니다",
			Err:     err,
		}
	}

	return user.RefreshToken, nil
}

func (m *userRepository) DeleteUser(ctx context.Context, userId string) error {
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "잘못된 아이디가 추출되었습니다. 토큰을 확인해주세요.",
			Err:     err,
		}
	}

	filter := bson.M{"_id": objId}
	result, err := m.collection.DeleteOne(ctx, filter)
	if err != nil {
		return &common.AppError{
			Code:    common.ErrServer,
			Message: "사용자 삭제에 실패했습니다",
			Err:     err,
		}
	}

	if result.DeletedCount == 0 {
		return &common.AppError{
			Code:    common.ErrNotFound,
			Message: "사용자가 존재하지 않습니다. 토큰을 확인해주세요.",
			Err:     err,
		}
	}

	return nil
}

func (m *userRepository) RemoveRefreshToken(ctx context.Context, userId string) error {
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "잘못된 아이디가 추출되었습니다. 토큰을 확인해주세요.",
			Err:     err,
		}
	}

	filter := bson.M{"_id": objId}
	update := bson.M{
		"$set": bson.M{
			"refreshToken": "",
			"tokenExpiry":  time.Time{},
		},
	}

	result, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &common.AppError{
			Code:    common.ErrServer,
			Message: "Refresh Token 삭제에 실패했습니다",
			Err:     err,
		}
	}

	if result.MatchedCount == 0 {
		return &common.AppError{
			Code:    common.ErrNotFound,
			Message: "사용자가 존재하지 않습니다. 토큰을 확인해주세요.",
			Err:     err,
		}
	}

	return nil
}
