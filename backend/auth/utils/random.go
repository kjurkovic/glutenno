package utils

import "math/rand"

var availableCharacters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*(){}|[]-_")

func RandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = availableCharacters[rand.Intn(len(availableCharacters))]
	}
	return string(b)
}
