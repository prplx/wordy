package repositories

import (
	"errors"

	"github.com/prplx/wordy/internal/models"
	"gorm.io/gorm"
)

type UsersRepo struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) Create(user *models.User) (uint, error) {
	result := r.db.Create(&user)
	return user.ID, result.Error
}

func (r *UsersRepo) Get(id uint) (models.User, error) {
	var user models.User
	result := r.db.First(&user, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, models.ErrRecordNotFound
	}

	return user, result.Error
}

func (r *UsersRepo) GetByTgId(id uint) (models.User, error) {
	var user models.User
	result := r.db.First(&user, "telegram_id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, models.ErrRecordNotFound
	}

	return user, result.Error
}

func (r *UsersRepo) Update(user *models.User) error {
	result := r.db.Model(&user).Updates(user)
	return result.Error
}
