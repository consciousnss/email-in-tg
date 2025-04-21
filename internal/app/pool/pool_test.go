package pool

import (
	"context"
	"errors"
	"testing"

	"github.com/un1uckyyy/email-in-tg/internal/domain/mail"

	"github.com/un1uckyyy/email-in-tg/internal/domain/models"

	"github.com/stretchr/testify/assert"
)

type mockImapService struct {
	loginCalled bool
	startCalled bool
	stopCalled  bool

	loginErr error
	startErr error
	stopErr  error
}

func (m *mockImapService) Login(email, password string) error {
	m.loginCalled = true
	return m.loginErr
}

func (m *mockImapService) Start(ctx context.Context, updates chan<- *models.Update) error {
	m.startCalled = true
	return m.startErr
}

func (m *mockImapService) Stop(ctx context.Context) error {
	m.stopCalled = true
	return m.stopErr
}

type mockFactory struct {
	service *mockImapService
	err     error
}

func (f *mockFactory) New(sd models.MailServiceData) (mail.MailboxWatcher, error) {
	return f.service, f.err
}

func TestPool_Add_Success(t *testing.T) {
	ctx := context.Background()
	mockSvc := &mockImapService{}
	mockFact := &mockFactory{service: mockSvc}

	p := &pool{
		clients: make(map[int64]mail.MailboxWatcher),
		updates: make(chan *models.Update),
		factory: mockFact,
	}

	group := &models.Group{
		ID: 123,
		Login: &models.EmailLogin{
			Email:    "test@mail.ru",
			Password: "pass",
		},
	}

	err := p.Add(ctx, group)
	assert.NoError(t, err)
	assert.True(t, mockSvc.loginCalled)
	assert.True(t, mockSvc.startCalled)
}

func TestPool_Add_LoginError(t *testing.T) {
	ctx := context.Background()
	mockSvc := &mockImapService{loginErr: errors.New("login failed")}
	mockFact := &mockFactory{service: mockSvc}

	p := &pool{
		clients: make(map[int64]mail.MailboxWatcher),
		updates: make(chan *models.Update),
		factory: mockFact,
	}

	group := &models.Group{
		ID: 1,
		Login: &models.EmailLogin{
			Email:    "x@mail.ru",
			Password: "badpass",
		},
	}

	err := p.Add(ctx, group)
	assert.ErrorContains(t, err, "error imap login")
	assert.True(t, mockSvc.loginCalled)
}

func TestPool_Add_StartError(t *testing.T) {
	ctx := context.Background()
	mockSvc := &mockImapService{
		startErr: errors.New("start failed"),
	}
	mockFact := &mockFactory{service: mockSvc}

	p := &pool{
		clients: make(map[int64]mail.MailboxWatcher),
		updates: make(chan *models.Update),
		factory: mockFact,
	}

	group := &models.Group{
		ID: 2,
		Login: &models.EmailLogin{
			Email:    "test@mail.ru",
			Password: "pass",
		},
	}

	err := p.Add(ctx, group)
	assert.ErrorContains(t, err, "error imap start")
	assert.True(t, mockSvc.loginCalled)
	assert.True(t, mockSvc.startCalled)
}

func TestPool_Delete_Success(t *testing.T) {
	ctx := context.Background()
	mockSvc := &mockImapService{}
	p := &pool{
		clients: make(map[int64]mail.MailboxWatcher),
		updates: make(chan *models.Update),
	}

	group := &models.Group{ID: 3}
	p.clients[group.ID] = mockSvc

	err := p.Delete(ctx, group)
	assert.NoError(t, err)
	assert.True(t, mockSvc.stopCalled)
	_, exists := p.clients[group.ID]
	assert.False(t, exists)
}

func TestPool_Delete_StopError(t *testing.T) {
	ctx := context.Background()
	mockSvc := &mockImapService{stopErr: errors.New("stop failed")}
	p := &pool{
		clients: make(map[int64]mail.MailboxWatcher),
		updates: make(chan *models.Update),
	}

	group := &models.Group{ID: 4}
	p.clients[group.ID] = mockSvc

	err := p.Delete(ctx, group)
	assert.ErrorContains(t, err, "error stopping imap client")
	assert.True(t, mockSvc.stopCalled)
}

func TestPool_Add_NilLogin(t *testing.T) {
	ctx := context.Background()
	mockSvc := &mockImapService{}
	mockFact := &mockFactory{service: mockSvc}

	p := &pool{
		clients: make(map[int64]mail.MailboxWatcher),
		updates: make(chan *models.Update),
		factory: mockFact,
	}

	group := &models.Group{ID: 5, Login: nil}
	err := p.Add(ctx, group)
	assert.ErrorContains(t, err, "group login is nil")
	assert.False(t, mockSvc.loginCalled)
	assert.False(t, mockSvc.startCalled)
}

func TestPool_Updates_Channel(t *testing.T) {
	p := &pool{
		updates: make(chan *models.Update, 1),
	}

	expected := &models.Update{
		Email: &models.Email{
			Text: "hello",
		},
	}
	p.updates <- expected

	recv := <-p.Updates()
	assert.Equal(t, expected, recv)
}
