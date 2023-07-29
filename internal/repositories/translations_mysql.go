package repositories

import (
	"errors"

	"github.com/prplx/wordy/internal/models"
	"gorm.io/gorm"
)

type TranslationsRepo struct {
	db *gorm.DB
}

func NewTranslationsRepository(db *gorm.DB) *TranslationsRepo {
	return &TranslationsRepo{
		db: db,
	}
}

func (r *TranslationsRepo) QueryByExpressionID(expressionID int) ([]models.Translation, error) {
	var translations []models.Translation
	result := r.db.Find(&translations, "expression_id = ?", expressionID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return translations, models.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return translations, nil
}

func (r *TranslationsRepo) Create(translations []models.Translation) (int64, error) {
	result := r.db.Create(translations)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
