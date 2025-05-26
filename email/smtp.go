package email

import (
"bytes"
"encoding/base64"
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

func SendMailWithAttachment(cfg SMTPConfig, to, subject, body, filename string, attachment []byte) error {
	auth := smtp.PlainAuth("", cfg.EmailSender, cfg.EmailPass, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	boundary := "f46d043c813270fc6b04c2d223da"

	header := make(map[string]string)
	header["From"] = cfg.EmailSender
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "multipart/mixed; boundary=" + boundary
	var msg bytes.Buffer
	for k, v := range header {
		fmt.Fprintf(&msg, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(&msg, "\r\n--%s\r\n", boundary)

	fmt.Fprintf(&msg, "Content-Type: text/plain; charset=utf-8\r\n")
	fmt.Fprintf(&msg, "Content-Transfer-Encoding: 7bit\r\n\r\n")
	fmt.Fprintf(&msg, "%s\r\n", body)

	fmt.Fprintf(&msg, "\r\n--%s\r\n", boundary)
	fmt.Fprintf(&msg, "Content-Type: application/pdf\r\n")
	fmt.Fprintf(&msg, "Content-Transfer-Encoding: base64\r\n")
	fmt.Fprintf(&msg, "Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", filename)

	b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment)))
	base64.StdEncoding.Encode(b, attachment)

	for i, l := 0, len(b); i < l; i++ {
		fmt.Fprintf(&msg, "%c", b[i])
		if (i+1)%76 == 0 {
			fmt.Fprintf(&msg, "\r\n")
		}
	}
	fmt.Fprintf(&msg, "\r\n--%s--", boundary)

	return smtp.SendMail(addr, auth, cfg.EmailSender, []string{to}, msg.Bytes())
}
