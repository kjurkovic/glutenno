package models

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
)

type Auth struct {
	AccessToken  string `json:"accessToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshToken struct {
	Token string `json:"refreshToken"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password" validate:"required,password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (auth *Auth) Serialize(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(auth)
}

func (token *RefreshToken) Serialize(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(token)
}

func (token *RefreshToken) Deserialize(reader io.Reader) error {
	encoder := json.NewDecoder(reader)
	return encoder.Decode(token)
}

func (token *ForgetPasswordRequest) Deserialize(reader io.Reader) error {
	encoder := json.NewDecoder(reader)
	return encoder.Decode(token)
}

func (token *ResetPasswordRequest) Deserialize(reader io.Reader) error {
	encoder := json.NewDecoder(reader)
	return encoder.Decode(token)
}

func (request *ForgetPasswordRequest) ValidateForgotPassword() error {
	validator := validator.New()
	validator.RegisterValidation("email", validateEmail)
	return validator.Struct(request)
}

func (request *ResetPasswordRequest) ValidateResetPassword() error {
	validator := validator.New()
	validator.RegisterValidation("password", validatePassword)
	return validator.Struct(request)
}
