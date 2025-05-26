package kafka

import (
	"fmt"
	"log"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/Prototype-1/freelanceX_message.notification_service/email"
	"github.com/Prototype-1/freelanceX_message.notification_service/internal/client"
)

type InvoiceCreatedEvent struct {
    InvoiceID      string `json:"invoice_id"`
    ProjectTitle   string `json:"project_title"`
    FreelancerName string `json:"freelancer_name"`
    ClientID       string `json:"client_id"`
}

func ConsumeInvoiceEvents(broker, topic string, smtpCfg email.SMTPConfig, invoiceClient *client.InvoiceServiceClient, userClient *client.UserServiceClient) {
    log.Println("Started consuming invoice events from Kafka topic:", topic)
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{broker},
        Topic:   topic,
        GroupID: "invoice-notification-group",
    })

    for {
        msg, err := r.ReadMessage(context.Background())
        if err != nil {
            log.Printf("error reading message: %v", err)
            continue
        }

        var event InvoiceCreatedEvent
        if err := json.Unmarshal(msg.Value, &event); err != nil {
            log.Printf("error unmarshaling invoice event: %v", err)
            continue
        }

        log.Printf("Received invoice event: %+v", event)
        clientEmail, err := userClient.GetUserEmail(context.Background(), event.ClientID)
        if err != nil {
            log.Printf("Failed to get client email: %v", err)
            continue
        }
        pdfBytes, err := invoiceClient.GetInvoicePDF(event.InvoiceID)
        if err != nil {
            log.Printf("Failed to fetch invoice PDF: %v", err)
            continue
        }

        subject := fmt.Sprintf("New Invoice for Project: %s", event.ProjectTitle)
        body := fmt.Sprintf("Dear Client,\n\nA new invoice has been issued for the project \"%s\" by %s.\n\nPlease find the attached invoice.\n\nRegards,\nFreelanceX Team", event.ProjectTitle, event.FreelancerName)
        filename := fmt.Sprintf("invoice_%s.pdf", event.InvoiceID)

        if err := email.SendMailWithAttachment(smtpCfg, clientEmail, subject, body, filename, pdfBytes); err != nil {
            log.Printf("Failed to send invoice email: %v", err)
        } else {
            log.Printf("Invoice email sent to %s", clientEmail)
        }
    }
}
