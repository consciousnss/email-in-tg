package pool

import (
	"context"
	"fmt"

	"github.com/un1uckyyy/email-in-tg/internal/models"
	"github.com/un1uckyyy/email-in-tg/internal/services/imap"
)

const (
	mailRuImap = "imap.mail.ru:993"
)

type Pool struct {
	Clients    map[int64]imap.ImapService
	Updates    chan *models.Update
	Register   chan *models.Group
	Unregister chan *models.Group
}

func NewPool() *Pool {
	return &Pool{
		Clients:    make(map[int64]imap.ImapService),
		Updates:    make(chan *models.Update),
		Register:   make(chan *models.Group),
		Unregister: make(chan *models.Group),
	}
}

func (p *Pool) Start(ctx context.Context) error {
	logger.Debug("starting pool loop...")
	go p.run(ctx)
	return nil
}

func (p *Pool) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case group := <-p.Register:
			msg := fmt.Sprintf("starting group register: %v", group.ID)
			logger.Debug(msg)

			is, err := imap.NewImapService(mailRuImap, group.ID, p.Updates)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			if group.Login == nil {
				logger.Error("login is nil")
				continue
			}

			err = is.Login(group.Login.Email, group.Login.Password)
			if err != nil {
				logger.Error(err.Error())
				continue
			}

			err = is.Start(ctx)
			if err != nil {
				msg := fmt.Sprintf("failed to start imap service: %v", err)
				logger.Error(msg)
				continue
			}

			p.Clients[group.ID] = is
			msg = fmt.Sprintf("succesful register group: %v", group.ID)
			logger.Debug(msg)
		case group := <-p.Unregister:
			err := p.Clients[group.ID].Logout()
			if err != nil {
				logger.Error(err.Error())
			}
			delete(p.Clients, group.ID)

			msg := fmt.Sprintf("got group unregister update: %+v", group)
			logger.Debug(msg)
		}
	}
}
