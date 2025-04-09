package services

import (
	"context"

	"github.com/un1uckyyy/email-in-tg/internal/services/imap"
	"github.com/un1uckyyy/email-in-tg/internal/services/pool"
	"github.com/un1uckyyy/email-in-tg/internal/services/tg"
)

type Service interface {
	Start(ctx context.Context) error
}

var (
	_ Service = (*tg.TelegramService)(nil)
	_ Service = (*pool.Pool)(nil)
	_ Service = (*imap.ImapServiceImpl)(nil)
)
