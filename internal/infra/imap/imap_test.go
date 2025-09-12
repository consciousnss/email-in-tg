package imap

import (
	"image"
	"image/jpeg"
	"os"
	"strings"
	"testing"

	"github.com/consciousnss/email-in-tg/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func decodeFile(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return jpeg.Decode(f)
}

func TestParseOne_WithRealEmlAttachment(t *testing.T) {
	f, err := os.Open("./testdata/email.eml")
	assert.NoError(t, err)
	defer f.Close()
	require.NoError(t, err)

	email := &models.Email{}
	err = parseOne(f, email)
	require.NoError(t, err)

	assert.Equal(t, "2 фото и текст", email.Subject)
	assert.Equal(t, "mohamedlowskill@mail.ru", email.From)
	assert.Equal(t, []string{"redirect_test@mail.ru"}, email.To)
	assert.Equal(t, "Wed, 10 Sep 2025 23:26:12 +0300", email.Date)

	text, err := os.ReadFile("./testdata/email_text.html")
	require.NoError(t, err)
	assert.Equal(t, string(text), email.Text)

	require.Len(t, email.Files, 2)
	firstFilename := email.Files[0].Filename
	assert.Equal(t, "hockey.jpeg", firstFilename)
	secondFilename := email.Files[1].Filename
	assert.Equal(t, "man.jpeg", secondFilename)

	hockeyImgExpected, err := decodeFile("./testdata/hockey.jpeg")
	assert.NoError(t, err)
	hockeyImgActual, err := jpeg.Decode(email.Files[0].Data)
	assert.NoError(t, err)

	assert.Equal(t, hockeyImgExpected.Bounds(), hockeyImgActual.Bounds())
}

func TestParseOne_PlainTextNoFiles(t *testing.T) {
	raw := `From: a@a.com
To: b@b.com
Subject: hi
Date: Mon, 07 Apr 2025 23:43:15 +0300
MIME-Version: 1.0
Content-Type: text/plain; charset=utf-8

Hello world
`
	email := &models.Email{}
	err := parseOne(strings.NewReader(raw), email)
	assert.NoError(t, err)

	assert.Contains(t, email.Text, "Hello world")
	assert.Equal(t, "a@a.com", email.From)
	assert.Equal(t, []string{"b@b.com"}, email.To)
	assert.Equal(t, "hi", email.Subject)
	assert.Equal(t, "Mon, 07 Apr 2025 23:43:15 +0300", email.Date)
	assert.Len(t, email.Files, 0)
}

func TestParseOne_PlainTextNoFilesWithSender(t *testing.T) {
	raw := `From: a@a.com, c@b.com
Sender: sender@send.com
To: b@b.com
Subject: hi
Date: Mon, 07 Apr 2025 23:43:15 +0300
MIME-Version: 1.0
Content-Type: text/plain; charset=utf-8

Hello world
`
	email := &models.Email{}
	err := parseOne(strings.NewReader(raw), email)
	assert.NoError(t, err)

	assert.Contains(t, email.Text, "Hello world")
	assert.Equal(t, "sender@send.com", email.From)
	assert.Equal(t, []string{"b@b.com"}, email.To)
	assert.Equal(t, "hi", email.Subject)
	assert.Equal(t, "Mon, 07 Apr 2025 23:43:15 +0300", email.Date)
	assert.Len(t, email.Files, 0)
}
