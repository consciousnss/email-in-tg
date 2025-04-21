package models

import "time"

type Group struct {
	ID   int64 `validate:"required"`
	Type string

	Title string

	Provider     MailProvider
	Login        *EmailLogin
	PollInterval time.Duration

	IsActive bool
}

type MailProvider string

const (
	MailRuProvider MailProvider = "imap.mail.ru:993"
	YandexProvider MailProvider = "imap.yandex.ru:993"
)

type EmailLogin struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type Subscription struct {
	ID          string
	SenderEmail *string `validate:"omitempty,required_if=OtherSenders false,email"`
	GroupID     int64   `validate:"required"`
	ThreadID    int

	// if set to true, will be used if no SenderEmail matches found
	OtherSenders bool
}
