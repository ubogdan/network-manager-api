package email

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/ubogdan/network-manager-api/model"
)

func Notify(from, to, subject string, email model.Email, t Template) error {

	emailBody, err := t.GenerateHTML(email)
	if err != nil {
		return fmt.Errorf("generate html: %w", err)
	}

	emailText, err := t.GeneratePlainText(email)
	if err != nil {
		return fmt.Errorf("generate plaintext: %w", err)
	}

	auth, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}

	_, err = ses.New(auth).SendEmail(&ses.SendEmailInput{
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
					Data: aws.String(emailBody),
				},
				Text: &ses.Content{
					Data: aws.String(emailText),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
