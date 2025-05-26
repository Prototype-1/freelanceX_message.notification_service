package kafka

import (
    "context"
    "encoding/json"
    "log"
    "github.com/Prototype-1/freelanceX_message.notification_service/email"
    "github.com/Prototype-1/freelanceX_message.notification_service/internal/client"
    "github.com/segmentio/kafka-go"
)

type ProposalEvent struct {
    ProposalID   string `json:"proposal_id"`
    ClientID     string `json:"client_id"`
    FreelancerID string `json:"freelancer_id"`
    Title        string `json:"title"`
    EventType    string `json:"event_type"` 
    Status       string `json:"status"`     
}

func ConsumeProposalEvents(broker, topic string, smtpCfg email.SMTPConfig, userClient *client.UserServiceClient) {
    log.Println("Started consuming proposal events from Kafka topic:", topic)
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:  []string{broker},
        Topic:    topic,
        GroupID:  "proposal-notification-group",
    })

    for {
        msg, err := r.ReadMessage(context.Background())
        if err != nil {
            log.Printf("error reading message: %v", err)
            continue
        }

        var event ProposalEvent
        if err := json.Unmarshal(msg.Value, &event); err != nil {
            log.Printf("error unmarshaling event: %v", err)
            continue
        }

        log.Printf("Received event: %+v", event)

  switch event.EventType {

        case "proposal.created":
            clientEmail, err := userClient.GetUserEmail(context.Background(), event.ClientID)
            if err != nil {
                log.Printf("Failed to get client email: %v", err)
                continue
            }

            subject := "New Proposal Created"
            body := "A freelancer has submitted a new proposal for your project: \"" + event.Title + "\"."

            if err := email.SendMail(smtpCfg, clientEmail, subject, body); err != nil {
                log.Printf("Failed to send email to client: %v", err)
            } else {
                log.Printf("Email sent to client %s", clientEmail)
            }

        case "proposal.status.updated":
            if event.Status == "accepted" || event.Status == "rejected" {
                freelancerEmail, err := userClient.GetUserEmail(context.Background(), event.FreelancerID)
                if err != nil {
                    log.Printf("Failed to get freelancer email: %v", err)
                    continue
                }

                subject := "Proposal Status Update"
                body := "Your proposal titled \"" + event.Title + "\" was *" + event.Status + "* by the client."

                if err := email.SendMail(smtpCfg, freelancerEmail, subject, body); err != nil {
                    log.Printf(" Failed to send email to freelancer: %v", err)
                } else {
                    log.Printf("Email sent to freelancer %s", freelancerEmail)
                }
            } else {
                log.Printf("Ignored status update with unsupported status: %s", event.Status)
            }

        default:
            log.Printf("Unknown event type: %s", event.EventType)
        }
    }
}
