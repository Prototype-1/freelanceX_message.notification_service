package kafka

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"
    "github.com/segmentio/kafka-go"
       "github.com/Prototype-1/freelanceX_message.notification_service/internal/model"
    "github.com/Prototype-1/freelanceX_message.notification_service/internal/repository"
    "github.com/google/uuid"
)

type KafkaMessage struct {
    FromUserID string `json:"from_user_id"`
    ToUserID   string `json:"to_user_id"`
    ProjectID  string `json:"project_id"`
    Content    string `json:"content"`
    Timestamp  string `json:"timestamp"`
}

type Consumer struct {
    Repo *repository.MessageRepository
}

func NewConsumer(repo *repository.MessageRepository) *Consumer {
    return &Consumer{Repo: repo}
}

func (c *Consumer) ConsumeMessages(brokerAddr string, topic string, groupID string) {
    log.Printf("Attempting to connect to Kafka at: %s", brokerAddr)
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:  []string{brokerAddr},
        GroupID:  groupID,
        Topic:    topic,
        MinBytes: 1,
        MaxBytes: 10e6,
    })

    defer r.Close()

    fmt.Println("Kafka consumer started...")
    for {
        log.Printf("Waiting for message from topic: %s", topic)
        m, err := r.ReadMessage(context.Background())
        if err != nil {
            log.Printf("Error reading message: %v", err)
            time.Sleep(5 * time.Second)
            continue
        }

        var event KafkaMessage
        if err := json.Unmarshal(m.Value, &event); err != nil {
            log.Printf("JSON unmarshal error: %v", err)
            continue
        }

        log.Printf("Received message: %+v", event)
        c.handleNotification(event)
    }
}

func (c *Consumer) handleNotification(msg KafkaMessage) {
    log.Printf("Notifying %s about new message from %s on project %s: %s",
        msg.ToUserID, msg.FromUserID, msg.ProjectID, msg.Content)

    fromUUID, _ := uuid.Parse(msg.FromUserID)
    toUUID, _ := uuid.Parse(msg.ToUserID)
    projectUUID, _ := uuid.Parse(msg.ProjectID)

    message := model.Message{
        FromUserID:  fromUUID,
        ToUserID:    toUUID,
        ProjectID:   projectUUID,
        MessageText: msg.Content,
        Channel:     model.InAppChannel,
        SentAt:      time.Now(),
        Read:        false,
        Status:      model.DeliveredStatus,
        IsDeleted:   false,
    }

    _, err := c.Repo.SaveMessage(context.Background(), &message)
    if err != nil {
        log.Printf("Failed to store message: %v", err)
        return
    }
    go func() {
        log.Printf("[WebSocket] Notify user %s about new message on project %s", msg.ToUserID, msg.ProjectID)
    }()

    go func() {
        log.Printf("[Email] Sending email to %s: You have a new message", msg.ToUserID)
    }()
}
