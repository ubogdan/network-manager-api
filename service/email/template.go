package email

import (
	"bytes"
	"html/template"

	"github.com/Masterminds/sprig"
	"github.com/imdario/mergo"
	"github.com/jaytaylor/html2text"
	"github.com/ubogdan/network-manager-api/model"
	"github.com/vanng822/go-premailer/premailer"
)

// Generator email generator.
type Template struct {
	Theme              Theme
	Product            Product
	TextDirection      TextDirection
	DisableCSSInlining bool
}

// Theme is an interface to implement when creating a new theme
type Theme interface {
	Name() string              // The name of the theme
	HTMLTemplate() string      // The golang template for HTML emails
	PlainTextTemplate() string // The golang templte for plain text emails (can be basic HTML)
}

// Product represents your company product (brand)
// Appears in header & footer of e-mails
type Product struct {
	Name        string
	Link        string // e.g. https://matcornic.github.io
	Logo        string // e.g. https://matcornic.github.io/img/logo.png
	Copyright   string // Copyright © 2019 Hermes. All rights reserved.
	TroubleText string // TroubleText is the sentence at the end of the email for users having trouble with the button (default to `If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser.`)
}

var templateFuncs = template.FuncMap{
	"url": func(s string) template.URL {
		return template.URL(s)
	},
}

// TextDirection of the text in HTML email
type TextDirection string

// TDLeftToRight is the text direction from left to right (default)
const TDLeftToRight TextDirection = "ltr"

// TDRightToLeft is the text direction from right to left
const TDRightToLeft TextDirection = "rtl"

// GeneratePlainText generates the email body from data
// This is for old email clients
func (t *Template) GeneratePlainText(email model.Email) (string, error) {
	err := setDefaultValues(t)
	if err != nil {
		return "", err
	}
	template, err := t.generateTemplate(email, t.Theme.PlainTextTemplate())
	if err != nil {
		return "", err
	}
	return html2text.FromString(template, html2text.Options{PrettyTables: true})
}

// GenerateHTML generates the email body from data to an HTML Reader
// This is for modern email clients
func (t *Template) GenerateHTML(email model.Email) (string, error) {
	err := setDefaultValues(t)
	if err != nil {
		return "", err
	}
	return t.generateTemplate(email, t.Theme.HTMLTemplate())
}

// default values of the engine
func setDefaultValues(h *Template) error {
	defaultHermes := Template{
		Theme:         new(Default),
		TextDirection: TDLeftToRight,
		Product: Product{
			Name:        "Network Manager",
			Copyright:   "Copyright © 2020 Light Fast Panels. All rights reserved.",
			TroubleText: "If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser.",
		},
	}

	// Merge the given hermes engine configuration with default one
	// Default one overrides all zero values
	err := mergo.Merge(h, defaultHermes)
	if err != nil {
		return err
	}

	if h.TextDirection != TDLeftToRight && h.TextDirection != TDRightToLeft {
		h.TextDirection = TDLeftToRight
	}

	return nil
}

func (t *Template) generateTemplate(email model.Email, tplt string) (string, error) {
	defaultEmail := model.Email{
		Intros:     []string{},
		Dictionary: []model.Entry{},
		Outros:     []string{},
		Signature:  "Yours truly",
		Greeting:   "Hi",
	}
	err := mergo.Merge(&email, defaultEmail)
	if err != nil {
		return "", err
	}

	// Generate the email from Golang template
	// Allow usage of simple function from sprig : https://github.com/Masterminds/sprig
	tpl, err := template.New("template").Funcs(sprig.FuncMap()).Funcs(templateFuncs).Funcs(template.FuncMap{
		"safe": func(s string) template.HTML { return template.HTML(s) }, // Used for keeping comments in generated template
	}).Parse(tplt)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = tpl.Execute(&b, struct {
		Product       Product
		TextDirection TextDirection
		Email         model.Email
	}{
		Product:       t.Product,
		Email:         email,
		TextDirection: t.TextDirection,
	})
	if err != nil {
		return "", err
	}

	res := b.String()
	if t.DisableCSSInlining {
		return res, nil
	}

	// Inlining CSS
	prem, err := premailer.NewPremailerFromString(res, premailer.NewOptions())
	if err != nil {
		return "", err
	}
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}
	return html, nil
}
