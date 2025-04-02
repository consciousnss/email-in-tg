package tg

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/un1uckyyy/email-in-tg/internal/imap"
	tele "gopkg.in/telebot.v4"
)

const (
	pollerTimeout = 10 * time.Second
)

type TelegramService struct {
	bot  *tele.Bot
	imap imap.ImapService
}

func NewTelegramService(
	token string,
	is imap.ImapService,
) (*TelegramService, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: pollerTimeout},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &TelegramService{
		bot:  bot,
		imap: is,
	}, nil
}

func (t *TelegramService) registerButtons() {
	t.bot.Handle("/start", func(c tele.Context) error {
		msg := fmt.Sprintf("Привет, %s!", c.Sender().Username)
		return c.Send(msg)
	})
	t.bot.Handle("/list", func(c tele.Context) error {
		boxes, err := t.imap.ListBoxes()
		if err != nil {
			log.Println("/list error:", err)
			return err
		}
		boxesMsg := strings.Join(boxes, "\n")
		return c.Send(boxesMsg)
	})
}

func (t *TelegramService) Start() {
	t.registerButtons()

	go func() {
		t.bot.Start()
	}()
}
