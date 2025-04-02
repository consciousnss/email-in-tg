package imap

import (
	"fmt"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

type ImapService interface {
	ListBoxes() ([]string, error)
}

type imapService struct {
	c *imapclient.Client
}

func NewImapService(
	imapServer string,
	user string,
	pass string,
) (ImapService, error) {
	client, err := imapclient.DialTLS(imapServer, nil)
	if err != nil {
		return nil, fmt.Errorf("dial TLS error: %w", err)
	}

	err = client.Login(user, pass).Wait()
	if err != nil {
		return nil, fmt.Errorf("login error: %w", err)
	}

	return &imapService{
		c: client,
	}, nil
}

func (i *imapService) ListBoxes() ([]string, error) {
	options := imap.ListOptions{
		ReturnStatus: &imap.StatusOptions{
			NumMessages: true,
		},
	}
	mailboxes, err := i.c.List("", "%", &options).Collect()
	if err != nil {
		return nil, fmt.Errorf("list mailboxes error: %w", err)
	}

	boxes := make([]string, 0, len(mailboxes))
	for _, mailbox := range mailboxes {
		boxes = append(boxes, mailbox.Mailbox)
	}

	return boxes, nil
}

//selectData, err := client.Select("INBOX", nil).Wait()
//if err != nil {
//log.Fatalf("Select().Wait() = %v", err)
//}
//
//fmt.Println(selectData.NumMessages)
//
//criteria := imap.SearchCriteria{}
//sd, err := client.UIDSearch(&criteria, nil).Wait()
//if err != nil {
//log.Fatal("Search error: ", err)
//}
//fmt.Print("Searched: ", sd.AllUIDs())
