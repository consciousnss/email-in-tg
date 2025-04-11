package pool

import (
	"github.com/un1uckyyy/email-in-tg/internal/infra/imap"
	"github.com/un1uckyyy/email-in-tg/internal/models"
)

type ImapServiceFactory interface {
	New(address string, groupID int64, updates chan *models.Update) (imap.ImapService, error)
}

type defaultImapServiceFactory struct{}

func (f *defaultImapServiceFactory) New(
	address string,
	groupID int64,
	updates chan *models.Update,
) (imap.ImapService, error) {
	return imap.NewImapService(address, groupID, updates)
}
