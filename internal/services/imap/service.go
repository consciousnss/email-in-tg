package imap

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/mail"
	"github.com/un1uckyyy/email-in-tg/internal/models"
)

type ImapService interface {
	Login(username string, password string) error
	Logout() error
	ListBoxes() ([]string, error)
	Select(mailbox string) error
	FetchOne(num uint32, uid bool) (*models.Email, error)
}

type ImapServiceImpl struct {
	c *imapclient.Client
}

var _ ImapService = (*ImapServiceImpl)(nil)

func NewImapService(
	imapServer string,
) (ImapService, error) {
	client, err := imapclient.DialTLS(imapServer,
		&imapclient.Options{
			// DebugWriter: os.Stderr,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("dial TLS error: %w", err)
	}

	return &ImapServiceImpl{
		c: client,
	}, nil
}

func (i *ImapServiceImpl) Start(_ context.Context) error {
	return nil
}

func (i *ImapServiceImpl) Login(username string, password string) error {
	err := i.c.Login(username, password).Wait()
	if err != nil {
		return fmt.Errorf("login error: %w", err)
	}
	return nil
}

func (i *ImapServiceImpl) Logout() error {
	err := i.c.Logout().Wait()
	if err != nil {
		return fmt.Errorf("logout error: %w", err)
	}
	return nil
}

func (i *ImapServiceImpl) ListBoxes() ([]string, error) {
	options := imap.ListOptions{
		ReturnStatus: &imap.StatusOptions{
			NumMessages: true,
			NumUnseen:   true,
		},
	}

	mailboxes, err := i.c.List("", "%", &options).Collect()
	if err != nil {
		return nil, fmt.Errorf("list mailboxes error: %w", err)
	}

	boxes := make([]string, 0, len(mailboxes))
	for _, mailbox := range mailboxes {
		msg := fmt.Sprintf(
			"Mailbox %s contains %v messages (%v unseen)",
			mailbox.Mailbox,
			*mailbox.Status.NumMessages,
			*mailbox.Status.NumUnseen,
		)
		boxes = append(boxes, msg)
	}

	return boxes, nil
}

func (i *ImapServiceImpl) Select(mailbox string) error {
	_, err := i.c.Select(mailbox, nil).Wait()
	if err != nil {
		return fmt.Errorf("select error: %w", err)
	}

	return nil
}

func (i *ImapServiceImpl) FetchOne(num uint32, uid bool) (*models.Email, error) {
	email := &models.Email{}

	var seqSet imap.NumSet
	if uid {
		seqSet = imap.UIDSetNum(imap.UID(num))
	} else {
		seqSet = imap.SeqSetNum(num)
	}

	bodySection := &imap.FetchItemBodySection{}
	fetchOptions := &imap.FetchOptions{
		BodySection: []*imap.FetchItemBodySection{
			bodySection,
		},
	}
	fetchCmd := i.c.Fetch(seqSet, fetchOptions)
	defer func() {
		err := fetchCmd.Close()
		if err != nil {
			msg := fmt.Sprintf("fetch close err: %v", err)
			logger.Error(msg)
		}
	}()

	for {
		msg := fetchCmd.Next()
		if msg == nil {
			break
		}

		for {
			item := msg.Next()
			if item == nil {
				break
			}

			if item, ok := item.(imapclient.FetchItemDataBodySection); ok {
				mr, err := mail.CreateReader(item.Literal)
				if err != nil {
					return nil, fmt.Errorf("mail parse err: %w", err)
				}

				s, err := mr.Header.Subject()
				if err != nil {
					return nil, fmt.Errorf("get subject error: %w", err)
				}
				email.Subject = s
				logger.Debug(fmt.Sprintf("got subject: %v", s))

				d, err := mr.Header.Date()
				if err != nil {
					return nil, fmt.Errorf("get date error: %w", err)
				}
				email.Date = d
				logger.Debug(fmt.Sprintf("got date: %v", d))

				alf, err := mr.Header.AddressList("From")
				if err != nil {
					return nil, fmt.Errorf("get address list error: %w", err)
				}
				email.MailFrom = alf[0].Address
				logger.Debug(fmt.Sprintf("got address: %v", alf[0].Address))

				for {
					p, err := mr.NextPart()
					if err == io.EOF {
						break
					} else if err != nil {
						return nil, fmt.Errorf("")
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
		}
	}

	return email, nil
}
