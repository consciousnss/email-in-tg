package main

import (
	"log"

	"github.com/un1uckyyy/email-in-tg/internal/tg"

	"github.com/un1uckyyy/email-in-tg/internal/config"
	"github.com/un1uckyyy/email-in-tg/internal/imap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	is, err := imap.NewImapService(cfg.ImapServer, cfg.Username, cfg.Password)
	if err != nil {
		log.Fatal(err)
	}

	ts, err := tg.NewTelegramService(cfg.TelegramToken, is)
	if err != nil {
		log.Fatal(err)
	}
	ts.Start()

	log.Println("app started...")
	select {}
}
