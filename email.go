package main

import (
	"os"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

// Email struct to be used to send emails
type Email struct {
	Subject    string
	Sender     string
	Message    string
	Recipients []string
}

func email(data Email) error {
	domain := os.Getenv("MAILGUN_DOMAIN")
	key := os.Getenv("MAILGUN_API_KEY")
	publicKey := os.Getenv("MAILGUN_PUBLIC_KEY")
	sender := os.Getenv("MAILGUN_SENDER")

	mg := mailgun.NewMailgun(domain, key, publicKey)

	message := mailgun.NewMessage(
		sender,
		data.Subject,
		data.Message,
		data.Recipients...,
	)

	message.SetHtml(data.Message)
	s, id, err := mg.Send(message)
	logger.Println(s, id)
	return err
}
