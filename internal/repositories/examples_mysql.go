package repositories

import (
	"errors"

	"github.com/prplx/wordy/internal/models"
	"gorm.io/gorm"
)

type ExamplesRepo struct {
	db *gorm.DB
}

func NewExamplesRepository(db *gorm.DB) *ExamplesRepo {
	return &ExamplesRepo{
		db: db,
	}
}

func (r *ExamplesRepo) Create(examples []models.Example) (int64, error) {
	result := r.db.Create(examples)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (r *ExamplesRepo) QueryByExpressionID(expressionID int) ([]models.Example, error) {
	var examples []models.Example
	result := r.db.Find(&examples, "expression_id = ?", expressionID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return examples, models.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return examples, nil
}
