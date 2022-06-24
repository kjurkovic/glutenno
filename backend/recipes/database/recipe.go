package database

import (
	"recipes/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecipesDao struct {
	Conn Database
}

func (dao *RecipesDao) Insert(item *models.Recipe) (int64, error) {
	tx := dao.Conn.Db.Create(item)
	return tx.RowsAffected, tx.Error
}

func (dao *RecipesDao) Get(id uuid.UUID) (*models.Recipe, error) {
	var recipe models.Recipe
	tx := dao.Conn.Db.First(&recipe, "id = ?", id)
	return &recipe, tx.Error
}

func (dao *RecipesDao) GetAll(keyword string) (models.Recipes, error) {
	var recipes models.Recipes

	var tx *gorm.DB
	if keyword == "*" {
		tx = dao.Conn.Db.Preload("Steps").Find(&recipes)
	} else {
		tx = dao.Conn.Db.Find(&recipes, "LOWER(title) LIKE ? OR LOWER(description) LIKE ?", keyword, keyword)
	}

	return recipes, tx.Error
}

func (dao *RecipesDao) GetByUser(userId uuid.UUID) (models.Recipes, error) {
	var recipes models.Recipes
	tx := dao.Conn.Db.Preload("Steps").Find(&recipes, "user_id = ?", userId)
	return recipes, tx.Error
}

func (dao *RecipesDao) GetById(recipeId uuid.UUID) (models.Recipe, error) {
	var recipe models.Recipe
	tx := dao.Conn.Db.Preload("Steps").Find(&recipe, "id = ?", recipeId)
	return recipe, tx.Error
}

func (dao *RecipesDao) Update(id uuid.UUID, userId uuid.UUID, item *models.Recipe) (*models.Recipe, error) {
	var recipe models.Recipe
	tx := dao.Conn.Db.First(&recipe, "id = ? AND user_id = ?", id, userId)

	if tx.Error != nil {
		return nil, tx.Error
	}

	recipe.Title = item.Title
	recipe.Description = item.Description

	tx = dao.Conn.Db.Save(&recipe)
	return &recipe, tx.Error
}

func (dao *RecipesDao) UpdateView(id uuid.UUID) (*models.Recipe, error) {
	var recipe models.Recipe
	tx := dao.Conn.Db.First(&recipe, "id = ?", id)

	if tx.Error != nil {
		return nil, tx.Error
	}

	recipe.Views += 1

	tx = dao.Conn.Db.Save(&recipe)
	return &recipe, tx.Error
}

func (dao *RecipesDao) Delete(id uuid.UUID, userId uuid.UUID) (int64, error) {
	tx := dao.Conn.Db.Where("id = ? AND user_id = ?", id, userId).Delete(&models.Recipe{})
	return tx.RowsAffected, tx.Error
}
