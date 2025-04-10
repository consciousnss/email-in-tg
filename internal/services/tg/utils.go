package tg

import (
	"strings"
	"text/template"

	"github.com/microcosm-cc/bluemonday"
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

var (
	helpTmpl  = template.Must(template.New("help").Parse(helpTemplate))
	emailTmpl = template.Must(template.New("email").Parse(emailTemplate))
)

func renderHTMLTemplate(tmpl *template.Template, data any) (string, error) {
	var builder strings.Builder
	if err := tmpl.Execute(&builder, data); err != nil {
		return "", err
	}
	return builder.String(), nil
}
