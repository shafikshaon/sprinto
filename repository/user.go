package repository

import (
	"gorm.io/gorm"
	"sprinto/models"
)

type UserRepository interface {
	Create(u models.User) error
	ByEmail(email string) (models.User, error)
	ByID(id uint) (models.User, error)
}

type userRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepo{db: db} }

func (r *userRepo) Create(u models.User) error { return r.db.Create(&u).Error }

func (r *userRepo) ByEmail(email string) (models.User, error) {
	var u models.User
	result := r.db.Where("email = ?", email).First(&u)
	return u, result.Error
}

func (r *userRepo) ByID(id uint) (models.User, error) {
	var u models.User
	result := r.db.First(&u, id)
	return u, result.Error
}
