package service

import (
		"time"
		"context"
        "fmt"
	    "github.com/google/uuid"
		 "google.golang.org/protobuf/types/known/timestamppb"
	    "google.golang.org/grpc/status"
    	"google.golang.org/grpc/codes"
		"go.mongodb.org/mongo-driver/bson/primitive"
		pb "github.com/Prototype-1/freelanceX_message.notification_service/proto"
		 "github.com/Prototype-1/freelanceX_message.notification_service/email"
		 "github.com/Prototype-1/freelanceX_message.notification_service/config"
		 "github.com/Prototype-1/freelanceX_message.notification_service/internal/client"
		"github.com/Prototype-1/freelanceX_message.notification_service/internal/model"
)

type MessageRepository interface {
    SaveMessage(ctx context.Context, msg *model.Message) (primitive.ObjectID, error)
    GetMessages(ctx context.Context, fromID, toID uuid.UUID, limit, offset int) ([]model.Message, error)
}

type MessageService struct {
     pb.UnimplementedMessageServiceServer
    repo MessageRepository
     smtpCfg email.SMTPConfig
     userClient *client.UserServiceClient 
}

func NewMessageService(repo MessageRepository, smtpCfg email.SMTPConfig, userClient *client.UserServiceClient) *MessageService {
    return &MessageService{
        repo:    repo,
        smtpCfg: smtpCfg,
         userClient: userClient,
    }
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

senderEmail, err := s.userClient.GetUserEmail(ctx, req.GetFromUserId())
if err != nil {
    return nil, status.Errorf(codes.InvalidArgument, "invalid from_user_id: %v", err)
}

recipientEmail, err := s.userClient.GetUserEmail(ctx, req.GetToUserId())
if err != nil {
    return nil, status.Errorf(codes.InvalidArgument, "invalid to_user_id: %v", err)
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

      isOnline, err := config.IsUserOnline(ctx, req.GetToUserId())
    if err != nil {
        fmt.Printf("Error checking user online status: %v\n", err)
    }
    if !isOnline {
    subject := "You've received a new message"
    body := fmt.Sprintf("You have a new message from user %s: %s", senderEmail, req.GetMessage())

    err = email.SendMail(s.smtpCfg, recipientEmail, subject, body)
    if err != nil {
        fmt.Printf("Failed to send email: %v\n", err)
    } else {
        fmt.Println("Notification email sent successfully!")
    }
}
    return &pb.SendMessageResponse{
        MessageId: insertedID.Hex(),
        SentAt:    timestamppb.New(now),
    }, nil
}

func (s *MessageService) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
    fromID, err := uuid.Parse(req.SenderId)
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid sender ID: %v", err)
    }
    
    toID, err := uuid.Parse(req.ReceiverId)
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid receiver ID: %v", err)
    }
    
    messages, err := s.repo.GetMessages(ctx, fromID, toID, int(req.Limit), int(req.Offset))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to fetch messages: %v", err)
    }

    pbMessages := make([]*pb.Message, 0, len(messages))
    for _, msg := range messages {
        pbMessages = append(pbMessages, &pb.Message{
            Id:         msg.ID.Hex(),
            SenderId:   msg.FromUserID.String(),
            ReceiverId: msg.ToUserID.String(),
            Content:    msg.MessageText,
            Timestamp:  msg.SentAt.Format(time.RFC3339),
        })
    }
    
    return &pb.GetMessagesResponse{Messages: pbMessages}, nil
}