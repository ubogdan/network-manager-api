package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"gopkg.in/gomail.v2"
)

//go:embed templates/*.html
var templateFS embed.FS

func Send(from, to, subject, template string, data interface{}) error {
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}
	service := ses.New(sess)

	var buf bytes.Buffer

	err = Template(template, data, &buf)
	if err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	_, err = service.SendEmail(&ses.SendEmailInput{
		Source: aws.String(from),
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(buf.String()),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// SendRaw godoc.
func SendRaw(from, to, subbect, template string, data interface{}) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	service := ses.New(sess)

	mail := gomail.NewMessage()
	mail.SetHeader("From", from)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subbect)

	var Buff bytes.Buffer

	err = Template(template, data, &Buff)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	mail.SetBody("text/html", Buff.String())

	var rawMessage bytes.Buffer

	_, err = mail.WriteTo(&rawMessage)
	if err != nil {
		return fmt.Errorf("mail.WriteTo failed: %w", err)
	}

	_, err = service.SendRawEmail(&ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: rawMessage.Bytes(),
		},
	})
	if err != nil {
		return fmt.Errorf("ses.SendRawEmail failed %w", err)
	}

	return nil
}

// Template render a html template into an io.Writer.
func Template(layout string, data interface{}, writer io.Writer) error {
	tmpl, err := template.New("").Funcs(
		template.FuncMap{
			"ToLower": strings.ToLower,
			"ToUpper": strings.ToUpper,
		}).ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return fmt.Errorf("template.ParseFS error: %w", err)
	}

	err = tmpl.ExecuteTemplate(writer, layout, data)
	if err != nil {
		return fmt.Errorf("template.ExecuteTemplate error %w", err)
	}

	return nil
}
