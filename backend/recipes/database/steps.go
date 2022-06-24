package database

import (
	"recipes/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StepDao struct {
	Conn Database
}

func (dao *StepDao) Insert(item *models.Step) (tx *gorm.DB) {
	return dao.Conn.Db.Create(item)
}

func (dao *StepDao) GetForRecipe(recipeId uuid.UUID) (models.Steps, error) {
	var steps models.Steps
	tx := dao.Conn.Db.Find(&steps, "recipe_id = ?", recipeId)
	return steps, tx.Error
}

func (dao *StepDao) Upsert(recipeId uuid.UUID, item *models.Step) (*models.Step, error) {
	var step models.Step
	tx := dao.Conn.Db.FirstOrCreate(&step, "recipe_id = ?", recipeId)

	if tx.RowsAffected == 0 {
		step.Description = item.Description
		tx = dao.Conn.Db.Save(&step)
	}

	return &step, tx.Error
}

func (dao *StepDao) Delete(recipeId uuid.UUID) (int64, error) {
	tx := dao.Conn.Db.Where("recipe_id = ?", recipeId).Delete(&models.Step{})
	return tx.RowsAffected, tx.Error
}
