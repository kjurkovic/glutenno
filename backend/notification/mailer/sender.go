package mailer

import (
	"log"
	"notifications/config"
	"notifications/models"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Sender struct {
	config *config.Mailer
}

func SenderFunc(config *config.Mailer) *Sender {
	return &Sender{
		config: config,
	}
}

func (sender *Sender) Mail(message *models.Message) {

	from := mail.NewEmail(sender.config.From.Name, sender.config.From.Email)

	recipient := mail.NewEmail(message.To, message.Email)
	plainTextContent := message.Text
	htmlContent := message.Text
	body := mail.NewSingleEmail(from, message.Subject, recipient, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sender.config.Key)
	_, err := client.Send(body)

	if err != nil {
		log.Println(err)
	}
}
