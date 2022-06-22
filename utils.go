package main

import (
	"fmt"
	"html"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.riyazali.net/sqlite"
)

/**	html_valid(document)
 * Returns 1 if the given document is valid HTML, 0 otherwise.
 * @param document {text | html} - HTML document to read from.
 **/
 type HtmlValidFunc struct{}

 func (*HtmlValidFunc) Deterministic() bool { return true }
 func (*HtmlValidFunc) Args() int           { return 1 }
 func (*HtmlValidFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	 html := values[0].Text()
	 _, err := goquery.NewDocumentFromReader(strings.NewReader(html))
 
	 if err != nil {
		 c.ResultInt(0)
	 } else {
		 c.ResultInt(1)
	 }
 }

 /**	html_escape(content)
 * Returns an HTML escaped version of the given content.
 * @param content {text} - Text content to escape.
 **/
type HtmlEscapeFunc struct{}

func (*HtmlEscapeFunc) Deterministic() bool { return true }
func (*HtmlEscapeFunc) Args() int           { return 1 }
func (*HtmlEscapeFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	c.ResultText(html.EscapeString(values[0].Text()))
}

 /**	html_unescape(content)
 * Returns an HTML unescaped version of the given content.
 * @param content {text} - Text content to unescape.
 **/
 type HtmlUnescapeFunc struct{}

func (*HtmlUnescapeFunc) Deterministic() bool { return true }
func (*HtmlUnescapeFunc) Args() int           { return 1 }
func (*HtmlUnescapeFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	c.ResultText(html.UnescapeString(values[0].Text()))
}

 /**	html_trim(content)
 * Trim whitespace around the given text content. Useful for output of html_text
 * @param content {text} - Text content to trim.
 **/
type HtmlTrimFunc struct{}

func (*HtmlTrimFunc) Deterministic() bool { return true }
func (*HtmlTrimFunc) Args() int           { return 1 }
func (*HtmlTrimFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	s := values[0].Text()
	c.ResultText(strings.TrimSpace(s))
}


 /**	html_table(content)
 * Wrap the given content around a HTML table. Useful for parsing table rows. 
 * @param content {text} - Text content to wrap.
 **/
type HtmlTableFunc struct{}

func (*HtmlTableFunc) Deterministic() bool { return true }
func (*HtmlTableFunc) Args() int           { return 1 }
func (*HtmlTableFunc) Apply(c *sqlite.Context, values ...sqlite.Value) {
	s := values[0].Text()
	c.ResultText(fmt.Sprintf("<table>%s", s))
}

func RegisterUtils(api *sqlite.ExtensionApi) error {
	var err error
	if err = api.CreateFunction("html_valid", &HtmlValidFunc{}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_table", &HtmlTableFunc{}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_escape", &HtmlEscapeFunc{}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_unescape", &HtmlUnescapeFunc{}); err != nil {
		return err
	}
	if err = api.CreateFunction("html_trim", &HtmlTrimFunc{}); err != nil {
		return err
	}
	return nil
}
