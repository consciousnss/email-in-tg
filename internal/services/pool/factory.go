package pool

import (
	"github.com/un1uckyyy/email-in-tg/internal/domain/models"
	"github.com/un1uckyyy/email-in-tg/internal/infra/imap"
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
