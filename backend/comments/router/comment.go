package router

import (
	"comments/config"
	"comments/database"
	"comments/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var uuidRegex = "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}"

type CommentRouter struct {
	log    *log.Logger
	config *config.Config
}

func (cr *CommentRouter) pathPrefix() string {
	return "/comments"
}

func (cr *CommentRouter) setup(router *mux.Router) {
	ch := handlers.Comments(cr.log, &database.CommentsDao{Conn: *database.Instance}, cr.config)

	// GET router
	get := router.Methods(http.MethodGet).Subrouter()
	get.HandleFunc(fmt.Sprintf("/{id:%s}", uuidRegex), ch.Get)

	getProtected := router.Methods(http.MethodGet).Subrouter()
	getProtected.HandleFunc("/user", ch.GetPerUserRecipe)
	getProtected.Use(ch.MiddlewareAuthorization)

	// PUT router
	put := router.Methods(http.MethodPut).Subrouter()
	put.HandleFunc(fmt.Sprintf("/{id:%s}", uuidRegex), ch.Update)
	put.Use(ch.MiddlewareAuthorization, ch.MiddlewareComment)

	// POST router
	post := router.Methods(http.MethodPost).Subrouter()
	post.HandleFunc("", ch.Add)
	post.Use(ch.MiddlewareAuthorization, ch.MiddlewareComment)

	// DELETE router
	delete := router.Methods(http.MethodDelete).Subrouter()
	delete.HandleFunc(fmt.Sprintf("/{id:%s}", uuidRegex), ch.Delete)
	delete.Use(ch.MiddlewareAuthorization)
}
