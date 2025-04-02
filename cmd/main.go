package main

import (
	"fmt"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/un1uckyyy/email-in-tg/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := imapclient.DialTLS(cfg.ImapServer, nil)
	if err != nil {
		log.Fatal("Dial TLS error:", err)
	}

	err = client.Login(cfg.Username, cfg.Password).Wait()
	if err != nil {
		log.Fatal("Login error:", err)
	}

	fmt.Println("Logged in")
}
