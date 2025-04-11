package imap

import (
	"fmt"
	"time"

	"github.com/emersion/go-message/mail"
	"github.com/un1uckyyy/email-in-tg/internal/models"
)

func parseHeader(header mail.Header, email *models.Email) error {
	subject, err := header.Subject()
	if err != nil {
		return fmt.Errorf("get subject error: %w", err)
	}
	email.Subject = subject

	date, err := header.Date()
	if err != nil {
		return fmt.Errorf("get date error: %w", err)
	}
	email.Date = date.Format(time.RFC1123)

	from, err := header.AddressList("From")
	if err != nil {
		return fmt.Errorf("get from address list error: %w", err)
	}
	if len(from) == 0 {
		return fmt.Errorf("no From address found")
	}
	email.MailFrom = from[0].Address

	to, err := header.AddressList("To")
	if err != nil {
		return fmt.Errorf("get to address list error: %w", err)
	}
	if len(to) == 0 {
		return fmt.Errorf("no To address found")
	}
	email.MailTo = to[0].Address

	return nil
}
