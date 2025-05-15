package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/google/uuid"
	 "go.mongodb.org/mongo-driver/bson"
	  "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/Prototype-1/freelanceX_message.notification_service/internal/model"
)

type MessageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository(col *mongo.Collection) *MessageRepository {
	return &MessageRepository{collection: col}
}

func (r *MessageRepository) SaveMessage(ctx context.Context, msg *model.Message) (primitive.ObjectID, error) {
    result, err := r.collection.InsertOne(ctx, msg)
    if err != nil {
        return primitive.NilObjectID, err
    }

    insertedID := result.InsertedID.(primitive.ObjectID)
    return insertedID, nil
}

func (r *MessageRepository) GetMessages(ctx context.Context, fromID, toID uuid.UUID, limit, offset int) ([]model.Message, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from_user_id": fromID, "to_user_id": toID},
			{"from_user_id": toID, "to_user_id": fromID},
		},
		"is_deleted": false,
	}

	opts := options.Find().
		SetSort(bson.D{primitive.E{Key: "sent_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []model.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepository) CreateMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
    res, err := r.collection.InsertOne(ctx, message)
    if err != nil {
        return nil, err
    }
    message.ID = res.InsertedID.(primitive.ObjectID)
    return message, nil
}