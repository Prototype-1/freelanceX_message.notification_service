package client

import (
    "context"
    "fmt"
    pb "github.com/Prototype-1/freelanceX_message.notification_service/proto/user_service"
)

type UserServiceClient struct {
    Grpc pb.AuthServiceClient
}

func NewUserServiceClient(grpcClient pb.AuthServiceClient) *UserServiceClient {
    return &UserServiceClient{Grpc: grpcClient}
}

func (uc *UserServiceClient) GetUserEmail(ctx context.Context, userID string) (string, error) {
    resp, err := uc.Grpc.GetMe(ctx, &pb.SessionRequest{
        UserId: userID,
    })
    if err != nil {
        return "", fmt.Errorf("failed to fetch user email: %w", err)
    }
      if resp == nil || resp.Email == "" {
        return "", fmt.Errorf("user not found or email is empty")
    }
    return resp.Email, nil
}
