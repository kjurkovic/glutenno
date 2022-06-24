package handlers

import (
	"comments/config"
	"comments/database"
	"comments/models"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// types
type CommentHandler struct {
	logger *log.Logger
	dao    *database.CommentsDao
	config *config.Config
}

type KeyComment struct{}

func Comments(l *log.Logger, dao *database.CommentsDao, config *config.Config) *CommentHandler {
	return &CommentHandler{
		logger: l,
		dao:    dao,
		config: config,
	}
}

func (handler *CommentHandler) Add(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Print("POST")

	rw.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(models.UserKey{}).(models.User)
	comment := r.Context().Value(KeyComment{}).(models.Comment)

	comment.ID = uuid.New()
	comment.UserId = user.ID
	comment.User = user.Name

	fmt.Printf("Received object: %v", comment)
	_, err := handler.dao.Insert(&comment)

	if err != nil {
		http.Error(rw, "Unable to add", http.StatusInternalServerError)
		return
	}

	err = comment.Serialize(rw)

	if err != nil {
		http.Error(rw, "Serialization error", http.StatusInternalServerError)
		return
	}
}

func (handler *CommentHandler) Update(rw http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(mux.Vars(r)["id"])
	handler.logger.Printf("PUT id: %s", id.String())

	rw.Header().Set("Content-Type", "application/json")

	comment := r.Context().Value(KeyComment{}).(models.Comment)
	user := r.Context().Value(models.UserKey{}).(models.User)

	storedComment, _ := handler.dao.GetById(id)

	if storedComment.UserId != user.ID {
		http.Error(rw, "Not yours to change", http.StatusBadRequest)
		return
	}

	updated, err := handler.dao.Update(storedComment.ID, user.ID, &comment)

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

func (handler *CommentHandler) Get(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Print("GET")
	rw.Header().Set("Content-Type", "application/json")

	resourceId := uuid.MustParse(mux.Vars(r)["id"])
	comments, err := handler.dao.Get(resourceId)

	if err != nil {
		fmt.Print(err)
		http.Error(rw, "Fetch error", http.StatusInternalServerError)
		return
	}

	err = comments.Serialize(rw)

	if err != nil {
		http.Error(rw, "Serialization error", http.StatusInternalServerError)
		return
	}
}

func (handler *CommentHandler) GetPerUserRecipe(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Print("GET")
	rw.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(models.UserKey{}).(models.User)

	comments, err := handler.dao.GetByResourceOwnerId(user.ID)

	if err != nil {
		fmt.Print(err)
		http.Error(rw, "Fetch error", http.StatusInternalServerError)
		return
	}

	err = comments.Serialize(rw)

	if err != nil {
		http.Error(rw, "Serialization error", http.StatusInternalServerError)
		return
	}
}

func (handler *CommentHandler) Delete(rw http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(mux.Vars(r)["id"])
	handler.logger.Printf("DELETE id: %s", id.String())
	rw.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(models.UserKey{}).(models.User)

	comment, _ := handler.dao.GetById(id)

	if comment.UserId != user.ID {
		http.Error(rw, "Not yours to change", http.StatusBadRequest)
		return
	}

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
