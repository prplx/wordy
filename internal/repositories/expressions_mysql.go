package repositories

import (
	"errors"

	"github.com/prplx/wordy/internal/models"
	"gorm.io/gorm"
)

type ExpressionsRepo struct {
	db *gorm.DB
}

func NewExpressionsRepository(db *gorm.DB) *ExpressionsRepo {
	return &ExpressionsRepo{
		db: db,
	}
}

func (r *ExpressionsRepo) GetByText(text string) (models.Expression, error) {
	var expression models.Expression
	result := r.db.First(&expression, "text = ?", text)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return expression, models.ErrRecordNotFound
	}

	return expression, result.Error
}

func (r *ExpressionsRepo) GetByTextWithAllData(text string) (models.Expression, error) {
	var expression models.Expression
	result := r.db.Preload("Translations").Preload("Examples").Preload("Audio").Preload("Synonyms").First(&expression, "text = ?", text)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return expression, models.ErrRecordNotFound
	}

	return expression, result.Error
}

func (r *ExpressionsRepo) Create(expression *models.Expression) (uint, error) {
	result := r.db.Create(&expression)

	return expression.ID, result.Error
}
