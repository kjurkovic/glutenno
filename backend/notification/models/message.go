package models

import (
	"encoding/json"
	"io"
)

type Message struct {
	To      string `json:"to"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

func (message *Message) Deserialize(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(message)
}
