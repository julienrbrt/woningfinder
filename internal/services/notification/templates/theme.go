package templates

import _ "embed"

//go:embed woningfinder-html.tpl
var htmlTpl string

//go:embed woningfinder-text.tpl
var textTpl string

type theme struct{}

func (t *theme) Name() string {
	return "WoningFinder"
}

// HTMLTemplate returns a Golang template that will generate an HTML email.
func (t *theme) HTMLTemplate() string {
	return htmlTpl
}

// PlainTextTemplate returns a Golang template that will generate an plain text email.
func (t *theme) PlainTextTemplate() string {
	return textTpl
}
