package imap

import (
	"fmt"
	"time"

	"github.com/un1uckyyy/email-in-tg/internal/domain/models"

	"github.com/emersion/go-message/mail"
)

const (
	noSubject = "<Без темы>"
)

func parseHeader(header mail.Header, email *models.Email) error {
	err := getSubject(header, email)
	if err != nil {
		return fmt.Errorf("get subject error: %w", err)
	}

	err = getDate(header, email)
	if err != nil {
		return fmt.Errorf("get date error: %w", err)
	}

	err = getFrom(header, email)
	if err != nil {
		return fmt.Errorf("get from error: %w", err)
	}

	err = getTo(header, email)
	if err != nil {
		return fmt.Errorf("get to error: %w", err)
	}

	return nil
}

func getSubject(header mail.Header, email *models.Email) error {
	subject, err := header.Subject()
	if err != nil {
		return fmt.Errorf("get subject error: %w", err)
	}
	email.Subject = subject
	if subject == "" {
		email.Subject = noSubject
	}
	return nil
}

func getDate(header mail.Header, email *models.Email) error {
	date, err := header.Date()
	if err != nil {
		return fmt.Errorf("get date error: %w", err)
	}
	email.Date = date.Format(time.RFC1123)
	return nil
}

func getFrom(header mail.Header, email *models.Email) error {
	from, err := header.AddressList("From")
	if err != nil {
		return fmt.Errorf("get from address list error: %w", err)
	}

	if len(from) == 0 {
		return fmt.Errorf("no From address found")
	}

	if len(from) == 1 {
		email.From = from[0].Address
	}

	if len(from) > 1 {
		s, err := header.AddressList("Sender")
		if err != nil {
			return fmt.Errorf("get sender address list error: %w", err)
		}

		if len(s) == 0 {
			return fmt.Errorf("no Sender address found")
		}

		if len(s) > 1 {
			return fmt.Errorf("more than one Sender address found")
		}

		email.From = s[0].Address
	}

	return nil
}

func getTo(header mail.Header, email *models.Email) error {
	toAddrList, err := header.AddressList("To")
	if err != nil {
		return fmt.Errorf("get to address list error: %w", err)
	}

	if len(toAddrList) == 0 {
		return fmt.Errorf("no To address found")
	}

	to := make([]string, 0, len(toAddrList))
	for _, toAddr := range toAddrList {
		to = append(to, toAddr.Address)
	}
	email.To = to
	return nil
}
