package router

import (
	"auth/config"
	"auth/database"
	"auth/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type UserRouter struct {
	log    *log.Logger
	config *config.Config
}

func (wr *UserRouter) pathPrefix() string {
	return "/user"
}

func (wr *UserRouter) setup(router *mux.Router) {
	authHandler := handlers.Auth(wr.log, &database.UserDao{Conn: *database.Instance}, wr.config)
	userHandler := handlers.User(wr.log, &database.UserDao{Conn: *database.Instance})

	get := router.Methods(http.MethodGet).Subrouter()
	get.HandleFunc("", userHandler.GetUser)
	get.Use(authHandler.MiddlewareAuthorization)
}
