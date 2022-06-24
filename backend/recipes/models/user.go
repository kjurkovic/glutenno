package models

import (
	"encoding/json"
	"io"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}

type UserKey struct{}

func (user *User) Deserialize(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(user)
}
