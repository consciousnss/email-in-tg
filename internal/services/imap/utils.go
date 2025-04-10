package imap

import (
	"fmt"
	"time"

	"github.com/emersion/go-message/mail"
	"github.com/un1uckyyy/email-in-tg/internal/models"
)

func headerParse(header mail.Header, email *models.Email) error {
	s, err := header.Subject()
	if err != nil {
		return fmt.Errorf("get subject error: %w", err)
	}
	email.Subject = s

	d, err := header.Date()
	if err != nil {
		return fmt.Errorf("get date error: %w", err)
	}
	email.Date = d.Format(time.RFC1123)

	alf, err := header.AddressList("From")
	if err != nil {
		return fmt.Errorf("get address list error: %w", err)
	}
	email.MailFrom = alf[0].Address

	alt, err := header.AddressList("To")
	if err != nil {
		return fmt.Errorf("get address list error: %w", err)
	}
	email.MailTo = alt[0].Address

	return nil
}
