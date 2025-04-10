package tg

import (
	"bytes"
	"text/template"

	"github.com/microcosm-cc/bluemonday"
	"github.com/un1uckyyy/email-in-tg/internal/models"
)

var p = bluemonday.NewPolicy()

func init() {
	p.AllowElements("b", "strong", "i", "em", "code", "s", "strike", "del", "u", "pre")
	p.AllowAttrs("href").OnElements("a")
}

func cleanTelegramHTML(input string) string {
	html := p.Sanitize(input)

	return html
}

func renderEmailTemplate(email *models.Email) (string, error) {
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, email)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
