package models

import (
	"encoding/json"
	"io"
	"time"

	"github.com/google/uuid"
)

type Step struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey"`
	Description string    `json:"description" validate:"required"`
	Order       int       `json:"order"`
	RecipeId    uuid.UUID `json:"-" gorm:"index,column:recipe_id"`
	Recipe      Recipe    `json:"-"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type Steps []*Step

// serialization
func (step *Steps) Serialize(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(step)
}

func (step *Step) Serialize(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(step)
}

func (step *Step) Deserialize(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(step)
}
