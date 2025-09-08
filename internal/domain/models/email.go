package models

import (
	"io"
	"time"
)

type MailServiceData struct {
	GroupID      int64
	Provider     MailProvider
	PollInterval time.Duration
}

type Update struct {
	GroupID int64
	Email   *Email
}

type Email struct {
	From    string
	To      []string
	Date    string
	Subject string
	Text    string
	Files   []*File
}

type File struct {
	Filename string
	Data     io.Reader
}
