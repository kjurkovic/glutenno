package handlers

import (
	"log"
	"net/http"
	"notifications/mailer"
	"notifications/models"
)

type MessageHandler struct {
	logger *log.Logger
	sender *mailer.Sender
}

func MessageHandlerFunc(logger *log.Logger, sender *mailer.Sender) *MessageHandler {
	return &MessageHandler{
		logger: logger,
		sender: sender,
	}
}

func (handler *MessageHandler) Post(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Print("New notification received")

	message := &models.Message{}
	message.Deserialize(r.Body)

	handler.sender.Mail(message)
}
