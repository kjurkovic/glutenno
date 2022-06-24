package handlers

import (
	"auth/errors"
	"auth/models"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func (handler *AuthHandler) MiddlewareAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		authorizationHeader := r.Header.Get("authorization")
		authorizationHeader = strings.TrimPrefix(authorizationHeader, "Bearer ")

		if len(authorizationHeader) == 0 {
			errors.WrongCredentials.SendErrorResponse(rw, http.StatusUnauthorized)
			return
		}

		token, _ := jwt.ParseWithClaims(authorizationHeader, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(handler.config.Authentication.SecretKey), nil
		})

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), KeyClaims{}, claims)
			next.ServeHTTP(rw, r.WithContext(ctx))
		} else {
			errors.WrongCredentials.SendErrorResponse(rw, http.StatusUnauthorized)
		}
	})
}

func (handler *AuthHandler) MiddlewareUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := models.User{}

		err := user.Deserialize(r.Body)
		if err != nil {
			handler.logger.Println("[ERROR] deserializing user", err)
			errors.SerializationError.SendErrorResponse(rw, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func (handler *AuthHandler) MiddlewareLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(KeyUser{}).(models.User)

		err := user.ValidateLogin()
		if err != nil {
			handler.logger.Println("[ERROR] validating user", err)
			errors.UserValidationError.SendErrorResponse(rw, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func (handler *AuthHandler) MiddlewareRegistration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(KeyUser{}).(models.User)

		err := user.ValidateRegister()
		if err != nil {
			handler.logger.Println("[ERROR] validating user", err)
			errors.UserValidationError.SendErrorResponse(rw, http.StatusBadRequest)
			return
		}

		_, err = handler.dao.GetByEmail(user.Email)

		if err != gorm.ErrRecordNotFound {
			handler.logger.Println("[ERROR] user already exists", user.Email)
			errors.UserAlreadyExistError.SendErrorResponse(rw, http.StatusConflict)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
