package tg

import (
	"strings"
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

var emailTmpl = template.Must(template.New("email").Parse(emailTemplate))

func renderEmailTemplateHTML(email *models.Email) (string, error) {
	var builder strings.Builder
	if err := emailTmpl.Execute(&builder, email); err != nil {
		return "", err
	}
	return builder.String(), nil
}
