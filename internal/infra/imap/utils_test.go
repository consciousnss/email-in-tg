package imap

import (
	"testing"
	"time"

	"github.com/un1uckyyy/email-in-tg/internal/domain/models"

	"github.com/emersion/go-message/mail"
	"github.com/stretchr/testify/assert"
)

const (
	subject = "Test subject"
	from    = "alice@example.com"
)

var to = []string{"bob@example.com"}

func buildTestHeader(subject, from string, to []string, date time.Time) mail.Header {
	fromList := []*mail.Address{{Name: "", Address: from}}

	toList := make([]*mail.Address, 0, len(to))
	for _, addr := range to {
		toList = append(toList, &mail.Address{Name: "", Address: addr})
	}

	var h mail.Header
	h.SetSubject(subject)
	h.SetDate(date)
	h.SetAddressList("From", fromList)
	h.SetAddressList("To", toList)

	return h
}

func TestParseHeader_Success(t *testing.T) {
	date := time.Date(2025, 4, 10, 15, 4, 5, 0, time.UTC)

	header := buildTestHeader(subject, from, to, date)

	email := &models.Email{}

	err := parseHeader(header, email)

	assert.NoError(t, err)
	assert.Equal(t, subject, email.Subject)
	assert.Equal(t, from, email.From)
	assert.Equal(t, to, email.To)
}

func TestParseHeader_MissingSubject(t *testing.T) {
	date := time.Date(2025, 4, 10, 15, 4, 5, 0, time.UTC)

	header := buildTestHeader("", from, to, date)

	email := &models.Email{}

	err := parseHeader(header, email)

	assert.NoError(t, err)
	assert.Equal(t, email.Subject, noSubject)
}

func TestParseHeader_MissingFrom(t *testing.T) {
	date := time.Date(2025, 4, 10, 15, 4, 5, 0, time.UTC)

	header := buildTestHeader(subject, "", to, date)

	email := &models.Email{}

	err := parseHeader(header, email)

	assert.Error(t, err)
}

func TestParseHeader_MissingTo(t *testing.T) {
	date := time.Date(2025, 4, 10, 15, 4, 5, 0, time.UTC)

	header := buildTestHeader(subject, from, []string{}, date)

	email := &models.Email{}

	err := parseHeader(header, email)

	assert.Error(t, err)
}
