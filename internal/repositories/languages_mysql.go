package repositories

import (
	"errors"

	"github.com/prplx/wordy/internal/models"
	"gorm.io/gorm"
)

type LanguagesRepo struct {
	db *gorm.DB
}

func NewLanguagesRepository(db *gorm.DB) *LanguagesRepo {
	return &LanguagesRepo{
		db: db,
	}
}

func (r *LanguagesRepo) Query() ([]models.Language, error) {
	var languages []models.Language
	result := r.db.Find(&languages)
	if result.Error != nil {
		return nil, result.Error
	}
	return languages, nil
}

func (r *LanguagesRepo) GetByCode(code string) (models.Language, error) {
	var language models.Language
	result := r.db.First(&language, "code = ?", code)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return language, models.ErrRecordNotFound
	}

	return language, result.Error
}
