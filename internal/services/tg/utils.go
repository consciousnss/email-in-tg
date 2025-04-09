package tg

import "github.com/microcosm-cc/bluemonday"

var p = bluemonday.NewPolicy()

func init() {
	p.AllowElements("b", "strong", "i", "em", "code", "s", "strike", "del", "u", "pre")
	p.AllowAttrs("href").OnElements("a")
}

func cleanTelegramHTML(input string) string {
	html := p.Sanitize(input)

	return html
}
