package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

func CreateHash(salt string, password string) string {
	var sb strings.Builder
	sb.WriteString(salt)
	sb.WriteString(".")
	sb.WriteString(password)
	encryptedPassword := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(encryptedPassword[:])
}

func GeneratePasswordHashWithSalt(password string) string {
	var sb strings.Builder
	salt := sha256.Sum256([]byte(time.Now().GoString()))
	saltString := hex.EncodeToString(salt[:])
	encryptedPassword := CreateHash(saltString, password)
	sb.WriteString(saltString)
	sb.WriteString(":")
	sb.WriteString(encryptedPassword)
	return sb.String()
}
