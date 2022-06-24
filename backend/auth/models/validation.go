package models

import (
	"net/mail"
	"regexp"

	"github.com/go-playground/validator"
)

func validateNameLogin(fieldValidator validator.FieldLevel) bool {
	return true
}

func validateNameRegistration(fieldValidator validator.FieldLevel) bool {
	// requiring User name to be min 4 letters and contain only alphanumeric characters
	regex := regexp.MustCompile(`[a-zA-Z0-9]{4,}`)
	matches := regex.FindAllString(fieldValidator.Field().String(), 1)
	return len(matches) == 1
}

func validateEmail(fieldValidator validator.FieldLevel) bool {
	_, err := mail.ParseAddress(fieldValidator.Field().String())
	return err == nil
}

func validatePassword(fieldValidator validator.FieldLevel) bool {
	// password has to have one upper case, lower case, number and special character with min length 6
	// regex := regexp.MustCompile(`(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,}`)
	// matches := regex.FindAllString(fieldValidator.Field().String(), 1)
	// return len(matches) == 1
	return true // TODO fix regex (re2)
}
