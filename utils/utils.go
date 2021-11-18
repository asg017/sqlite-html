package utils

import (
	"fmt"
	"html"
	"strings"

	"go.riyazali.net/sqlite"
)

func htmlEscape(s string) string {
	return html.EscapeString(s)
}

type HtmlEscapeFunc struct{}

func (*HtmlEscapeFunc) Deterministic() bool { return true }
func (*HtmlEscapeFunc) Args() int           { return 1 }
func (*HtmlEscapeFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	c.ResultText(htmlEscape(values[0].Text()))
}

func htmlUnescape(s string) string {
	return html.UnescapeString(s)
}

type HtmlUnescapeFunc struct{}

func (*HtmlUnescapeFunc) Deterministic() bool { return true }
func (*HtmlUnescapeFunc) Args() int           { return 1 }
func (*HtmlUnescapeFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	c.ResultText(html.UnescapeString(values[0].Text()))
}

func htmlTrim(s string) string {
	return strings.TrimSpace(s)
}

type HtmlTrimFunc struct{}

func (*HtmlTrimFunc) Deterministic() bool { return true }
func (*HtmlTrimFunc) Args() int           { return 1 }
func (*HtmlTrimFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	s := values[0].Text()
	c.ResultText(htmlTrim(s))
}

type HtmlTableFunc struct{}

func (*HtmlTableFunc) Deterministic() bool { return true }
func (*HtmlTableFunc) Args() int           { return 1 }
func (*HtmlTableFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	s := values[0].Text()
	c.ResultText(fmt.Sprintf("<table>%s", s))
}
