package tg

import (
	"strings"
	"text/template"

	"github.com/k3a/html2text"

	"github.com/un1uckyyy/email-in-tg/internal/domain/models"

	tele "gopkg.in/telebot.v4"
)

// TODO add unit tests
func html2Text(input string) string {
	return html2text.HTML2TextWithOptions(input, html2text.WithLinksInnerText())
}

var (
	helpTmpl  = template.Must(template.New("help").Parse(helpTemplate))
	loginTmpl = template.Must(template.New("login").Parse(loginTemplate))
	emailTmpl = template.Must(template.New("email").Parse(emailTemplate))
)

// TODO update unit tests
func renderHTMLTemplate(tmpl *template.Template, data any) (string, error) {
	var builder strings.Builder
	if err := tmpl.Execute(&builder, data); err != nil {
		return "", err
	}
	return builder.String(), nil
}

const (
	telegramMessageLenLimit = 4096
	telegramAlbumMediaLimit = 10
)

// TODO add unit tests
func splitTextToMessages(text string, limit int) []string {
	msgCount := len(text) / limit
	if len(text)%limit != 0 {
		msgCount++
	}

	messages := make([]string, 0, msgCount)
	for i := 0; i < msgCount; i++ {
		messages = append(messages, text[i*limit:(i+1)*limit])
	}

	return messages
}

func splitFilesToAlbums(files []*models.File) []tele.Album {
	albumsNum := (len(files) + telegramAlbumMediaLimit - 1) / telegramAlbumMediaLimit
	albums := make([]tele.Album, 0, albumsNum)

	for i := 0; i < albumsNum; i++ {
		start, end := i*telegramAlbumMediaLimit, min(len(files), (i+1)*telegramAlbumMediaLimit)

		album := make(tele.Album, 0, end-start)

		for j := start; j < end; j++ {
			file := files[j]
			album = append(album,
				&tele.Document{
					File:     tele.FromReader(file.Data),
					FileName: file.Filename,
				},
			)
		}

		albums = append(albums, album)
	}

	return albums
}
