package imap

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/un1uckyyy/email-in-tg/internal/domain/models"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/mail"
	mailboxwatcher "github.com/un1uckyyy/email-in-tg/internal/domain/mail"
)

type imapService struct {
	serviceData models.MailServiceData
	c           *imapclient.Client
	updates     chan<- *models.Update
	done        chan struct{}
}

var _ mailboxwatcher.MailboxWatcher = (*imapService)(nil)

const (
	inbox = "INBOX"
)

func NewImapService(
	sd models.MailServiceData,
) (mailboxwatcher.MailboxWatcher, error) {
	is := &imapService{
		serviceData: sd,
		c:           nil,
		done:        make(chan struct{}),
	}

	client, err := is.connect()
	if err != nil {
		return nil, fmt.Errorf("dial TLS error: %w", err)
	}
	is.c = client
	return is, nil
}

func (i *imapService) connect() (*imapclient.Client, error) {
	provider := string(i.serviceData.Provider)
	client, err := imapclient.DialTLS(provider, nil)
	if err != nil {
		return nil, fmt.Errorf("dial TLS error: %w", err)
	}
	return client, nil
}

func (i *imapService) Start(ctx context.Context, updates chan<- *models.Update) error {
	i.updates = updates
	go i.run(ctx)
	return nil
}

func (i *imapService) Stop(_ context.Context) error {
	err := i.logout()
	if err != nil {
		return err
	}
	i.done <- struct{}{}
	return nil
}

func (i *imapService) run(ctx context.Context) {
	ticker := time.NewTicker(i.serviceData.PollInterval)
	defer ticker.Stop()

	uidNext, err := i.Status()
	if err != nil {
		msg := fmt.Sprintf("imap status error: %s", err)
		logger.Error(msg)
	}

	msg := fmt.Sprintf("got UIDNext: %d", uidNext)
	logger.Debug(msg)

	for {
		select {
		case <-ctx.Done():
			return
		case <-i.done:
			return
		case <-ticker.C:
			uidNextNext, err := i.Status()
			if err != nil {
				msg := fmt.Sprintf("imap status error: %s", err)
				logger.Error(msg)
				break
			}

			if uidNextNext == uidNext {
				break
			}

			// TODO add fetching all mails from changed delta.
			msg := fmt.Sprintf("UIDNext changed from: %d, to: %d", uidNext, uidNextNext)
			logger.Debug(msg)

			err = i.selectMailbox(inbox)
			if err != nil {
				msg := fmt.Sprintf("imap select error: %s", err)
				logger.Error(msg)
				break
			}

			email, err := i.fetchOne(uidNext)
			if err != nil {
				msg := fmt.Sprintf("fetch uidNextNext %d error: %s", uidNext, err)
				logger.Error(msg)
				break
			}
			i.updates <- &models.Update{
				Email:   email,
				GroupID: i.serviceData.GroupID,
			}
			uidNext = uidNextNext
		}
	}
}

func (i *imapService) Login(username string, password string) error {
	err := i.c.Login(username, password).Wait()
	if err != nil {
		return fmt.Errorf("login error: %w", err)
	}
	return nil
}

func (i *imapService) logout() error {
	err := i.c.Logout().Wait()
	if err != nil {
		return fmt.Errorf("logout error: %w", err)
	}
	return nil
}

func (i *imapService) selectMailbox(mailbox string) error {
	_, err := i.c.Select(mailbox, nil).Wait()
	if err != nil {
		return fmt.Errorf("select error: %w", err)
	}

	return nil
}

func (i *imapService) Status() (imap.UID, error) {
	data, err := i.c.Status(inbox, &imap.StatusOptions{UIDNext: true}).Wait()
	if err != nil {
		return 0, fmt.Errorf("status error: %w", err)
	}
	return data.UIDNext, nil
}

func (i *imapService) fetchOne(uid imap.UID) (*models.Email, error) {
	email := &models.Email{}

	seqSet := imap.UIDSetNum(uid)

	bodySection := &imap.FetchItemBodySection{}
	fetchOptions := &imap.FetchOptions{
		BodySection: []*imap.FetchItemBodySection{
			bodySection,
		},
	}
	fetchCmd := i.c.Fetch(seqSet, fetchOptions)
	defer fetchCmd.Close()

	msg := fetchCmd.Next()
	if msg == nil {
		return nil, fmt.Errorf("got nil fetch result")
	}

	for {
		item := msg.Next()
		if item == nil {
			break
		}

		dataBodySection, ok := item.(imapclient.FetchItemDataBodySection)
		if !ok {
			continue
		}

		mr, err := mail.CreateReader(dataBodySection.Literal)
		if err != nil {
			return nil, fmt.Errorf("mail parse err: %w", err)
		}

		err = parseHeader(mr.Header, email)
		if err != nil {
			return nil, fmt.Errorf("header parse err: %w", err)
		}
		logger.Debug(fmt.Sprintf("got %+v", email))

		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, fmt.Errorf("mail reader error: %w", err)
			}

			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				b, err := io.ReadAll(p.Body)
				if err != nil {
					return nil, fmt.Errorf("read text error %w", err)
				}
				email.Text = string(b)
			case *mail.AttachmentHeader:
				filename, err := h.Filename()
				if err != nil {
					return nil, fmt.Errorf("get filename error %w", err)
				}

				b, err := io.ReadAll(p.Body)
				if err != nil {
					return nil, fmt.Errorf("read attachment error: %w", err)
				}

				email.Files = append(email.Files, &models.File{
					Filename: filename,
					Data:     bytes.NewReader(b),
				})
			}
		}
	}

	return email, nil
}
