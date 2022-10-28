package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ubogdan/network-manager-api/service/email"
)

type SNSEvent struct {
	Records []SNSEventRecord `json:"Records"`
}

type SNSEventRecord struct {
	EventVersion         string    `json:"EventVersion"`
	EventSubscriptionArn string    `json:"EventSubscriptionArn"`
	EventSource          string    `json:"EventSource"`
	SNS                  SNSEntity `json:"Sns"`
}

type SNSEntity struct {
	Signature         string                         `json:"Signature"`
	MessageID         string                         `json:"MessageId"`
	Type              string                         `json:"Type"`
	TopicArn          string                         `json:"TopicArn"`
	MessageAttributes map[string]SNSMessageAttribute `json:"MessageAttributes"`
	SignatureVersion  string                         `json:"SignatureVersion"`
	Timestamp         time.Time                      `json:"Timestamp"`
	SigningCertURL    string                         `json:"SigningCertUrl"`
	Message           string                         `json:"Message"`
	UnsubscribeURL    string                         `json:"UnsubscribeUrl"`
	Subject           string                         `json:"Subject"`
}

type SNSMessageAttribute struct {
	Type  string `json:"Type"`
	Value string `json:"Value,omitempty"`
}

func main() {
	defaultFrom := os.Getenv("DEFAULT_EMAIL_FROM")

	lambda.Start(func(ctx context.Context, snsEvent SNSEvent) {
		for _, record := range snsEvent.Records {
			log.Printf("Got attributes %#v", record.SNS.MessageAttributes)

			var from, to, subject, template string

			attributes := make(map[string]string)

			for name, attr := range record.SNS.MessageAttributes {
				switch name {
				case "From":
					from = attr.Value
				case "To":
					to = attr.Value
				case "Template":
					template = attr.Value
				case "Subject":
					subject = attr.Value
					fallthrough
				default:
					attributes[name] = attr.Value
				}
			}

			if len(from) == 0 {
				from = defaultFrom
			}

			if len(to) == 0 {
				log.Printf("missing To attribute for SNS message %#v", record.SNS.MessageID)

				continue
			}

			if len(template) == 0 {
				log.Printf("missing Template attribute for SNS message %#v", record.SNS.MessageID)

				continue
			}

			err := email.Send(from, to, subject, template, attributes)
			if err != nil {
				log.Println("send email error:", err)
			}
		}
	})
}
