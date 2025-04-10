package models

import (
	"io"
	"time"
)

type Email struct {
	MailFrom string
	MailTo   string
	Date     time.Time
	Subject  string
	Text     string
	Files    []*File
}

type File struct {
	Filename string
	Data     io.Reader
}
