package service

import (
		"time"
		"context"
	    "github.com/google/uuid"
		 "google.golang.org/protobuf/types/known/timestamppb"
	    "google.golang.org/grpc/status"
    	"google.golang.org/grpc/codes"
		"go.mongodb.org/mongo-driver/bson/primitive"
		pb "github.com/Prototype-1/freelanceX_message.notification_service/proto"
		"github.com/Prototype-1/freelanceX_message.notification_service/internal/model"
)

type MessageRepository interface {
    SaveMessage(ctx context.Context, msg *model.Message) (primitive.ObjectID, error)
}

type MessageService struct {
    repo MessageRepository
}

func NewMessageService(repo MessageRepository) *MessageService {
    return &MessageService{repo: repo}
}

func (s *MessageService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
    fromUUID, err := uuid.Parse(req.GetFromUserId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid from_user_id")
    }
    toUUID, err := uuid.Parse(req.GetToUserId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid to_user_id")
    }
    projectUUID, err := uuid.Parse(req.GetProjectId())
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid project_id")
    }

    now := time.Now()

    msg := &model.Message{
        FromUserID:  fromUUID,
        ToUserID:    toUUID,
        ProjectID:   projectUUID,
        MessageText: req.GetMessage(),
        Channel:     model.InAppChannel,
        SentAt:      now,
        Read:        false,
        Status:      model.SentStatus,
        IsDeleted:   false,
        Attachments: req.GetAttachments(),
    }

    insertedID, err := s.repo.SaveMessage(ctx, msg)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to save message: %v", err)
    }

    return &pb.SendMessageResponse{
        MessageId: insertedID.Hex(),
        SentAt:    timestamppb.New(now),
    }, nil
}
