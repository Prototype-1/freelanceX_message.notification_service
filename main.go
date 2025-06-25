package main

import (
"os"
"log"
"context"
"net"
 "os/signal"
"syscall"
"github.com/Prototype-1/freelanceX_message.notification_service/config"
"github.com/Prototype-1/freelanceX_message.notification_service/pkg"
"github.com/Prototype-1/freelanceX_message.notification_service/email"
proto "github.com/Prototype-1/freelanceX_message.notification_service/proto"
"github.com/Prototype-1/freelanceX_message.notification_service/internal/service"
clt "github.com/Prototype-1/freelanceX_message.notification_service/internal/client"
authPb "github.com/Prototype-1/freelanceX_message.notification_service/proto/user_service"
pb "github.com/Prototype-1/freelanceX_message.notification_service/proto/invoice_service"
"github.com/Prototype-1/freelanceX_message.notification_service/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
"github.com/Prototype-1/freelanceX_message.notification_service/kafka"

)

func main() {
	cfg := config.LoadConfig()
	pkg.InitRedis(cfg.RedisAddr)

	ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

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

 go pkg.StartDeliveryStatusCron(ctx, messageRepo)

 	go func() {
broker := os.Getenv("KAFKA_BROKER")
if broker == "" {
    broker = "kafka:9092"
}
		topic := "new.message"     
		groupID := "notification-group" 

		 log.Printf("Starting Kafka consumer with broker: %s, topic: %s, groupID: %s", 
               broker, topic, groupID)
		consumer := kafka.NewConsumer(messageRepo)
        consumer.ConsumeMessages(broker, topic, groupID)
	}()

	go func() {
    		broker := os.Getenv("KAFKA_BROKER")
if broker == "" {
    broker = "kafka:9092"
}
    topic := "proposal-events"
    log.Printf("Starting Proposal Kafka consumer with broker: %s, topic: %s", broker, topic)
    kafka.ConsumeProposalEvents(broker, topic, emailCfg, userClient)
}()

go func() {
			broker := os.Getenv("KAFKA_BROKER")
if broker == "" {
    broker = "kafka:9092"
}
	topic := "invoice-events"
	log.Printf("Starting Invoice Kafka consumer with broker: %s, topic: %s", broker, topic)

	invoiceConn, err := grpc.NewClient(
		cfg.InvoiceServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to connect to Invoice Service: %v", err)
		return
	}
rawInvoiceClient := pb.NewInvoiceServiceClient(invoiceConn)
invoiceClient := clt.NewInvoiceServiceClient(rawInvoiceClient)
	kafka.ConsumeInvoiceEvents(broker, topic, emailCfg, invoiceClient, userClient)
}()

	   go func() {
        stopChan := make(chan os.Signal, 1)
        signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
        <-stopChan
        log.Println("Shutdown signal received, stopping services...")
        cancel()            
        grpcServer.GracefulStop()
    }()

	log.Printf("Starting gRPC server on port %s...\n", cfg.ServerPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}