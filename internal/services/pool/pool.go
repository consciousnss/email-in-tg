package pool

import (
	"context"
	"errors"
	"fmt"

	"github.com/un1uckyyy/email-in-tg/internal/infra/imap"

	"github.com/un1uckyyy/email-in-tg/internal/models"
)

const (
	mailRuImap = "imap.mail.ru:993"
)

type Pool interface {
	Updates() <-chan *models.Update
	Add(ctx context.Context, group *models.Group) error
	Delete(ctx context.Context, group *models.Group) error
}

type pool struct {
	clients map[int64]imap.ImapService
	updates chan *models.Update
	factory ImapServiceFactory
}

var _ Pool = (*pool)(nil)

func NewPool() Pool {
	return &pool{
		clients: make(map[int64]imap.ImapService),
		updates: make(chan *models.Update),
		factory: &defaultImapServiceFactory{},
	}
}

func (p *pool) Updates() <-chan *models.Update {
	return p.updates
}

func (p *pool) Add(ctx context.Context, group *models.Group) error {
	msg := fmt.Sprintf("starting group register: %v", group.ID)
	logger.Debug(msg)

	is, err := p.factory.New(mailRuImap, group.ID, p.updates)
	if err != nil {
		return fmt.Errorf("error creating imap service: %v", err)
	}
	if group.Login == nil {
		return errors.New("group login is nil")
	}

	err = is.Login(group.Login.Email, group.Login.Password)
	if err != nil {
		return fmt.Errorf("error imap login: %v", err)
	}

	err = is.Start(ctx)
	if err != nil {
		return fmt.Errorf("error imap start: %v", err)
	}

	p.clients[group.ID] = is
	msg = fmt.Sprintf("succesful register group: %v", group.ID)
	logger.Debug(msg)

	return nil
}

func (p *pool) Delete(ctx context.Context, group *models.Group) error {
	err := p.clients[group.ID].Stop(ctx)
	if err != nil {
		return fmt.Errorf("error stopping imap client: %v", err)
	}
	delete(p.clients, group.ID)

	msg := fmt.Sprintf("got group unregister update: %+v", group)
	logger.Debug(msg)
	return nil
}
