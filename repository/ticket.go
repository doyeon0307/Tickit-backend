package repository

import (
	"context"
	"std/github.com/dodo/Tickit-backend/common"
	"std/github.com/dodo/Tickit-backend/domain"
	"std/github.com/dodo/Tickit-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ticketRepository struct {
	collection *mongo.Collection
}

func NewTicketRepository(db *mongo.Database) domain.TicketRepository {
	return &ticketRepository{
		collection: db.Collection("tickets"),
	}
}

func (m *ticketRepository) GetPreviews(ctx context.Context) ([]*models.TicketPreview, error) {
	previews := make([]*models.TicketPreview, 0)

	opts := options.Find().SetProjection(bson.M{
		"_id":   1,
		"image": 1,
	})

	cursor, err := m.collection.Find(ctx, bson.M{}, opts)
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
		previews = make([]*models.TicketPreview, 0)
	}

	return previews, nil
}

func (m *ticketRepository) GetById(ctx context.Context, id string) (*models.Ticket, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "아이디 형식이 잘못되었습니다",
			Err:     err,
		}
	}

	var ticket models.Ticket
	err = m.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &common.AppError{
				Code:    common.ErrNotFound,
				Message: "티켓이 존재하지 않습니다. 아이디를 확인해주세요.",
				Err:     err,
			}
		}
		return nil, &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}

	return &ticket, nil
}

func (m *ticketRepository) Create(ctx context.Context, ticket *models.Ticket) (string, error) {
	result, err := m.collection.InsertOne(ctx, ticket)
	if err != nil {
		return "", &common.AppError{
			Code:    common.ErrServer,
			Message: "데이터베이스 오류가 발생했습니다",
			Err:     err,
		}
	}

	ticket.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return ticket.Id, nil
}

func (m *ticketRepository) Update(ctx context.Context, id string, ticket *models.Ticket) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "아이디 형식이 잘못되었습니다",
			Err:     err,
		}
	}

	update := bson.M{
		"$set": bson.M{
			"image":           ticket.Image,
			"title":           ticket.Title,
			"location":        ticket.Location,
			"datetime":        ticket.Datetime,
			"backgroundColor": ticket.BackgroundColor,
			"foregroundColor": ticket.ForegroundColor,
			"fields":          ticket.Fields,
		},
	}

	result, err := m.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
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
			Message: "존재하지 않는 티켓은 수정할 수 없습니다",
			Err:     err,
		}
	}

	return nil
}

func (m *ticketRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &common.AppError{
			Code:    common.ErrBadRequest,
			Message: "아이디 형식이 잘못되었습니다",
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
			Message: "존재하지 않는 티켓은 삭제할 수 없습니다",
			Err:     err,
		}
	}

	return nil
}
