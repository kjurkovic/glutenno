package handlers

import (
	"fmt"
	"log"
	"net/http"
	"recipes/config"
	"recipes/database"
	"recipes/models"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// types
type RecipeHandler struct {
	logger  *log.Logger
	dao     *database.RecipesDao
	stepDao *database.StepDao
	config  *config.Authorization
}

type KeyRecipe struct{}

func Recipes(l *log.Logger, dao *database.RecipesDao, stepDao *database.StepDao, config *config.Authorization) *RecipeHandler {
	return &RecipeHandler{
		logger:  l,
		dao:     dao,
		stepDao: stepDao,
		config:  config,
	}
}

func (handler *RecipeHandler) Add(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Print("POST")

	rw.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(models.UserKey{}).(models.User)
	recipe := r.Context().Value(KeyRecipe{}).(models.Recipe)

	recipe.ID = uuid.New()
	recipe.UserId = user.ID

	for i := 0; i < len(recipe.Steps); i++ {
		step := &recipe.Steps[i]
		step.ID = uuid.New()
	}

	fmt.Printf("Received object: %v", recipe)
	_, err := handler.dao.Insert(&recipe)

	if err != nil {
		http.Error(rw, "Unable to add", http.StatusInternalServerError)
		return
	}

	err = recipe.Serialize(rw)

	if err != nil {
		http.Error(rw, "Serialization error", http.StatusInternalServerError)
		return
	}
}

func (handler *RecipeHandler) Update(rw http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(mux.Vars(r)["id"])
	handler.logger.Printf("PUT id: %s", id.String())

	rw.Header().Set("Content-Type", "application/json")

	recipe := r.Context().Value(KeyRecipe{}).(models.Recipe)
	user := r.Context().Value(models.UserKey{}).(models.User)

	storedRecipe, _ := handler.dao.Get(id)

	if storedRecipe.UserId != user.ID {
		http.Error(rw, "Not yours to change", http.StatusBadRequest)
		return
	}

	handler.stepDao.Delete(storedRecipe.ID)

	for i := 0; i < len(recipe.Steps); i++ {
		step := &recipe.Steps[i]
		step.ID = uuid.New()
		step.RecipeId = storedRecipe.ID
		handler.stepDao.Insert(step)
	}

	updated, err := handler.dao.Update(storedRecipe.ID, user.ID, &recipe)

	if err == gorm.ErrRecordNotFound {
		http.Error(rw, "Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(rw, "Not Found", http.StatusInternalServerError)
		return
	}

	response, _ := handler.dao.GetById(updated.ID)

	err = response.Serialize(rw)

	if err != nil {
		http.Error(rw, "Serialization error", http.StatusInternalServerError)
		return
	}
}

func (handler *RecipeHandler) UpdateView(rw http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(mux.Vars(r)["id"])
	handler.logger.Printf("PUT id: %s", id.String())

	rw.Header().Set("Content-Type", "application/json")

	response, _ := handler.dao.UpdateView(id)

	err := response.Serialize(rw)

	if err != nil {
		http.Error(rw, "Serialization error", http.StatusInternalServerError)
		return
	}
}

func (handler *RecipeHandler) GetSingle(rw http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(mux.Vars(r)["id"])
	handler.logger.Printf("GET id: %s", id.String())
	rw.Header().Set("Content-Type", "application/json")

	recipes, err := handler.dao.GetById(id)

	if err != nil {
		fmt.Print(err)
		http.Error(rw, "Fetch error", http.StatusInternalServerError)
		return
	}

	err = recipes.Serialize(rw)

	if err != nil {
		http.Error(rw, "Serialization error", http.StatusInternalServerError)
		return
	}
}

func (handler *RecipeHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	user := r.Context().Value(models.UserKey{}).(models.User)

	recipes, err := handler.dao.GetByUser(user.ID)

	if err != nil {
		fmt.Print(err)
		http.Error(rw, "Fetch error", http.StatusInternalServerError)
		return
	}

	err = recipes.Serialize(rw)

	if err != nil {
		http.Error(rw, "Serialization error", http.StatusInternalServerError)
		return
	}
}

func (handler *RecipeHandler) Get(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Print("GET")
	rw.Header().Set("Content-Type", "application/json")

	var recipes models.Recipes
	var err error
	if keyword, ok := mux.Vars(r)["search"]; ok {
		recipes, err = handler.dao.GetAll("%" + strings.ToLower(keyword) + "%")
	} else {
		recipes, err = handler.dao.GetAll("*")
	}

	if err != nil {
		fmt.Print(err)
		http.Error(rw, "Fetch error", http.StatusInternalServerError)
		return
	}

	err = recipes.Serialize(rw)

	if err != nil {
		http.Error(rw, "Serialization error", http.StatusInternalServerError)
		return
	}
}

func (handler *RecipeHandler) Delete(rw http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(mux.Vars(r)["id"])
	handler.logger.Printf("DELETE id: %s", id.String())
	rw.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(models.UserKey{}).(models.User)

	recipe, _ := handler.dao.Get(id)

	if recipe.UserId != user.ID {
		http.Error(rw, "Not yours to change", http.StatusBadRequest)
		return
	}

	handler.stepDao.Delete(id)
	affectedRows, err := handler.dao.Delete(id, user.ID)

	if err == gorm.ErrRecordNotFound || affectedRows == 0 {
		http.Error(rw, "Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(rw, "Not Found", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
