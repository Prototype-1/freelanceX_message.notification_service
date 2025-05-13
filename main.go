package main

import (
"log"
"context"
"net"
"fmt"
"github.com/Prototype-1/freelanceX_message.notification_service/config"
proto "github.com/Prototype-1/freelanceX_message.notification_service/proto"
"github.com/Prototype-1/freelanceX_message.notification_service/internal/service"
"github.com/Prototype-1/freelanceX_message.notification_service/internal/repository"
	"google.golang.org/grpc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
		cfg := config.LoadConfig()
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	 db := client.Database(cfg.DatabaseName)
    messagesCollection := db.Collection("messages")

	messageRepo := repository.NewMessageRepository(messagesCollection)
	messageService := service.NewMessageService(messageRepo)
	
	lis, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.ServerPort, err)
	}

	grpcServer := grpc.NewServer()
 proto.RegisterMessageServiceServer(grpcServer, messageService)

	fmt.Printf("Starting gRPC server on port %s...\n", cfg.ServerPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}