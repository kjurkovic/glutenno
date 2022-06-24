package models

import (
	"bytes"
	"encoding/json"
)

type Message struct {
	To      string `json:"to"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

func (message *Message) Serialize() (*bytes.Buffer, error) {
	data, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(data), nil
}
