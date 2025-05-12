package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
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
