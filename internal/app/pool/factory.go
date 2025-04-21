package pool

import (
	"github.com/un1uckyyy/email-in-tg/internal/domain/mail"
	"github.com/un1uckyyy/email-in-tg/internal/domain/models"
	"github.com/un1uckyyy/email-in-tg/internal/infra/imap"
)

type MailboxWatcherFactory interface {
	New(sd models.MailServiceData) (mail.MailboxWatcher, error)
}

type defaultMailboxWatcherFactory struct{}

var _ MailboxWatcherFactory = (*defaultMailboxWatcherFactory)(nil)

func (f *defaultMailboxWatcherFactory) New(
	sd models.MailServiceData,
) (mail.MailboxWatcher, error) {
	return imap.NewImapService(sd)
}
