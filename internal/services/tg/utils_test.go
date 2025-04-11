package tg

import (
	"bytes"
	"io"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/un1uckyyy/email-in-tg/internal/models"
	tele "gopkg.in/telebot.v4"
)

func mockFile(name string) *models.File {
	return &models.File{
		Filename: name,
		Data:     bytes.NewBufferString("mock content"),
	}
}

func TestSplitFilesToAlbums(t *testing.T) {
	tests := []struct {
		name       string
		fileCount  int
		wantAlbums int
	}{
		{"no files", 0, 0},
		{"less than limit", 5, 1},
		{"exactly one album", telegramAlbumMediaLimit, 1},
		{"slightly over limit", telegramAlbumMediaLimit + 1, 2},
		{"multiple full albums", telegramAlbumMediaLimit * 3, 3},
		{"multiple and remainder", telegramAlbumMediaLimit*2 + 3, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := make([]*models.File, 0, tt.fileCount)
			for i := 0; i < tt.fileCount; i++ {
				files = append(files, mockFile("file"+string(rune(i))))
			}

			albums := splitFilesToAlbums(files)

			assert.Len(t, albums, tt.wantAlbums)

			// check that all files are present
			totalMedia := 0
			for _, album := range albums {
				totalMedia += len(album)
				for _, media := range album {
					doc, ok := media.(*tele.Document)
					assert.True(t, ok)
					assert.Implements(t, (*io.Reader)(nil), doc.FileReader)
				}
			}
			assert.Equal(t, tt.fileCount, totalMedia)
		})
	}
}

func TestRenderHTMLTemplate_EmailTemplate(t *testing.T) {
	email := &models.Email{
		MailFrom: "alice@example.com",
		MailTo:   "bob@example.com",
		Date:     "2025-04-10",
		Subject:  "Hello!",
		Text:     "This is a test email.",
	}

	output, err := renderHTMLTemplate(emailTmpl, email)
	assert.NoError(t, err)
	assert.Contains(t, output, "alice@example.com")
	assert.Contains(t, output, "Hello!")
	assert.Contains(t, output, "This is a test email.")
	assert.Contains(t, output, "✉️ <b>Новое письмо</b> ✉️")
}

func TestRenderHTMLTemplate_HelpTemplate(t *testing.T) {
	output, err := renderHTMLTemplate(helpTmpl, "support@example.com")
	assert.NoError(t, err)
	assert.Contains(t, output, "<b>/start 'email' 'password'</b>")
	assert.Contains(t, output, "<i>support:</i> support@example.com.")
}

func TestRenderHTMLTemplate_InvalidTemplate(t *testing.T) {
	invalidTmpl := template.Must(template.New("bad").Parse("{{.NonexistentField}}"))

	output, err := renderHTMLTemplate(invalidTmpl, "just a string")

	assert.Error(t, err)
	assert.Empty(t, output)
}
