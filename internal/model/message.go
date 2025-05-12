package model

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/google/uuid"
    "time"
)

type MessageChannel string
const (
    EmailChannel MessageChannel = "email"
    InAppChannel MessageChannel = "in-app"
)

type MessageStatus string
const (
    SentStatus      MessageStatus = "sent"
    DeliveredStatus MessageStatus = "delivered"
    ReadStatus      MessageStatus = "read"
)

type Message struct {
    ID            primitive.ObjectID `bson:"_id,omitempty"`
    FromUserID    uuid.UUID          `bson:"from_user_id"`
    ToUserID      uuid.UUID          `bson:"to_user_id"`
    ProjectID     uuid.UUID          `bson:"project_id"`
    MessageText   string             `bson:"message"`
    Channel       MessageChannel     `bson:"channel"`
    SentAt        time.Time          `bson:"sent_at"`
    Read          bool               `bson:"read"`
    Attachments   []string           `bson:"attachments"`
    Status        MessageStatus      `bson:"status"`
    IsDeleted     bool               `bson:"is_deleted"`
}
