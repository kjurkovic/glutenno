package database

import (
	"comments/models"

	"github.com/google/uuid"
)

type CommentsDao struct {
	Conn Database
}

func (dao *CommentsDao) Insert(item *models.Comment) (int64, error) {
	tx := dao.Conn.Db.Create(item)
	return tx.RowsAffected, tx.Error
}

func (dao *CommentsDao) GetById(id uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	tx := dao.Conn.Db.First(&comment, "id = ?", id)
	return &comment, tx.Error
}

func (dao *CommentsDao) Get(resourceId uuid.UUID) (models.Comments, error) {
	var comments models.Comments
	tx := dao.Conn.Db.Find(&comments, "resource_id = ?", resourceId)

	return comments, tx.Error
}

func (dao *CommentsDao) GetByResourceOwnerId(ownerId uuid.UUID) (models.Comments, error) {
	var comments models.Comments
	tx := dao.Conn.Db.Find(&comments, "resource_owner_id = ?", ownerId)
	return comments, tx.Error
}

func (dao *CommentsDao) Update(id uuid.UUID, userId uuid.UUID, item *models.Comment) (*models.Comment, error) {
	var comment models.Comment
	tx := dao.Conn.Db.First(&comment, "id = ? AND user_id = ?", id, userId)

	if tx.Error != nil {
		return nil, tx.Error
	}

	comment.Text = item.Text

	tx = dao.Conn.Db.Save(&comment)
	return &comment, tx.Error
}

func (dao *CommentsDao) Delete(id uuid.UUID, userId uuid.UUID) (int64, error) {
	tx := dao.Conn.Db.Where("id = ? AND user_id = ?", id, userId).Delete(&models.Comment{})
	return tx.RowsAffected, tx.Error
}
