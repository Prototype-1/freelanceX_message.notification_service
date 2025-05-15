package email

import (
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
    EmailSender string
    EmailPass   string
    SMTPHost    string
    SMTPPort    string
}

func SendMail(cfg SMTPConfig, to string, subject, body string) error {
	auth := smtp.PlainAuth("", cfg.EmailSender, cfg.EmailPass, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	err := smtp.SendMail(addr, auth, cfg.EmailSender, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("send mail failed: %w", err)
	}
	return nil
}

