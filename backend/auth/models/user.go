package models

import (
	"encoding/json"
	"io"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID `json:"id" gorm:"primaryKey"`
	Name                string    `json:"name" validate:"name"`
	Email               string    `json:"email" validate:"required,email"`
	Password            string    `json:"password,omitempty" validate:"required,password"`
	RefreshToken        string    `json:"-"`
	ForgotPasswordToken string    `json:"-"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

func (user *User) Serialize(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(user)
}

func (user *User) Deserialize(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(user)
}

func (user *User) ValidateRegister() error {
	validator := validator.New()
	validator.RegisterValidation("name", validateNameRegistration)
	validator.RegisterValidation("email", validateEmail)
	validator.RegisterValidation("password", validatePassword)
	return validator.Struct(user)
}

func (user *User) ValidateLogin() error {
	validator := validator.New()
	validator.RegisterValidation("name", validateNameLogin)
	validator.RegisterValidation("email", validateEmail)
	validator.RegisterValidation("password", validatePassword)
	return validator.Struct(user)
}
