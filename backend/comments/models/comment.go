package models

import (
	"encoding/json"
	"io"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID              uuid.UUID `json:"id" gorm:"primaryKey"`
	Text            string    `json:"text" validate:"required"`
	ResourceOwnerId uuid.UUID `json:"resourceOwnerId" validate:"required"`
	ResourceId      uuid.UUID `json:"resourceId" validate:"required"`
	User            string    `json:"user"`
	UserId          uuid.UUID `json:"userId" gorm:"index"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

type Comments []*Comment

// serialization
func (comment *Comments) Serialize(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(comment)
}

func (comment *Comment) Serialize(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(comment)
}

func (comment *Comment) Deserialize(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(comment)
}
