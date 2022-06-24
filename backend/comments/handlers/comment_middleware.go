package handlers

import (
	"comments/errors"
	"comments/models"
	"context"
	"net/http"
)

func (handler *CommentHandler) MiddlewareAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")

		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodGet, handler.config.Authorization.Address, nil)
		req.Header.Set("Authorization", authorizationHeader)

		res, err := client.Do(req)

		if err != nil {
			errors.ServerError.SendErrorResponse(rw, http.StatusInternalServerError)
			return
		}

		if res.StatusCode == http.StatusOK {
			user := models.User{}
			err = user.Deserialize(res.Body)

			if err != nil {
				errors.SerializationError.SendErrorResponse(rw, http.StatusBadGateway)
				return
			}

			ctx := context.WithValue(r.Context(), models.UserKey{}, user)
			next.ServeHTTP(rw, r.WithContext(ctx))
		} else {
			errors.WrongCredentials.SendErrorResponse(rw, http.StatusUnauthorized)
		}
	})
}

func (handler *CommentHandler) MiddlewareComment(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		comment := models.Comment{}

		err := comment.Deserialize(r.Body)
		if err != nil {
			handler.logger.Println("[ERROR] deserializing", err)
			http.Error(rw, "Unable to parse request", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyComment{}, comment)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
