package handlers

import (
	"auth/database"
	"auth/errors"
	"auth/models"
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserHandler struct {
	logger *log.Logger
	dao    *database.UserDao
}

func User(l *log.Logger, dao *database.UserDao) *UserHandler {
	return &UserHandler{
		logger: l,
		dao:    dao,
	}
}

// middleware
func (handler *UserHandler) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := models.User{}

		err := user.Deserialize(r.Body)
		if err != nil {
			handler.logger.Println("[ERROR] deserializing user", err)
			errors.SerializationError.SendErrorResponse(rw, http.StatusBadRequest)
			return
		}

		err = user.ValidateRegister()
		if err != nil {
			handler.logger.Println("[ERROR] validating user", err)
			errors.UserValidationError.SendErrorResponse(rw, http.StatusBadRequest)
			return
		}

		existing, _ := handler.dao.GetByEmail(user.Email)
		if existing != nil {
			handler.logger.Println("[ERROR] user already exists", user.Email)
			errors.UserAlreadyExistError.SendErrorResponse(rw, http.StatusConflict)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

// REST methods

func (handler *UserHandler) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(mux.Vars(r)["id"])
	handler.logger.Printf("PUT Update user#id: %s", id.String())

	rw.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(KeyUser{}).(models.User)

	response, err := handler.dao.Update(id, &user)

	if err == gorm.ErrRecordNotFound {
		errors.UserNotFoundError.SendErrorResponse(rw, http.StatusNotFound)
		return
	} else if err != nil {
		errors.UserNotFoundError.SendErrorResponse(rw, http.StatusInternalServerError)
		return
	}

	err = response.Serialize(rw)

	if err != nil {
		errors.SerializationError.SendErrorResponse(rw, http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(KeyClaims{}).(*models.Claims)

	handler.logger.Printf("GET User#email: %s", claims.Username)
	user, err := handler.dao.GetByEmail(claims.Username)

	if err != nil {
		errors.UserNotFoundError.SendErrorResponse(rw, http.StatusBadGateway)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	user.Password = ""
	err = user.Serialize(rw)
	if err != nil {
		errors.SerializationError.SendErrorResponse(rw, http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(mux.Vars(r)["id"])
	handler.logger.Printf("DELETE User#id: %s", id.String())
	rw.Header().Set("Content-Type", "application/json")

	affectedRows, err := handler.dao.Delete(id)

	if err == gorm.ErrRecordNotFound || affectedRows == 0 {
		errors.UserNotFoundError.SendErrorResponse(rw, http.StatusNotFound)
		return
	} else if err != nil {
		errors.UserNotFoundError.SendErrorResponse(rw, http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
