package kafka

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "github.com/segmentio/kafka-go"
)

type KafkaMessage struct {
    FromUserID string `json:"from_user_id"`
    ToUserID   string `json:"to_user_id"`
    ProjectID  string `json:"project_id"`
    Content    string `json:"content"`
    Timestamp  string `json:"timestamp"`
}

func ConsumeMessages(brokerAddr string, topic string, groupID string) {
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
        m, err := r.ReadMessage(context.Background())
        if err != nil {
            log.Printf("Error reading message: %v", err)
            continue
        }

        var event KafkaMessage
        if err := json.Unmarshal(m.Value, &event); err != nil {
            log.Printf("JSON unmarshal error: %v", err)
            continue
        }

        log.Printf("Received message: %+v", event)
        handleNotification(event)
    }
}

func handleNotification(msg KafkaMessage) {
    log.Printf("📣 Notifying %s about new message from %s on project %s: %s",
        msg.ToUserID, msg.FromUserID, msg.ProjectID, msg.Content)

    // TODO:
    // - Store in DB (if needed)
    // - Push real-time event to WebSocket layer (optional)
    // - Send email using SMTP or a service
}
