package pool

import (
	"github.com/un1uckyyy/email-in-tg/internal/domain/mail"
	"github.com/un1uckyyy/email-in-tg/internal/infra/imap"
)

type MailboxWatcherFactory interface {
	New(address string, groupID int64) (mail.MailboxWatcher, error)
}

type defaultMailboxWatcherFactory struct{}

var _ MailboxWatcherFactory = (*defaultMailboxWatcherFactory)(nil)

func (f *defaultMailboxWatcherFactory) New(
	address string,
	id int64,
) (mail.MailboxWatcher, error) {
	return imap.NewImapService(address, id)
}
