package tg

import (
	"strings"
	"text/template"

	"github.com/un1uckyyy/email-in-tg/internal/domain/models"

	tele "gopkg.in/telebot.v4"

	"github.com/microcosm-cc/bluemonday"
)

var p = bluemonday.NewPolicy()

func init() {
	p.AllowElements("b", "strong", "i", "em", "code", "s", "strike", "del", "u", "pre")
	p.AllowAttrs("href").OnElements("a")
}

// TODO rework it
func cleanTelegramHTML(input string) string {
	html := p.Sanitize(input)

	return html
}

var (
	helpTmpl  = template.Must(template.New("help").Parse(helpTemplate))
	loginTmpl = template.Must(template.New("login").Parse(loginTemplate))
	emailTmpl = template.Must(template.New("email").Parse(emailTemplate))
)

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

// TODO split messages with max telegramMessageLenLimit each
// nolint
func splitTextToMessages(text string) []string {
	return strings.Split(text, "\n")
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
