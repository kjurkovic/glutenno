package router

import (
	"fmt"
	"log"
	"net/http"
	"recipes/config"
	"recipes/database"
	"recipes/handlers"

	"github.com/gorilla/mux"
)

var recipeIdRegex = "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}"

type RecipeRouter struct {
	log    *log.Logger
	config *config.Authorization
}

func (rr *RecipeRouter) pathPrefix() string {
	return "/recipes"
}

func (rr *RecipeRouter) setup(router *mux.Router) {
	rh := handlers.Recipes(rr.log, &database.RecipesDao{Conn: *database.Instance}, &database.StepDao{Conn: *database.Instance}, rr.config)

	// GET router
	get := router.Methods(http.MethodGet).Subrouter()
	get.HandleFunc("", rh.Get)
	get.HandleFunc(fmt.Sprintf("/{id:%s}", recipeIdRegex), rh.GetSingle)

	protectedGet := router.Methods(http.MethodGet).Subrouter()
	protectedGet.HandleFunc("/user", rh.GetUser)
	protectedGet.Use(rh.MiddlewareAuthorization)

	// PUT router
	protectedPut := router.Methods(http.MethodPut).Subrouter()
	protectedPut.HandleFunc(fmt.Sprintf("/{id:%s}", recipeIdRegex), rh.Update)
	protectedPut.Use(rh.MiddlewareAuthorization, rh.MiddlewareValidateRecipe)

	put := router.Methods(http.MethodPut).Subrouter()
	put.HandleFunc(fmt.Sprintf("/{id:%s}/view", recipeIdRegex), rh.UpdateView)

	// POST router
	post := router.Methods(http.MethodPost).Subrouter()
	post.HandleFunc("", rh.Add)
	post.Use(rh.MiddlewareAuthorization, rh.MiddlewareValidateRecipe)

	// DELETE router
	delete := router.Methods(http.MethodDelete).Subrouter()
	delete.HandleFunc(fmt.Sprintf("/{id:%s}", recipeIdRegex), rh.Delete)
	delete.Use(rh.MiddlewareAuthorization)
}
