package mail

import (
	"context"

	"github.com/consciousnss/email-in-tg/internal/domain/models"
)

type MailboxWatcher interface {
	Login(username string, password string) error
	Start(ctx context.Context, updates chan<- *models.Update) error
	Stop(ctx context.Context) error
}

type MailboxWatcherType string

const (
	ImapMailboxWatcher MailboxWatcherType = "imap"
	SMTPMailboxWatcher MailboxWatcherType = "smtp"
)
