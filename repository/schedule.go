package repository

import (
	"context"

	"github.com/doyeon0307/tickit-backend/common"
	"github.com/doyeon0307/tickit-backend/domain"
	"github.com/doyeon0307/tickit-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type scheduleRepository struct {
	collection *mongo.Collection
}

func NewScheduleRepository(db *mongo.Database) domain.ScheduleRepository {
	return &scheduleRepository{
		collection: db.Collection("schedules"),
	}
}

func (m *scheduleRepository) GetPreviewsForTicket(ctx context.Context, userId, date string) ([]*models.Schedule, error) {
	previews := make([]*models.Schedule, 0)

	filter := bson.M{
		"userId": userId,
		"date": bson.M{
			"$lte": date,
		},
	}

	opts := options.Find().SetSort(bson.M{"date": -1})

	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &previews); err != nil {
		return nil, &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}

	if previews == nil {
		previews = make([]*models.Schedule, 0)
	}

	return previews, nil
}

func (m *scheduleRepository) GetPreviewsForCalendar(ctx context.Context, userId, startDate, endDate string) ([]*models.Schedule, error) {
	previews := make([]*models.Schedule, 0)

	filter := bson.M{
		"userId": userId,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	opts := options.Find().SetSort(bson.M{"date": 1})

	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &previews); err != nil {
		return nil, &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}

	if previews == nil {
		previews = make([]*models.Schedule, 0)
	}

	return previews, nil
}

func (m *scheduleRepository) GetById(ctx context.Context, userId, id string) (*models.Schedule, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "아이디 형식이 잘못되었습니다",
			Err:     err,
		}
	}

	var schedule models.Schedule
	err = m.collection.FindOne(ctx, bson.M{"_id": objID, "userId": userId}).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &common.AppError{
				Code:    common.ErrNotFound,
				Message: "일정이 존재하지 않습니다. 아이디를 확인해주세요.",
				Err:     err,
			}
		}
		return nil, &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}

	return &schedule, nil
}

func (m *scheduleRepository) Create(ctx context.Context, schedule *models.Schedule) (string, error) {
	// Schedule에 이미 UserId가 설정되어 있다고 가정
	result, err := m.collection.InsertOne(ctx, schedule)
	if err != nil {
		return "", &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}

	schedule.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return schedule.Id, nil
}

func (m *scheduleRepository) Update(ctx context.Context, userId, id string, schedule *models.Schedule) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "아이디 형식이 잘못되었습니다",
			Err:     err,
		}
	}

	// userId 검증을 위한 필터 추가
	filter := bson.M{
		"_id":    objID,
		"userId": userId,
	}

	update := bson.M{
		"$set": bson.M{
			"date":      schedule.Date,
			"title":     schedule.Title,
			"number":    schedule.Number,
			"image":     schedule.Image,
			"thumbnail": schedule.Thumbnail,
			"location":  schedule.Location,
			"time":      schedule.Time,
			"seat":      schedule.Seat,
			"casting":   schedule.Casting,
			"company":   schedule.Company,
			"link":      schedule.Link,
			"memo":      schedule.Memo,
			"userId":    userId, // userId도 함께 업데이트
		},
	}

	result, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}

	if result.MatchedCount == 0 {
		return &common.AppError{
			Code:    common.ErrNotFound,
			Message: "존재하지 않는 일정이거나 수정 권한이 없습니다",
			Err:     err,
		}
	}

	return nil
}

func (m *scheduleRepository) Delete(ctx context.Context, userId, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "아이디 형식이 잘못되었습니다",
			Err:     err,
		}
	}

	// userId 검증을 위한 필터 추가
	filter := bson.M{
		"_id":    objID,
		"userId": userId,
	}

	result, err := m.collection.DeleteOne(ctx, filter)
	if err != nil {
		return &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}

	if result.DeletedCount == 0 {
		return &common.AppError{
			Code:    common.ErrNotFound,
			Message: "존재하지 않는 일정이거나 삭제 권한이 없습니다",
			Err:     err,
		}
	}

	return nil
}
