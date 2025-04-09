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
	Updates    chan *models.Email
	Register   chan *models.Group
	Unregister chan *models.Group
}

func NewPool() *Pool {
	return &Pool{
		Clients:    make(map[int64]imap.ImapService),
		Updates:    make(chan *models.Email),
		Register:   make(chan *models.Group),
		Unregister: make(chan *models.Group),
	}
}

func (p *Pool) Start(_ context.Context) error {
	logger.Debug("starting pool loop...")
	go p.run()
	return nil
}

func (p *Pool) run() {
	for {
		select {
		case email := <-p.Updates:
			msg := fmt.Sprintf("got email update: %+v", email)
			logger.Debug(msg)
		case group := <-p.Register:
			msg := fmt.Sprintf("starting group register: %+v", group)
			logger.Debug(msg)

			is, err := imap.NewImapService(mailRuImap)
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
			p.Clients[group.ID] = is
			msg = fmt.Sprintf("succesful register group: %+v", group)
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
