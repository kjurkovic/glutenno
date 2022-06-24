package handlers

import (
	"auth/config"
	"auth/database"
	"auth/errors"
	"auth/models"
	"auth/utils"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"crypto/sha256"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// types
type AuthHandler struct {
	logger *log.Logger
	dao    *database.UserDao
	config *config.Config
}

type KeyUser struct{}
type KeyClaims struct{}

func Auth(l *log.Logger, dao *database.UserDao, config *config.Config) *AuthHandler {
	return &AuthHandler{
		logger: l,
		dao:    dao,
		config: config,
	}
}

// HTTP handlers

// swagger:route POST /auth/register auth register
// Returns auth object
// responses:
// 	200: authResponse
func (handler *AuthHandler) Register(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Print("POST Auth Register")

	rw.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(KeyUser{}).(models.User)
	user.ID = uuid.New()
	user.Password = utils.GeneratePasswordHashWithSalt(user.Password)

	handler.logger.Printf("Received object: %v", user)

	response := handler.generateAuthResponse(&user, rw)
	user.RefreshToken = response.RefreshToken
	_, err := handler.dao.Insert(&user)

	if err != nil {
		handler.logger.Println(err)
		http.Error(rw, "Unable to add user", http.StatusInternalServerError)
		return
	}

	response.Serialize(rw)
}

func (handler *AuthHandler) Login(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Print("POST Auth Login")

	user := r.Context().Value(KeyUser{}).(models.User)
	storedUser, err := handler.dao.GetByEmail(user.Email)

	if err != nil {
		errors.WrongCredentials.SendErrorResponse(rw, http.StatusForbidden)
		return
	}

	isUserAuthorized := handler.validatePassword(storedUser, &user)

	if !isUserAuthorized {
		errors.WrongCredentials.SendErrorResponse(rw, http.StatusForbidden)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	handler.generateAuthResponse(storedUser, rw).Serialize(rw)
}

func (handler *AuthHandler) RefreshToken(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Print("POST Auth Refresh Token")

	refreshToken := &models.RefreshToken{}
	err := refreshToken.Deserialize(r.Body)

	if err != nil {
		errors.SerializationError.SendErrorResponse(rw, http.StatusBadGateway)
		return
	}

	user, err := handler.dao.GetUserByRefreshToken(refreshToken.Token)

	if err == gorm.ErrRecordNotFound {
		errors.WrongCredentials.SendErrorResponse(rw, http.StatusForbidden)
		return
	} else if err != nil {
		errors.ServerError.SendErrorResponse(rw, http.StatusInternalServerError)
		return
	}

	response := handler.generateAuthResponse(user, rw)
	rw.Header().Set("Content-Type", "application/json")
	response.Serialize(rw)
}

func (handler *AuthHandler) ForgetPassword(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Print("POST Auth Forget Password")

	requestBody := &models.ForgetPasswordRequest{}
	err := requestBody.Deserialize(r.Body)

	if err != nil {
		errors.SerializationError.SendErrorResponse(rw, http.StatusBadGateway)
		return
	}

	user, err := handler.dao.GetByEmail(requestBody.Email)

	if err == nil {
		randomToken := sha256.Sum256([]byte(utils.RandomString(20)))
		forgotPasswordToken := hex.EncodeToString(randomToken[:])
		handler.dao.UpdateForgotPasswordToken(user.ID, forgotPasswordToken)

		err = handler.sendMail(&models.Message{
			To:      user.Name,
			Email:   user.Email,
			Subject: "Forgot password",
			Text:    fmt.Sprintf("To reset your password click on the link below:\\n <a href=\"%s%s%s\">Reset password</a>", handler.config.Frontend.Url, handler.config.Frontend.ResetPassword, forgotPasswordToken),
		})
	}

	if err != nil {
		fmt.Print(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent)
}

func (handler *AuthHandler) ResetPassword(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Print("POST Auth Reset Password")

	requestBody := &models.ResetPasswordRequest{}
	err := requestBody.Deserialize(r.Body)

	if err != nil {
		errors.SerializationError.SendErrorResponse(rw, http.StatusBadGateway)
		return
	}

	user, err := handler.dao.GetUserByForgotPasswordToken(requestBody.Token)

	if err == gorm.ErrRecordNotFound {
		errors.WrongCredentials.SendErrorResponse(rw, http.StatusForbidden)
		return
	} else if err != nil {
		errors.ServerError.SendErrorResponse(rw, http.StatusInternalServerError)
		return
	}

	handler.dao.UpdatePassword(user.ID, utils.GeneratePasswordHashWithSalt(requestBody.Password))

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent)
}

func (handler *AuthHandler) generateAuthResponse(user *models.User, rw http.ResponseWriter) *models.Auth {
	expirationTime := time.Now().Add(handler.config.Authentication.AccessTokenExpiration * time.Minute)
	claims := &models.Claims{
		Username: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(handler.config.Authentication.SecretKey))

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	response := &models.Auth{
		AccessToken:  tokenString,
		ExpiresIn:    expirationTime.UnixMilli(),
		RefreshToken: handler.generateRefreshToken(user),
	}
	return response
}

func (handler *AuthHandler) generateRefreshToken(user *models.User) string {
	var sb strings.Builder
	sb.WriteString(user.ID.String())
	sb.WriteString(handler.config.Authentication.SecretKey)
	sb.WriteString(time.Now().GoString())
	hash := sha256.Sum256([]byte(sb.String()))
	refreshToken := hex.EncodeToString(hash[:])
	handler.dao.UpdateRefreshToken(user.ID, refreshToken)
	return refreshToken
}

func (handler *AuthHandler) validatePassword(storedUser *models.User, user *models.User) bool {
	storedEncryptedPassword := storedUser.Password
	components := strings.Split(storedEncryptedPassword, ":")
	salt := components[0]
	encryptedPassword := components[1]

	requestPasswordHash := utils.CreateHash(salt, user.Password)
	return encryptedPassword == requestPasswordHash
}

func (handler *AuthHandler) sendMail(message *models.Message) error {
	buffer, err := message.Serialize()

	if err != nil {
		return err
	}

	handler.logger.Print("Sending email", handler.config.Notification.Address, message)
	_, err = http.Post(handler.config.Notification.Address, "application/json", buffer)

	return err
}
