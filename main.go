package main

import (
"log"
"context"
"net"
"fmt"
"github.com/Prototype-1/freelanceX_message.notification_service/config"
"github.com/Prototype-1/freelanceX_message.notification_service/email"
proto "github.com/Prototype-1/freelanceX_message.notification_service/proto"
"github.com/Prototype-1/freelanceX_message.notification_service/internal/service"
clt "github.com/Prototype-1/freelanceX_message.notification_service/internal/client"
authPb "github.com/Prototype-1/freelanceX_message.notification_service/proto/user_service"
"github.com/Prototype-1/freelanceX_message.notification_service/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
"github.com/Prototype-1/freelanceX_message.notification_service/kafka"

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
	emailCfg := email.SMTPConfig{
	EmailSender: cfg.SMTP.EmailSender,
	EmailPass:   cfg.SMTP.EmailPass,
	SMTPHost:    cfg.SMTP.SMTPHost,
	SMTPPort:    cfg.SMTP.SMTPPort,
}

 userConn, err := grpc.NewClient(
        cfg.UserServiceAddress,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatalf("Failed to connect to User Service: %v", err)
    }
	defer userConn.Close()

	authClient := authPb.NewAuthServiceClient(userConn)
	userClient := clt.NewUserServiceClient(authClient)
	messageService := service.NewMessageService(messageRepo, emailCfg, userClient)
	
	lis, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.ServerPort, err)
	}

	grpcServer := grpc.NewServer()
 proto.RegisterMessageServiceServer(grpcServer, messageService)

 	go func() {
		broker := "localhost:9092" 
		topic := "new.message"     
		groupID := "notification-group" 

		 log.Printf("Starting Kafka consumer with broker: %s, topic: %s, groupID: %s", 
               broker, topic, groupID)
		consumer := kafka.NewConsumer(messageRepo)
        consumer.ConsumeMessages(broker, topic, groupID)
	}()

	fmt.Printf("Starting gRPC server on port %s...\n", cfg.ServerPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}