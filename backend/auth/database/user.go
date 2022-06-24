package database

import (
	"auth/models"

	"github.com/google/uuid"
)

type UserDao struct {
	Conn Database
}

func (dao *UserDao) Insert(item *models.User) (int64, error) {
	tx := dao.Conn.Db.Create(item)
	return tx.RowsAffected, tx.Error
}

func (dao *UserDao) Get(id uuid.UUID) (*models.User, error) {
	var user models.User
	tx := dao.Conn.Db.First(&user, "id = ?", id)
	return &user, tx.Error
}

func (dao *UserDao) GetByEmail(email string) (*models.User, error) {
	var user models.User
	tx := dao.Conn.Db.First(&user, "email = ?", email)
	return &user, tx.Error
}

func (dao *UserDao) GetUserByRefreshToken(token string) (*models.User, error) {
	var user models.User
	tx := dao.Conn.Db.First(&user, "refresh_token = ?", token)
	return &user, tx.Error
}

func (dao *UserDao) GetUserByForgotPasswordToken(token string) (*models.User, error) {
	var user models.User
	tx := dao.Conn.Db.First(&user, "forgot_password_token = ?", token)
	return &user, tx.Error
}

func (dao *UserDao) GetAll() ([]models.User, error) {
	var users []models.User
	tx := dao.Conn.Db.Find(&users)
	return users, tx.Error
}

func (dao *UserDao) Update(id uuid.UUID, item *models.User) (*models.User, error) {
	var user models.User
	tx := dao.Conn.Db.First(&user, "id = ?", id)

	if tx.Error != nil {
		return nil, tx.Error
	}

	// on user only name and email are allowed to be updated
	tx = dao.Conn.Db.Model(&user).Where("id = ?", id).Updates(models.User{Name: item.Name, Email: item.Email})
	return &user, tx.Error
}

func (dao *UserDao) UpdateRefreshToken(id uuid.UUID, refreshToken string) (*models.User, error) {
	var user models.User
	tx := dao.Conn.Db.First(&user, "id = ?", id)

	if tx.Error != nil {
		return nil, tx.Error
	}

	tx = dao.Conn.Db.Model(&user).Where("id = ?", id).Updates(models.User{RefreshToken: refreshToken})
	return &user, tx.Error
}

func (dao *UserDao) UpdateForgotPasswordToken(id uuid.UUID, token string) (*models.User, error) {
	var user models.User
	tx := dao.Conn.Db.First(&user, "id = ?", id)

	if tx.Error != nil {
		return nil, tx.Error
	}

	tx = dao.Conn.Db.Model(&user).Where("id = ?", id).Updates(models.User{ForgotPasswordToken: token})
	return &user, tx.Error
}

func (dao *UserDao) UpdatePassword(id uuid.UUID, password string) (*models.User, error) {
	var user models.User
	tx := dao.Conn.Db.First(&user, "id = ?", id)

	if tx.Error != nil {
		return nil, tx.Error
	}

	tx = dao.Conn.Db.Model(&user).Where("id = ?", id).Updates(models.User{Password: password})
	return &user, tx.Error
}

func (dao *UserDao) Delete(id uuid.UUID) (int64, error) {
	tx := dao.Conn.Db.Where("id = ?", id).Delete(&models.User{})
	return tx.RowsAffected, tx.Error
}
