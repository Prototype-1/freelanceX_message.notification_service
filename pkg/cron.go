package pkg

import (
    "context"
    "log"
    "time"
    "github.com/Prototype-1/freelanceX_message.notification_service/internal/model"
    "github.com/Prototype-1/freelanceX_message.notification_service/internal/repository"
)

func StartDeliveryStatusCron(ctx context.Context, repo *repository.MessageRepository) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cutoff := time.Now().Add(-5 * time.Minute)
            messages, err := repo.FindPendingDeliveryMessages(ctx, cutoff)
            if err != nil {
                log.Printf("Error fetching pending delivery messages: %v", err)
                continue
            }

            for _, msg := range messages {
                err := repo.UpdateMessageStatus(ctx, msg.ID, model.DeliveredStatus)
                if err != nil {
                    log.Printf("Failed to update message status for ID %s: %v", msg.ID.Hex(), err)
                } else {
                    log.Printf("Message ID %s status updated to delivered", msg.ID.Hex())
                }
            }
        case <-ctx.Done():
            log.Println("Stopping delivery status cron")
            return
        }
    }
}
