package main

import (
	"context"
	"log"
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
	defaultFrom, defaultTemplate := "no-reply@lfpanels.com", "default.html"

	lambda.Start(func(ctx context.Context, snsEvent SNSEvent) {
		for _, record := range snsEvent.Records {
			log.Printf("Got attributes %#v", record.SNS.MessageAttributes)

			var from = defaultFrom
			fromAttr, ok := record.SNS.MessageAttributes["From"]
			if ok {
				from = fromAttr.Value
			}

			var to string
			toAttr, ok := record.SNS.MessageAttributes["To"]
			if ok {
				to = toAttr.Value
			}

			if len(to) == 0 {
				log.Printf("missing To attribute for SNS message %#v", record.SNS)
				continue
			}

			var subject string
			subjectAttr, ok := record.SNS.MessageAttributes["Subject"]
			if ok {
				subject = subjectAttr.Value
			}

			var template = defaultTemplate
			templateAttr, ok := record.SNS.MessageAttributes["Template"]
			if ok {
				template = templateAttr.Value
			}

			log.Printf("Layout %s", template)

			err := email.Send(from, to, subject, template, map[string]interface{}{
				"Title":   "Test email",
				"Message": "message",
			})
			if err != nil {
				log.Println("send email error:", err)
			}
		}
	})
}
