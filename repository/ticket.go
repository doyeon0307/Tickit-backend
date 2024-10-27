package repository

import (
	"context"
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
	opts := options.Find().SetProjection(bson.M{
		"_id":   1,
		"image": 1,
	})

	cursor, err := m.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []*models.TicketPreview
	if err = cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}

	return tickets, nil
}

func (m *ticketRepository) GetById(ctx context.Context, id string) (*models.Ticket, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var ticket models.Ticket
	err = m.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &ticket, nil
}

func (m *ticketRepository) Create(ctx context.Context, ticket *models.Ticket) error {
	result, err := m.collection.InsertOne(ctx, ticket)
	if err != nil {
		return err
	}

	ticket.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (m *ticketRepository) Update(ctx context.Context, id string, ticket *models.Ticket) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
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
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (m *ticketRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := m.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
