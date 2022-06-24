package models

import (
	"encoding/json"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type Recipe struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" validate:"required,title" gorm:"index"`
	Description string    `json:"description" validate:"required"`
	Steps       []Step    `json:"steps" gorm:"ForeignKey:RecipeId"`
	UserId      uuid.UUID `json:"ownerId" gorm:"index"`
	Views       int64     `json:"views"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type Recipes []*Recipe

// serialization
func (recipe *Recipes) Serialize(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(recipe)
}

func (recipe *Recipe) Serialize(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(recipe)
}

func (recipe *Recipe) Deserialize(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(recipe)
}

func (recipe *Recipe) Validate() error {
	validator := validator.New()
	validator.RegisterValidation("title", validateTitle)
	return validator.Struct(recipe)
}

func validateTitle(fieldValidator validator.FieldLevel) bool {
	regex := regexp.MustCompile(`[a-zA-Z0-9]{4,}`)
	matches := regex.FindAllString(fieldValidator.Field().String(), 1)

	return len(matches) == 1
}
